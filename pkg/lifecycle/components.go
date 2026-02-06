package lifecycle

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// HTTPServer returns a Group actor for a *http.Server.
func HTTPServer(srv *http.Server) (start func() error, interrupt func(error)) {
	start = func() error { return srv.ListenAndServe() }
	interrupt = func(error) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}
	return
}

// Signals returns an actor that cancels the provided cancel func on SIGINT/SIGTERM (and optional others).
func Signals(cancel context.CancelFunc, sigs ...os.Signal) (func() error, func(error)) {
	if len(sigs) == 0 {
		sigs = []os.Signal{os.Interrupt, syscall.SIGTERM}
	}
	return func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, sigs...)
			s := <-c // block until a signal arrives
			_ = s
			cancel()
			return nil
		},
		func(error) { cancel() }
}

// Webview actor: pass in a Run func and a Terminate func.
type Webviewer interface {
	Run()       // blocking
	Terminate() // stops Run()
}

func Webview(wv Webviewer) (func() error, func(error)) {
	return func() error {
			wv.Run()
			return nil
		},
		func(error) { wv.Terminate() }
}
