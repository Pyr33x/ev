package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pyr33x/ev/pkg/config"
)

type Server struct {
	port int
}

func New(cfg *config.Server) *http.Server {
	srv := &Server{
		port: cfg.Port,
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", srv.port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
