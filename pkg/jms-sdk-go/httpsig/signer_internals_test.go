package httpsig

import (
	"github.com/stretchr/testify/assert"
	u "net/url"
	"testing"
)

func doGetPathAndQueryFromURL(urlString string) string {
	url, _ := u.Parse(urlString)
	return getPathAndQueryFromURL(url)
}

func TestGetPathAndQueryFromURL(t *testing.T) {
	assert.Equal(t, "/", doGetPathAndQueryFromURL("http://example.com"))
	assert.Equal(t, "/", doGetPathAndQueryFromURL("http://example.com/"))
	assert.Equal(t, "/?withquery=string", doGetPathAndQueryFromURL("http://example.com/?withquery=string"))
	assert.Equal(t, "/test/path", doGetPathAndQueryFromURL("http://example.com/test/path"))
	assert.Equal(t, "/test/path?withquery=string%20andspaces", doGetPathAndQueryFromURL("http://example.com/test/path?withquery=string%20andspaces"))
}
