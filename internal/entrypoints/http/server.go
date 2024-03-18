package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sivistrukov/vk-assigment/internal/models"
)

type authService interface {
	Authenticate(string, string) (models.User, error)
}

// NewServer returns new http server instance
func NewServer(cfg Config, handler http.Handler) http.Server {
	return http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
}