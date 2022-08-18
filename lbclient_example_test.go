package rawfasthttp_test

import (
	"fmt"
	"log"

	"github.com/vidocsecurity/rawfasthttp"
)

func ExampleLBClient() {
	// Requests will be spread among these servers.
	servers := []string{
		"google.com:80",
		"foobar.com:8080",
		"127.0.0.1:123",
	}

	// Prepare clients for each server
	var lbc rawfasthttp.LBClient
	for _, addr := range servers {
		c := &rawfasthttp.HostClient{
			Addr: addr,
		}
		lbc.Clients = append(lbc.Clients, c)
	}

	// Send requests to load-balanced servers
	var req rawfasthttp.Request
	var resp rawfasthttp.Response
	for i := 0; i < 10; i++ {
		url := fmt.Sprintf("http://abcedfg/foo/bar/%d", i)
		req.SetRequestURI(url)
		if err := lbc.Do(&req, &resp); err != nil {
			log.Fatalf("Error when sending request: %s", err)
		}
		if resp.StatusCode() != rawfasthttp.StatusOK {
			log.Fatalf("unexpected status code: %d. Expecting %d", resp.StatusCode(), rawfasthttp.StatusOK)
		}

		useResponseBody(resp.Body())
	}
}
