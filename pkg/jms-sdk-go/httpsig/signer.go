package httpsig

import (
	"crypto"
	"crypto/dsa"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	u "net/url"
	"strings"
	"time"
)

// A Signer is a function that takes a signing string, and returns a byte array
// containing the signature. Used when creating a custom RequestSigner
type Signer func(string) ([]byte, error)

// RequestSigner contains the SignRequest method which signs a request and
// populates the Authorization Header
type RequestSigner struct {
	keyID     string
	algorithm string
	signer    Signer
}

// SignStrict specifies strict signing. If true, then obsolete versions of the HTTP-Signature standard will not be supported.
// Currently, this controls whether to allow the "request-line" psuedo-header.
var SignStrict = false

// NewRequestSigner creates a new RequestSigner with the given keyID, key, and algorithm.
// If algorithm is the empty string, then the signing algorithm will be determined from the key type,
// and the hashing algorithm will be sha256. RSA, DSA, and ECDSA keys must be in PEM format.
func NewRequestSigner(keyID string, key string, algorithm string) (*RequestSigner, error) {
	var alg *hashAlgorithm
	var err error
	if algorithm == "" {
		alg, err = autoDetectAlgorithm(key)
		if err != nil {
			return nil, err
		}
	} else {
		alg, err = validateAlgorithm(algorithm)
		if err != nil {
			return nil, err
		}
	}
	signer, err := getSigner(alg, key)
	if err != nil {
		return nil, err
	}
	a := alg.String()
	return &RequestSigner{
		keyID:     keyID,
		algorithm: a,
		signer:    signer,
	}, nil
}

// NewCustomRequestSigner creates a new RequestSigner with a custom signing algorithm.
func NewCustomRequestSigner(keyID string, algorithm string, signer Signer) *RequestSigner {
	return &RequestSigner{
		keyID:     keyID,
		algorithm: algorithm,
		signer:    signer,
	}
}

// SignRequest signs a request, populating the Authorization Header with the resulting
// HTTP Signature
func (rs *RequestSigner) SignRequest(request *http.Request, headers []string, ext map[string]string) error {
	if _, ok := request.Header["Date"]; !ok {
		request.Header["Date"] = []string{time.Now().Format(time.RFC1123)}
	}
	if len(headers) == 0 {
		headers = []string{"date"}
	}
	lines := make([]string, 0, len(headers))
	for _, h := range headers {
		h = strings.ToLower(h)
		if h == "request-line" {
			if SignStrict {
				return errors.New("request-line is not a valid header with strict parsing enabled")
			}
			lines = append(lines, fmt.Sprintf("%s %s %s", request.Method, getPathAndQueryFromURL(request.URL), request.Proto))
		} else if h == "(request-target)" {
			lines = append(lines, fmt.Sprintf("(request-target): %s %s", strings.ToLower(request.Method), getPathAndQueryFromURL(request.URL)))
		} else if h == "host" {
			lines = append(lines, fmt.Sprintf("%s: %s", h, request.URL.Host))
		} else if h == "content-length" {
			lines = append(lines, fmt.Sprintf("%s: %d", h, request.ContentLength))
		} else {
			values, ok := request.Header[http.CanonicalHeaderKey(h)]
			if !ok {
				return fmt.Errorf("No value for header \"%s\"", h)
			}
			lines = append(lines, fmt.Sprintf("%s: %s", h, values[0]))
		}
	}
	stringToSign := strings.Join(lines, "\n")
	signature, err := rs.signer(stringToSign)
	if err != nil {
		return err
	}
	request.Header["Authorization"] = []string{formatSignature(rs.keyID, rs.algorithm, headers, ext, signature)}
	return nil
}

func getPathAndQueryFromURL(url *u.URL) (pathAndQuery string) {
	pathAndQuery = url.Path
	if pathAndQuery == "" {
		pathAndQuery = "/"
	}
	if url.RawQuery != "" {
		pathAndQuery += "?" + url.RawQuery
	}
	return pathAndQuery
}

func formatSignature(keyID string, algorithm string, headers []string, ext map[string]string, signature []byte) string {
	sig := fmt.Sprintf("Signature keyId=\"%s\",algorithm=\"%s\",headers=\"%s\"", keyID, algorithm, strings.Join(headers, " "))

	for key, val := range ext {
		sig += fmt.Sprintf(",%s=\"%s\"", key, val)
	}
	sig += fmt.Sprintf(",signature=\"%s\"", base64.StdEncoding.EncodeToString(signature))
	return sig
}

func hmacSigner(secret string, hash crypto.Hash) (Signer, error) {
	return func(data string) ([]byte, error) {
		h := hmac.New(hash.New, []byte(secret))
		h.Write([]byte(data))
		return h.Sum(nil), nil
	}, nil
}

func rsaSigner(key string, hash crypto.Hash) (Signer, error) {
	block, _ := pem.Decode([]byte(key))
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return noSigner, err
	}
	privateKey.Precompute()
	return func(data string) ([]byte, error) {
		hashed := calcHash(data, hash)
		return rsa.SignPKCS1v15(rand.Reader, privateKey, hash, hashed)
	}, nil
}

func dsaSigner(key string, hash crypto.Hash) (Signer, error) {
	privateKey, err := getDsaKey(key)
	if err != nil {
		return noSigner, err
	}
	return func(data string) ([]byte, error) {
		hashed := calcHash(data, hash)
		qlen := len(privateKey.Q.Bytes())
		if len(hashed) > qlen {
			hashed = hashed[:qlen]
		}
		r, s, err := dsa.Sign(rand.Reader, privateKey, hashed)
		if err != nil {
			return nil, err
		}

		return asn1.Marshal(dsaSignature{r, s})
	}, nil
}

func ecdsaSigner(key string, hash crypto.Hash) (Signer, error) {
	block, _ := pem.Decode([]byte(key))
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return noSigner, err
	}
	return func(data string) ([]byte, error) {
		hashed := calcHash(data, hash)
		return privateKey.Sign(rand.Reader, hashed, nil)
	}, nil
}

func getSigner(alg *hashAlgorithm, key string) (Signer, error) {
	switch alg.sign {
	case "hmac":
		return hmacSigner(key, alg.hash)
	case "rsa":
		return rsaSigner(key, alg.hash)
	case "dsa":
		return dsaSigner(key, alg.hash)
	case "ecdsa":
		return ecdsaSigner(key, alg.hash)
	}
	return nil, fmt.Errorf("Unsupported signing algorithm: %v", alg)
}

func noSigner(data string) ([]byte, error) {
	return nil, errors.New("Invalid signer")
}
