package reuseport_test

import (
	"fmt"
	"log"

	"github.com/vidocsecurity/runner/pkg/rawfasthttp"
	"github.com/vidocsecurity/runner/pkg/rawfasthttp/reuseport"
)

func ExampleListen() {
	ln, err := reuseport.Listen("tcp4", "localhost:12345")
	if err != nil {
		log.Fatalf("error in reuseport listener: %s", err)
	}

	if err = rawfasthttp.Serve(ln, requestHandler); err != nil {
		log.Fatalf("error in fasthttp Server: %s", err)
	}
}

func requestHandler(ctx *rawfasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, world!")
}
