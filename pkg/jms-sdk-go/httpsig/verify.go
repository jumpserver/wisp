package httpsig

import (
	"bytes"
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

type verifier func(string, []byte) (bool, error)

// VerifySignature takes a ParsedSignature and a public key, and validates that
// the signature was signed with the given algorithm and a corresponding private
// key.
func VerifySignature(parsedSignature *ParsedSignature, pubKey string) (bool, error) {
	v, err := getVerifier(parsedSignature.Algorithm(), pubKey)
	if err != nil {
		return false, err
	}
	sig, err := base64.StdEncoding.DecodeString(parsedSignature.Signature())
	if err != nil {
		return false, err
	}
	return v(parsedSignature.SigningString(), sig)
}

func hmacVerifier(secret string, hash crypto.Hash) verifier {
	return func(data string, sig []byte) (bool, error) {
		h := hmac.New(hash.New, []byte(secret))
		h.Write([]byte(data))
		expected := h.Sum(nil)
		return bytes.Equal(expected, sig), nil
	}
}

func rsaVerifier(key *rsa.PublicKey, hash crypto.Hash) verifier {
	return func(data string, sig []byte) (bool, error) {
		hashed := calcHash(data, hash)
		err := rsa.VerifyPKCS1v15(key, hash, hashed, sig)
		return err == nil, err
	}
}

func dsaVerifier(key *dsa.PublicKey, hash crypto.Hash) verifier {
	return func(data string, sig []byte) (bool, error) {
		hashed := calcHash(data, hash)
		qlen := len(key.Q.Bytes())
		if len(hashed) > qlen {
			hashed = hashed[:qlen]
		}
		s := dsaSignature{}
		if _, err := asn1.Unmarshal(sig, &s); err != nil {
			return false, err
		}
		return dsa.Verify(key, hashed, s.R, s.S), nil
	}
}

func ecdsaVerifier(key *ecdsa.PublicKey, hash crypto.Hash) verifier {
	return func(data string, sig []byte) (bool, error) {
		hashed := calcHash(data, hash)
		s := dsaSignature{}
		if _, err := asn1.Unmarshal(sig, &s); err != nil {
			return false, err
		}
		return ecdsa.Verify(key, hashed, s.R, s.S), nil
	}
}

func getVerifier(algorithm string, pubKey string) (verifier, error) {
	alg, err := validateAlgorithm(algorithm)
	if err != nil {
		return nil, err
	}
	if alg.sign == "hmac" {
		return hmacVerifier(pubKey, alg.hash), nil
	}
	k, err := getPublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	switch key := k.(type) {
	case *rsa.PublicKey:
		if alg.sign != "rsa" {
			return nil, fmt.Errorf("Algorithm %s doesn't match public key of type %T", algorithm, key)
		}
		return rsaVerifier(key, alg.hash), nil
	case *dsa.PublicKey:
		if alg.sign != "dsa" {
			return nil, fmt.Errorf("Algorithm %s doesn't match public key of type %T", algorithm, key)
		}
		return dsaVerifier(key, alg.hash), nil
	case *ecdsa.PublicKey:
		if alg.sign != "ecdsa" {
			return nil, fmt.Errorf("Algorithm %s doesn't match public key of type %T", algorithm, key)
		}
		return ecdsaVerifier(key, alg.hash), nil
	}
	return nil, fmt.Errorf("Unsupported signing algorithm: %s", algorithm)
}

func getPublicKey(key string) (interface{}, error) {
	block, _ := pem.Decode([]byte(key))
	return x509.ParsePKIXPublicKey(block.Bytes)
}
