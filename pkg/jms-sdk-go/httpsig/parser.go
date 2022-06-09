package httpsig

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// DefaultClockSkew specifies the allowed time difference between requests and signature parsing
var DefaultClockSkew = 300 * time.Second

// DefaultHeaders specify the headers to be parsed by default
var DefaultHeaders = []string{"date"}

// ParseStrict specifies sprict parsing. If true, then obsolete versions of the HTTP-Signature standard will not be supported.
// Currently, this controls whether to allow the "request-line" psuedo-header.
var ParseStrict = false

const (
	stateNew    = 0
	stateParams = 1
)

const (
	paramsStateName  = 0
	paramsStateQuote = 1
	paramsStateValue = 2
	paramsStateComma = 3
)

// ParsedSignature contains the details of a parsed signature
type ParsedSignature struct {
	scheme        string
	params        map[string]string
	signingString string
}

// ParseRequest parses an http.Request and returns a ParsedSignature with the values
// from the Authorization Header.
func ParseRequest(request *http.Request) (*ParsedSignature, error) {
	state := stateNew
	substate := paramsStateName
	tmpName := ""
	tmpValue := ""

	parsed := ParsedSignature{
		params: make(map[string]string),
	}

	var authz string
	if auth, ok := request.Header["Authorization"]; ok {
		authz = auth[0]
	} else {
		return nil, errors.New("Missing Header error: no authorization header present in the request")
	}

	for _, code := range authz {
		c := string(code)
		switch state {
		case stateNew:
			if c != " " {
				parsed.scheme += c
			} else {
				state = stateParams
			}
			break
		case stateParams:
			switch substate {
			case paramsStateName:
				if (code >= 0x41 && code <= 0x5a) || // A-Z
					(code >= 0x61 && code <= 0x7a) { // a-z
					tmpName += c
				} else if c == "=" {
					if tmpName == "" {
						return nil, errors.New("Invalid Header Error: bad param format")
					}
					substate = paramsStateQuote
				} else {
					return nil, errors.New("Invalid Header Error: bad param format")
				}
				break
			case paramsStateQuote:
				if c == "\"" {
					tmpValue = ""
					substate = paramsStateValue
				} else {
					return nil, errors.New("Invalid Header Error: bad param format")
				}
				break
			case paramsStateValue:
				if c == "\"" {
					parsed.params[tmpName] = tmpValue
					substate = paramsStateComma
				} else {
					tmpValue += c
				}
				break
			case paramsStateComma:
				if c == "," {
					tmpName = ""
					substate = paramsStateName
				} else {
					return nil, errors.New("Invalid Header Error: bad param format")
				}
				break
			}
			break
		}
	}

	var h string
	if val, ok := parsed.params["headers"]; ok {
		h = strings.ToLower(val)
	} else {
		h = "date"
	}
	parsed.params["headers"] = h
	headers := strings.Split(h, " ")

	if parsed.scheme != "Signature" {
		return nil, errors.New("Scheme was not \"Signature\"")
	}

	if _, ok := parsed.params["keyId"]; !ok {
		return nil, errors.New("Invalid Header Error: keyId was not specified")
	}

	if _, ok := parsed.params["algorithm"]; !ok {
		return nil, errors.New("Invalid Header Error: algorithm was not specified")
	}

	if _, ok := parsed.params["signature"]; !ok {
		return nil, errors.New("Invalid Header Error: signature was not specified'")
	}

	if _, err := validateAlgorithm(parsed.Algorithm()); err != nil {
		return nil, err
	}

	// Build the signingString
	var signingParts = make([]string, len(headers), len(headers))
	for i, h := range headers {
		if h == "request-line" {
			if !ParseStrict {
				/*
				 * We allow headers from the older spec drafts if strict parsing isn't
				 * specified in options.
				 */
				signingParts[i] = fmt.Sprintf("%s %s %s", request.Method, request.RequestURI, request.Proto)
			} else {
				/* Strict parsing doesn't allow older draft headers. */
				return nil, errors.New("Strict Parsing error: request-line is not a valid header with strict parsing enabled.")
			}
		} else if h == "(request-target)" {
			signingParts[i] = fmt.Sprintf("(request-target): %s %s", strings.ToLower(request.Method), request.RequestURI)
		} else if h == "host" {
			signingParts[i] = fmt.Sprintf("%s: %s", h, request.Host)
		} else if h == "content-length" {
			signingParts[i] = fmt.Sprintf("%s: %d", h, request.ContentLength)
		} else {
			if value, ok := request.Header[http.CanonicalHeaderKey(h)]; ok {
				signingParts[i] = h + ": " + value[0]
			} else {
				return nil, fmt.Errorf("Missing Header error: \"%s\" was not in the request.", h)
			}
		}
	}
	parsed.signingString = strings.Join(signingParts, "\n")

	var dateString string
	hasDate := false
	if ds, ok := request.Header["Date"]; ok {
		dateString = ds[0]
		hasDate = true
	}

	if hasDate {
		date, _ := time.Parse(time.RFC1123, dateString)
		now := time.Now()
		skew := now.Sub(date)
		if skew < 0 {
			skew = -skew
		}

		if skew > DefaultClockSkew {
			return nil, fmt.Errorf("Expired Request error: clock skew of %v was greater than %v", skew, DefaultClockSkew)
		}
	}

	for _, dh := range DefaultHeaders {
		found := false
		for _, h := range headers {
			if h == dh {
				found = true
			}
		}
		if !found {
			return nil, fmt.Errorf("%s was not a signed header", dh)
		}
	}

	return &parsed, nil
}

// Scheme is the scheme of the parsed signature. Currently, the only supported scheme is "Signature"
func (s *ParsedSignature) Scheme() string {
	return s.scheme
}

// Params is a map of sttrings containing the signature parameters
func (s *ParsedSignature) Params() map[string]string {
	return s.params
}

// SigningString is the string that was signed by the client. Verifiers should use this to check that a singature is valid
func (s *ParsedSignature) SigningString() string {
	return s.signingString
}

// Algorithm is the signing and hash algorithm, like "rsa-sha256"
func (s *ParsedSignature) Algorithm() string {
	return strings.ToUpper(s.params["algorithm"])
}

// KeyId is the keyId parameter from the signature
//
// Deprecated: since 1.1
func (s *ParsedSignature) KeyId() string {
	return s.params["keyId"]
}

// KeyID is the keyId parameter from the signature
func (s *ParsedSignature) KeyID() string {
	return s.params["keyId"]
}

// Signature is the base-64 encoded signature
func (s *ParsedSignature) Signature() string {
	return s.params["signature"]
}

// Headers is the slice of headers that were signed by the client
func (s *ParsedSignature) Headers() []string {
	return strings.Split(s.params["headers"], " ")
}
