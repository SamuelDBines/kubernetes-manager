package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/SamuelDBines/kubernetes-manager/pkg/httpserver"
	"github.com/SamuelDBines/kubernetes-manager/pkg/httpserver/handlers/health"
	"github.com/SamuelDBines/kubernetes-manager/pkg/utils/env"
	"github.com/SamuelDBines/kubernetes-manager/pkg/utils/lifecycle"
)

type Config struct {
	ServerConfig httpserver.Config
	// DatabaseConfig database.Config
}

func main() {
	_, err := env.LoadDefault(&env.Options{Overwrite: false, Expand: true})
	if err != nil {
		log.Fatal(err)
	}

	print("Environment variables loaded successfully\n")

	var cfg Config = Config{
		ServerConfig: httpserver.Config{
			Port: env.Int("MEDIA_SERVICE_PORT", 3333),
			Name: "",
		},
	}

	mux := http.NewServeMux()
	srv := httpserver.NewServer(cfg.ServerConfig, mux)
	health.Routes(mux)

	log.Printf("🚀 %s starting on http://localhost%s (pid %d)",
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
