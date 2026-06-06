package http_adapter

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/caio-bernardo/dragonite/internal/delivery/http_adapter/client"
)

// Registra os endpoints do servidor
func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	clientHandler := client.NewHandler(s.usuarioService)
	// federationHandler := federation.NewHandler()

	// Registra rotas
	mux.HandleFunc("GET /health", s.healthHandler)
	clientHandler.RegisterRoutes(mux, s.TokenBearerMiddleware)
	// federationHandler.RegisterRoutes(mux)

	// wildcard
	mux.HandleFunc("GET /", s.HelloWorldHandler)

	// Adiciona middlewares
	// NOTE: a ordem dos middleware importa! O mais interno é chamado primeiro.
	return s.logMiddleware(s.corsMiddleware(mux))
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Hello World"}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(s.systemService.PingDB())
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
