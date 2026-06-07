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
	dirService     usecase.DirectoryService
	profileService usecase.ProfileService
	syncService    usecase.SyncService
}

// Cria um novo servidor http
func NewServer(port int,
	jwtSecret string,
	usuarioService usecase.UsuarioService,
	systemService usecase.HealthService,
	dirService usecase.DirectoryService,
	profileService usecase.ProfileService,
	syncService usecase.SyncService,
) *http.Server {

	NewServer := &Server{
		port:           port,
		jwtSecret:      jwtSecret,
		usuarioService: usuarioService,
		systemService:  systemService,
		dirService:     dirService,
		profileService: profileService,
		syncService:    syncService,
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
