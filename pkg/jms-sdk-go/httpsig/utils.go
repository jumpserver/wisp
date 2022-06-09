package httpsig

import (
	"crypto"
	"crypto/dsa"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

var validAlgorithms = map[string]bool{
	"hmac":  true,
	"rsa":   true,
	"dsa":   true,
	"ecdsa": true,
}

var validHashAlgorithms = map[string]crypto.Hash{
	"sha1":   crypto.SHA1,
	"sha256": crypto.SHA256,
	"sha512": crypto.SHA512,
}

type hashAlgorithm struct {
	sign string
	hash crypto.Hash
}

type dsaSignature struct {
	R, S *big.Int
}

func autoDetectAlgorithm(key string) (*hashAlgorithm, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, errors.New("Could not determine key format, key was not in PEM format")
	}
	switch block.Type {
	case "RSA PRIVATE KEY":
		return &hashAlgorithm{"rsa", crypto.SHA256}, nil
	case "DSA PRIVATE KEY":
		return &hashAlgorithm{"dsa", crypto.SHA256}, nil
	case "EC PRIVATE KEY":
		return &hashAlgorithm{"ecdsa", crypto.SHA256}, nil
	default:
		return nil, fmt.Errorf("Could not determine key format (pem block type '%s')", block.Type)
	}

}

func validateAlgorithm(algorithm string) (*hashAlgorithm, error) {
	alg := strings.Split(strings.ToLower(algorithm), "-")
	if len(alg) != 2 {
		return nil, fmt.Errorf("%s is not a valid algorithm", strings.ToUpper(algorithm))
	}
	if hash, ok := validHashAlgorithms[alg[1]]; ok {
		if _, ok := validAlgorithms[alg[0]]; !ok {
			return nil, fmt.Errorf("%s type keys are not supported", strings.ToUpper(alg[0]))
		}
		return &hashAlgorithm{alg[0], hash}, nil
	}
	return nil, fmt.Errorf("%s is not a supported hash algorithm", strings.ToUpper(alg[0]))
}

func (alg *hashAlgorithm) String() string {
	return fmt.Sprintf("%s-%s", strings.ToLower(alg.sign), strings.ToLower(hashName(alg.hash)))
}

func hashName(hash crypto.Hash) string {
	switch hash {
	case crypto.SHA1:
		return "SHA1"
	case crypto.SHA256:
		return "sha256"
	case crypto.SHA512:
		return "sha512"
	}
	return "unknown"
}

func calcHash(data string, hash crypto.Hash) []byte {
	h := hash.New()
	h.Write([]byte(data))
	return h.Sum(nil)
}

type tempDsaKey struct {
	E1, P, Q, G, Y, X *big.Int
}

// PEM DSA format doesn't parse using any of the built-in crypto methods, so this function parses a DSA private key from a PEM file
func getDsaKey(key string) (privateKey *dsa.PrivateKey, err error) {
	block, _ := pem.Decode([]byte(key))

	tmpKey := tempDsaKey{}
	_, err = asn1.Unmarshal(block.Bytes, &tmpKey)
	if err != nil {
		return nil, err
	}
	privateKey = &dsa.PrivateKey{}
	privateKey.P = tmpKey.P
	privateKey.Q = tmpKey.Q
	privateKey.G = tmpKey.G
	privateKey.Y = tmpKey.Y
	privateKey.X = tmpKey.X
	return
}
