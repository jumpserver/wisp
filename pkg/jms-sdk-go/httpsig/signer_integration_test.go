// +build integration

package httpsig

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"os/exec"
	"testing"
)

func TestClientCanCallNodeServer(t *testing.T) {
	npmInstall()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	ready := make(chan *exec.Cmd)
	go startNodeServer(t, port, ready)
	cmd := <-ready
	defer cmd.Process.Kill()

	verifyClientCanCallNodeServer(t, port, "hmac-sha1")
	verifyClientCanCallNodeServer(t, port, "hmac-sha256")
	verifyClientCanCallNodeServer(t, port, "hmac-sha512")
	verifyClientCanCallNodeServer(t, port, "rsa-sha1")
	verifyClientCanCallNodeServer(t, port, "rsa-sha256")
	verifyClientCanCallNodeServer(t, port, "rsa-sha512")
	verifyClientCanCallNodeServer(t, port, "dsa-sha1")
	verifyClientCanCallNodeServer(t, port, "dsa-sha256")
	verifyClientCanCallNodeServer(t, port, "dsa-sha512")
	verifyClientCanCallNodeServer(t, port, "ecdsa-sha1")
	verifyClientCanCallNodeServer(t, port, "ecdsa-sha256")
	verifyClientCanCallNodeServer(t, port, "ecdsa-sha512")
}

func verifyClientCanCallNodeServer(t *testing.T, port string, algorithm string) {
	t.Logf("Calling node server with %s algorithm", algorithm)
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:%s/", port), nil)
	signer, _ := NewRequestSigner(getKeyIDForTests(algorithm), getPrivateKeyForTests(algorithm), algorithm)
	err := signer.SignRequest(req, []string{"date", "(request-target)"}, nil)
	if err != nil {
		t.Error(err)
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, http.StatusOK, res.StatusCode)
	t.Log(res.StatusCode)
}

func getKeyIDForTests(alg string) string {
	algorithm, _ := validateAlgorithm(alg)
	if algorithm.sign == "hmac" {
		return getPrivateKeyForTests(alg)
	}
	return fmt.Sprintf("%s_public.pem", algorithm.sign)
}

func startNodeServer(t *testing.T, port string, ch chan *exec.Cmd) {
	r, w := io.Pipe()
	read := bufio.NewReader(r)

	cmd := exec.Command("node", "server.js", port)
	cmd.Dir = "./test"
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Start()

	out := make(chan string, 12)
	go logOut(t, read, out)

	<-out // wait for Listening

	ch <- cmd

	return
}

func logOut(t *testing.T, r *bufio.Reader, out chan string) {
	for {
		line, e := r.ReadString('\n')
		if e != nil {
			t.Log(e)
			return
		}
		t.Logf("server: %s", line)
		out <- line
	}
}

func npmInstall() {
	cmd := exec.Command("npm", "install")
	cmd.Dir = "./test"
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
