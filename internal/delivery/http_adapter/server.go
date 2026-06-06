package http_adapter

import (
	"net/http"
	"strconv"
	"time"

	"github.com/caio-bernardo/dragonite/internal/usecase"
)

// AppServer representa o servidor em nível de aplicação
type Server struct {
	port           int
	jwtSecret      string
	usuarioService usecase.UsuarioService
	systemService  usecase.HealthService
}

// Cria um novo servidor http
func NewServer(port int, jwtSecret string, usuarioService usecase.UsuarioService, systemService usecase.HealthService) *http.Server {

	NewServer := &Server{
		port:           port,
		jwtSecret:      jwtSecret,
		usuarioService: usuarioService,
		systemService:  systemService,
	}

	// servidor http, com endpoints registrados e timeout para operações R/W
	server := http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &server
}
