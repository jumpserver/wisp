// +build integration

package httpsig

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"
)

func TestNodeClientCanCallServer(t *testing.T) {
	npmInstall()

	handler := &TestHandler{t}
	server := httptest.NewServer(handler)
	address := server.Listener.Addr().String()
	defer server.Close()

	verifyNodeClientCanCallServer(t, address, "hmac-sha1")
	verifyNodeClientCanCallServer(t, address, "hmac-sha256")
	verifyNodeClientCanCallServer(t, address, "hmac-sha512")
	verifyNodeClientCanCallServer(t, address, "rsa-sha1")
	verifyNodeClientCanCallServer(t, address, "rsa-sha256")
	verifyNodeClientCanCallServer(t, address, "rsa-sha512")
	verifyNodeClientCanCallServer(t, address, "dsa-sha1")
	verifyNodeClientCanCallServer(t, address, "dsa-sha256")
	verifyNodeClientCanCallServer(t, address, "dsa-sha512")
	verifyNodeClientCanCallServer(t, address, "ecdsa-sha1")
	verifyNodeClientCanCallServer(t, address, "ecdsa-sha256")
	verifyNodeClientCanCallServer(t, address, "ecdsa-sha512")
}

func verifyNodeClientCanCallServer(t *testing.T, address string, algorithm string) {
	t.Logf("Calling go server with %s algorithm\n", algorithm)
	output := runNodeClient(t, address, algorithm)
	t.Log(output)
	assert.Equal(t, "200\n", output)
}

func runNodeClient(t *testing.T, address string, algorithm string) string {
	cmd := exec.Command("node", "client.js", address, algorithm)
	cmd.Dir = "./test"
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	rs, _ := out.ReadString(byte(0))
	if err != nil {
		t.Log(rs)
		panic(err)
	}
	return rs
}

type TestHandler struct {
	*testing.T
}

func (h *TestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	parsed, err := ParseRequest(req)
	if err != nil {
		h.Log(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	publicKey := getPublicKeyForTests(parsed.Algorithm(), parsed.KeyId())
	verified, err := VerifySignature(parsed, publicKey)
	if err != nil || !verified {
		h.Logf("Unverified: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Authoirzation Passed"))
}

func getPublicKeyForTests(alg string, keyID string) string {
	algorithm, _ := validateAlgorithm(alg)
	if algorithm.sign == "hmac" {
		return getPrivateKeyForTests(alg)
	}
	return getTestKey(keyID)
}
