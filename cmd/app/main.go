package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/SamuelDBines/kubernetes-manager/pkg/env"
	"github.com/SamuelDBines/kubernetes-manager/pkg/httpserver"
	"github.com/SamuelDBines/kubernetes-manager/pkg/httpserver/handlers/health"
	ui "github.com/SamuelDBines/kubernetes-manager/pkg/httpserver/handlers/ui"
	"github.com/SamuelDBines/kubernetes-manager/pkg/lifecycle"
	"github.com/SamuelDBines/kubernetes-manager/pkg/store"
	"github.com/SamuelDBines/kubernetes-manager/pkg/web"
)

type Config struct {
	ServerConfig httpserver.Config
}

func main() {
	_, err := env.LoadDefault(&env.Options{Overwrite: false, Expand: true})
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config = Config{
		ServerConfig: httpserver.Config{
			Port: env.Int("MEDIA_SERVICE_PORT", 3333),
			Name: "kubernetes-manager",
		},
	}

	const outDir = "out"

	// Ensure output folder exists
	if err := store.EnsureOut(outDir); err != nil {
		log.Fatal(err)
	}

	// Templates live in ./templates
	renderer, err := web.NewRenderer("templates")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	srv := httpserver.NewServer(cfg.ServerConfig, mux)

	health.Routes(mux)

	// Static assets (css/js) served from templates/
	mux.Handle("/static/css/",
		http.StripPrefix("/static/css/", http.FileServer(http.Dir(filepath.Join("templates", "css")))),
	)
	mux.Handle("/static/js/",
		http.StripPrefix("/static/js/", http.FileServer(http.Dir(filepath.Join("templates", "js")))),
	)

	// UI routes
	mux.HandleFunc("/", ui.Index(renderer, outDir))

	// Optional: serve generated files (handy for debugging)
	mux.Handle("/out/",
		http.StripPrefix("/out/", http.FileServer(http.Dir(outDir))),
	)

	log.Printf("🚀 %s starting on http://localhost:%d (pid %d)",
		cfg.ServerConfig.Name, cfg.ServerConfig.Port, os.Getpid())

	var g lifecycle.Group
	ctx, cancel := context.WithCancel(context.Background())

	g.Add(lifecycle.Signals(cancel))

	startHTTP, stopHTTP := lifecycle.HTTPServer(srv)
	g.Add(startHTTP, stopHTTP)

	g.Add(func() error { <-ctx.Done(); return nil }, func(error) { cancel() })

	if err := g.Run(); err != nil && err != http.ErrServerClosed {
		log.Printf("exit: %v", err)
	}
}
