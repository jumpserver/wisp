# httpsig

httpsig is a go package with for [HTTP Signature](https://tools.ietf.org/html/draft-cavage-http-signatures-05). It also implements [extensions](https://tools.ietf.org/html/draft-cavage-http-signatures-05#appendix-B) to the standard.

## Usage

```go
import "gopkg.in/twindagger/httpsig.v1"
```

### Client

This example signs a request and includes the date, and (request-target) header components in the signature.
```go
// set key as a string from file read, memory, etc.
req, _ := http.NewRequest("GET", "http://example.com/path/to/resource", nil)
signer, _ := httpsig.NewRequestSigner("my-key-id", key, "rsa-sha256")
err := signer.SignRequest(req, []string{"date", "(request-target)"}, jwt)
```


### Server

This example verifies that a request contains a signature and returns a 401 Unauthorized response if a signature is not present or not verifiable.

```go

func HandleReq(w http.ResponseWriter, r *http.Request) {
    parsed, err := ParseRequest(req)
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    publicKey := lookupPubKey(parsed.KeyId())
    verified, err := VerifySignature(parsed, publicKey)
    if err != nil || !verified {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write("Authoirzation Passed")
}

func main() {
    http.HandleFunc("/", HandleReq)
    http.ListenAndServe(":8080", nil)
}
```

## Installation

    go get gopkg.in/twindagger/httpsig.v1

## License

MIT.
