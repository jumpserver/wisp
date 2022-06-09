package httpsig

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"
)

const (
	SampleKeyID     = "keyId"
	SampleAlgorithm = "rsa-sha256"
	SampleHeaders   = "host date"
	SampleSignature = "abcdefg"
)

func TestParseRequestWithNoAuthorizationHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/path/to/resource", nil)
	parsed, err := ParseRequest(req)
	assert.NotNil(t, err)
	assert.Nil(t, parsed)
	assert.Contains(t, err.Error(), "no authorization header present")
}

func TestParseRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/path/to/resource", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Signature keyId=\"%s\",algorithm=\"%s\",headers=\"%s\",signature=\"%s\"", SampleKeyID, SampleAlgorithm, SampleHeaders, SampleSignature))
	req.Header.Set("Date", time.Now().Format(time.RFC1123))
	parsed, err := ParseRequest(req)
	assert.Nil(t, err)
	assert.NotNil(t, parsed)
	assert.Equal(t, "Signature", parsed.Scheme())
	assert.NotNil(t, parsed.Params())
	assert.Equal(t, SampleKeyID, parsed.KeyId())
	assert.Equal(t, strings.ToUpper(SampleAlgorithm), parsed.Algorithm())
	assert.Equal(t, SampleHeaders, strings.Join(parsed.Headers(), " "))
	assert.Equal(t, SampleSignature, parsed.Signature())

	rx, _ := regexp.Compile("host: example.com\ndate: [^\n]+")
	assert.Regexp(t, rx, parsed.SigningString())
}

func TestParseRequestWithNoHeaders(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/path/to/resource", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Signature keyId=\"%s\",algorithm=\"%s\",signature=\"%s\"", SampleKeyID, SampleAlgorithm, SampleSignature))
	req.Header.Set("Date", time.Now().Format(time.RFC1123))
	parsed, err := ParseRequest(req)
	assert.Nil(t, err)
	assert.NotNil(t, parsed)
	assert.Equal(t, "date", strings.Join(parsed.Headers(), " "))
	rx, _ := regexp.Compile("date: [^\n]+")
	assert.Regexp(t, rx, parsed.SigningString())
}

func TestParseRequestWithRequestLine(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/path/to/resource", nil)
	req.RequestURI = "/path/to/resource"
	req.Header.Set("Authorization", fmt.Sprintf("Signature keyId=\"%s\",algorithm=\"%s\",headers=\"%s request-line\",signature=\"%s\"", SampleKeyID, SampleAlgorithm, SampleHeaders, SampleSignature))
	req.Header.Set("Date", time.Now().Format(time.RFC1123))
	parsed, err := ParseRequest(req)
	assert.Nil(t, err)
	assert.NotNil(t, parsed)
	assert.Equal(t, fmt.Sprintf("%s %s", SampleHeaders, "request-line"), strings.Join(parsed.Headers(), " "))
	assert.Equal(t, SampleSignature, parsed.Signature())

	rx, _ := regexp.Compile("host: example.com\ndate: [^\n]+\nGET /path/to/resource HTTP/1.1")
	assert.Regexp(t, rx, parsed.SigningString())
}

func TestParseRequestWithRequestTarget(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/path/to/resource", nil)
	req.RequestURI = "/path/to/resource"
	req.Header.Set("Authorization", fmt.Sprintf("Signature keyId=\"%s\",algorithm=\"%s\",headers=\"%s (request-target)\",signature=\"%s\"", SampleKeyID, SampleAlgorithm, SampleHeaders, SampleSignature))
	req.Header.Set("Date", time.Now().Format(time.RFC1123))
	parsed, err := ParseRequest(req)
	assert.Nil(t, err)
	assert.NotNil(t, parsed)
	assert.Equal(t, fmt.Sprintf("%s %s", SampleHeaders, "(request-target)"), strings.Join(parsed.Headers(), " "))
	assert.Equal(t, SampleSignature, parsed.Signature())

	rx, _ := regexp.Compile("host: example.com\ndate: [^\n]+\n\\(request-target\\): get /path/to/resource")
	assert.Regexp(t, rx, parsed.SigningString())
}

func TestParseRequestWithDateBeforeClockSkewRange(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/path/to/resource", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Signature keyId=\"%s\",algorithm=\"%s\",signature=\"%s\"", SampleKeyID, SampleAlgorithm, SampleSignature))
	req.Header.Set("Date", time.Now().Add(-DefaultClockSkew-time.Second).Format(time.RFC1123))
	parsed, err := ParseRequest(req)
	assert.NotNil(t, err)
	assert.Nil(t, parsed)
	assert.Contains(t, err.Error(), "Expired Request error")
}

func TestParseRequestWithDateAfterClockSkewRange(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/path/to/resource", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Signature keyId=\"%s\",algorithm=\"%s\",signature=\"%s\"", SampleKeyID, SampleAlgorithm, SampleSignature))
	req.Header.Set("Date", time.Now().Add(DefaultClockSkew+time.Second).Format(time.RFC1123))
	parsed, err := ParseRequest(req)
	assert.NotNil(t, err)
	assert.Nil(t, parsed)
	assert.Contains(t, err.Error(), "Expired Request error")
}
