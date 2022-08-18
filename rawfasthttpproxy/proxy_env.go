package rawfasthttpproxy

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net"
	"net/url"
	"sync/atomic"
	"time"

	"golang.org/x/net/http/httpproxy"

	"github.com/vidocsecurity/runner/pkg/rawfasthttp"
)

const (
	httpsScheme = "https"
	httpScheme  = "http"
	tlsPort     = "443"
)

// FasthttpProxyHTTPDialer returns a rawfasthttp.DialFunc that dials using
// the the env(HTTP_PROXY, HTTPS_PROXY and NO_PROXY) configured HTTP proxy.
//
// Example usage:
//	c := &rawfasthttp.Client{
//		Dial: FasthttpProxyHTTPDialer(),
//	}
func FasthttpProxyHTTPDialer() rawfasthttp.DialFunc {
	return FasthttpProxyHTTPDialerTimeout(0)
}

// FasthttpProxyHTTPDialer returns a rawfasthttp.DialFunc that dials using
// the env(HTTP_PROXY, HTTPS_PROXY and NO_PROXY) configured HTTP proxy using the given timeout.
//
// Example usage:
//	c := &rawfasthttp.Client{
//		Dial: FasthttpProxyHTTPDialerTimeout(time.Second * 2),
//	}
func FasthttpProxyHTTPDialerTimeout(timeout time.Duration) rawfasthttp.DialFunc {
	proxier := httpproxy.FromEnvironment().ProxyFunc()

	// encoded auth barrier for http and https proxy.
	authHTTPStorage := &atomic.Value{}
	authHTTPSStorage := &atomic.Value{}

	return func(addr string) (net.Conn, error) {

		port, _, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, fmt.Errorf("unexpected addr format: %v", err)
		}

		reqURL := &url.URL{Host: addr, Scheme: httpScheme}
		if port == tlsPort {
			reqURL.Scheme = httpsScheme
		}
		proxyURL, err := proxier(reqURL)
		if err != nil {
			return nil, err
		}

		if proxyURL == nil {
			if timeout == 0 {
				return rawfasthttp.Dial(addr)
			}
			return rawfasthttp.DialTimeout(addr, timeout)
		}

		var conn net.Conn
		if timeout == 0 {
			conn, err = rawfasthttp.Dial(proxyURL.Host)
		} else {
			conn, err = rawfasthttp.DialTimeout(proxyURL.Host, timeout)
		}
		if err != nil {
			return nil, err
		}

		req := "CONNECT " + addr + " HTTP/1.1\r\n"

		if proxyURL.User != nil {
			authBarrierStorage := authHTTPStorage
			if port == tlsPort {
				authBarrierStorage = authHTTPSStorage
			}

			auth := authBarrierStorage.Load()
			if auth == nil {
				authBarrier := base64.StdEncoding.EncodeToString([]byte(proxyURL.User.String()))
				auth := &authBarrier
				authBarrierStorage.Store(auth)
			}

			req += "Proxy-Authorization: Basic " + *auth.(*string) + "\r\n"
		}
		req += "\r\n"

		if _, err := conn.Write([]byte(req)); err != nil {
			return nil, err
		}

		res := rawfasthttp.AcquireResponse()
		defer rawfasthttp.ReleaseResponse(res)

		res.SkipBody = true

		if err := res.Read(bufio.NewReader(conn)); err != nil {
			if connErr := conn.Close(); connErr != nil {
				return nil, fmt.Errorf("conn close err %v followed by read conn err %v", connErr, err)
			}
			return nil, err
		}
		if res.Header.StatusCode() != 200 {
			if connErr := conn.Close(); connErr != nil {
				return nil, fmt.Errorf(
					"conn close err %v followed by connect to proxy: code: %d body %s",
					connErr, res.StatusCode(), string(res.Body()))
			}
			return nil, fmt.Errorf("could not connect to proxy: code: %d body %s", res.StatusCode(), string(res.Body()))
		}
		return conn, nil
	}
}
