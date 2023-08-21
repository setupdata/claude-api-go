package claudeapi

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"
)

const (
	idleConnsPerHost = 25
	defaultTimeout   = 30 * time.Second
)

var defaultHeaders = map[string]string{
	"Content-Type": "application/json",
	"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
	"Referer":      "https://claude.ai/chats",
	"Origin":       "https://claude.ai",
}

type RequestHandler func(req *http.Request) error
type ClaudeOption func(claude *Claude) error

type Claude struct {
	organizationId string
	cookies        []*http.Cookie
	baseReqHandler RequestHandler
	client         *http.Client
}

// WithCookies set cookies
func WithCookies(cookieInfos map[string]string) ClaudeOption {
	return func(claude *Claude) error {
		for key, value := range cookieInfos {
			claude.cookies = append(claude.cookies, &http.Cookie{
				Name:  key,
				Value: value,
			})
		}
		return nil
	}
}

func WithBaseHttpHandler(reqHandler RequestHandler) ClaudeOption {
	return func(claude *Claude) error {
		claude.baseReqHandler = reqHandler
		return nil
	}
}

func (c *Claude) getOrganizationId(apiReqHandles ...RequestHandler) (string, error) {
	url := "https://claude.ai/api/organizations"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	for key, value := range defaultHeaders {
		req.Header.Add(key, value)
	}

	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	if c.baseReqHandler != nil {
		err := c.baseReqHandler(req)
		if err != nil {
			return "", err
		}
	}

	if len(apiReqHandles) != 0 {
		for _, apiReqHandle := range apiReqHandles {
			err := apiReqHandle(req)
			if err != nil {
				return "", err
			}
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	var organizations Organizations
	err = json.NewDecoder(resp.Body).Decode(&organizations)
	if err != nil {
		return "", err
	}

	if len(organizations) == 0 {
		return "", errors.New("organizations is empty")
	}
	return organizations[0].Uuid, nil
}

func NewClaude(config *Config, options ...ClaudeOption) (*Claude, error) {
	proxy := http.ProxyFromEnvironment
	if config.Proxy != nil {
		proxy = config.Proxy
	}

	var tr http.RoundTripper
	if config.Transport != nil {
		tr = config.Transport
	} else {
		tr = &http.Transport{
			Proxy:               proxy,
			TLSHandshakeTimeout: 10 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConnsPerHost: idleConnsPerHost,
		}
	}

	timeout := defaultTimeout
	if config.Timeout != 0 {
		timeout = config.Timeout
	}

	claude := &Claude{
		client: &http.Client{
			Transport: tr,
			Timeout:   timeout,
		},
	}

	for _, option := range options {
		if err := option(claude); err != nil {
			return nil, err
		}
	}

	if claude.cookies == nil {
		return nil, errors.New("cookies is nil")
	}
	organizationId, err := claude.getOrganizationId()
	if err != nil {
		return nil, err
	}

	claude.organizationId = organizationId

	return claude, nil
}
