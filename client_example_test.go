package rawfasthttp_test

import (
	"log"

	"github.com/vidocsecurity/runner/pkg/rawfasthttp"
)

func ExampleHostClient() {
	// Perpare a client, which fetches webpages via HTTP proxy listening
	// on the localhost:8080.
	c := &rawfasthttp.HostClient{
		Addr: "localhost:8080",
	}

	// Fetch google page via local proxy.
	statusCode, body, err := c.Get(nil, "http://google.com/foo/bar")
	if err != nil {
		log.Fatalf("Error when loading google page through local proxy: %s", err)
	}
	if statusCode != rawfasthttp.StatusOK {
		log.Fatalf("Unexpected status code: %d. Expecting %d", statusCode, rawfasthttp.StatusOK)
	}
	useResponseBody(body)

	// Fetch foobar page via local proxy. Reuse body buffer.
	statusCode, body, err = c.Get(body, "http://foobar.com/google/com")
	if err != nil {
		log.Fatalf("Error when loading foobar page through local proxy: %s", err)
	}
	if statusCode != rawfasthttp.StatusOK {
		log.Fatalf("Unexpected status code: %d. Expecting %d", statusCode, rawfasthttp.StatusOK)
	}
	useResponseBody(body)
}

func useResponseBody(body []byte) {
	// Do something with body :)
}
