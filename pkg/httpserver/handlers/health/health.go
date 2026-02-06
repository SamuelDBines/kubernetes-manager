package health

import (
	"net/http"

	"github.com/SamuelDBines/golang-framework/pkg/httpserver"
)

func Routes(m *http.ServeMux) {
	httpserver.With(m, "/health", http.HandlerFunc(HealthCheckHandler))
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	httpserver.OK(w,map[string]string{"status": "ok"})
}