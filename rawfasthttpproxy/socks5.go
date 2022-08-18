package rawfasthttpproxy

import (
	"net"
	"net/url"

	"github.com/vidocsecurity/runner/pkg/rawfasthttp"

	"golang.org/x/net/proxy"
)

// FasthttpSocksDialer returns a rawfasthttp.DialFunc that dials using
// the provided SOCKS5 proxy.
//
// Example usage:
//	c := &rawfasthttp.Client{
//		Dial: fasthttpproxy.FasthttpSocksDialer("socks5://localhost:9050"),
//	}
func FasthttpSocksDialer(proxyAddr string) rawfasthttp.DialFunc {
	var (
		u      *url.URL
		err    error
		dialer proxy.Dialer
	)
	if u, err = url.Parse(proxyAddr); err == nil {
		dialer, err = proxy.FromURL(u, proxy.Direct)
	}
	// It would be nice if we could return the error here. But we can't
	// change our API so just keep returning it in the returned Dial function.
	// Besides the implementation of proxy.SOCKS5() at the time of writing this
	// will always return nil as error.

	return func(addr string) (net.Conn, error) {
		if err != nil {
			return nil, err
		}
		return dialer.Dial("tcp", addr)
	}
}
