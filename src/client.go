package src

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type RequestOptions struct {
	Method     string
	URL        string
	Headers    map[string]string
	Body       []byte
	ProxyURL   string
	Timeout    time.Duration
	SkipVerify bool
	BasicAuth  *BasicAuth
}

type BasicAuth struct {
	Username string
	Password string
}

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
	Cookies    []*http.Cookie
	Error      error
}

func MakeRequest(opts RequestOptions) (*Response, error) {
	if opts.Method == "" {
		opts.Method = "GET"
	}
	if opts.Timeout == 0 {
		opts.Timeout = 30 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	jar, err := cookiejar.New(nil)
	if err != nil {
		return &Response{Error: err}, err
	}

	client := &http.Client{
		Timeout: opts.Timeout,
		Jar:     jar,
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()

	if opts.ProxyURL != "" {
		var proxyURL *url.URL
		var err error

		if strings.Contains(opts.ProxyURL, "@") && !strings.HasPrefix(opts.ProxyURL, "http://") && !strings.HasPrefix(opts.ProxyURL, "https://") {
			parts := strings.Split(opts.ProxyURL, "@")
			if len(parts) == 2 {
				auth := parts[0]
				hostPort := parts[1]

				formattedProxyURL := fmt.Sprintf("http://%s@%s", auth, hostPort)
				proxyURL, err = url.Parse(formattedProxyURL)
			} else {
				err = fmt.Errorf("invalid proxy format: %s", opts.ProxyURL)
			}
		} else {
			proxyURL, err = url.Parse(opts.ProxyURL)
		}

		if err != nil {
			return &Response{Error: err}, err
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	if opts.SkipVerify {
		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		} else {
			transport.TLSClientConfig.InsecureSkipVerify = true
		}
	}

	client.Transport = transport

	req, err := http.NewRequestWithContext(ctx, opts.Method, opts.URL, bytes.NewReader(opts.Body))
	if err != nil {
		return &Response{Error: err}, err
	}

	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	if opts.BasicAuth != nil {
		req.SetBasicAuth(opts.BasicAuth.Username, opts.BasicAuth.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return &Response{Error: err}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Response{Error: err}, err
	}

	headers := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	cookies := resp.Cookies()

	response := &Response{
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Body:       body,
		Cookies:    cookies,
	}

	return response, nil
}
