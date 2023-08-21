package claudeapi

import (
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	// Transport may be used for custom HTTP behavior.
	Transport http.RoundTripper

	// The maximum length of time to wait before giving up on a server request. A value of zero means no timeout.
	Timeout time.Duration

	// Proxy is the proxy func to be used for all requests made by this
	// transport. If Proxy is nil, http.ProxyFromEnvironment is used. If Proxy
	// returns a nil *url.URL, no proxy is used.
	//
	// socks5 proxying does not currently support spdy streaming endpoints.
	Proxy func(*http.Request) (*url.URL, error)
}
