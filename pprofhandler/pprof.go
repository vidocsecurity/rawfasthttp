package pprofhandler

import (
	"net/http/pprof"
	rtp "runtime/pprof"
	"strings"

	"github.com/vidocsecurity/rawfasthttp"
	rawfasthttpadaptor "github.com/vidocsecurity/rawfasthttp/fasthttpadaptor"
)

var (
	cmdline = rawfasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Cmdline)
	profile = rawfasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Profile)
	symbol  = rawfasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Symbol)
	trace   = rawfasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Trace)
	index   = rawfasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Index)
)

// PprofHandler serves server runtime profiling data in the format expected by the pprof visualization tool.
//
// See https://golang.org/pkg/net/http/pprof/ for details.
func PprofHandler(ctx *rawfasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html")
	if strings.HasPrefix(string(ctx.Path()), "/debug/pprof/cmdline") {
		cmdline(ctx)
	} else if strings.HasPrefix(string(ctx.Path()), "/debug/pprof/profile") {
		profile(ctx)
	} else if strings.HasPrefix(string(ctx.Path()), "/debug/pprof/symbol") {
		symbol(ctx)
	} else if strings.HasPrefix(string(ctx.Path()), "/debug/pprof/trace") {
		trace(ctx)
	} else {
		for _, v := range rtp.Profiles() {
			ppName := v.Name()
			if strings.HasPrefix(string(ctx.Path()), "/debug/pprof/"+ppName) {
				namedHandler := rawfasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Handler(ppName).ServeHTTP)
				namedHandler(ctx)
				return
			}
		}
		index(ctx)
	}
}
