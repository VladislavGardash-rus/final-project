package http_server

import (
	"context"
	"fmt"
	http_handlers "github.com/gardashvs/final-project/internal/handlers/http"
	"github.com/gardashvs/final-project/internal/logger"
	"github.com/gardashvs/final-project/internal/transport/http/middleware"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Server struct {
	helloWorldHandler *http_handlers.PreviewerHandler
	srv               *http.Server
	router            *mux.Router
	alias             string
}

func NewServer(addr string, alias string) *Server {
	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	server := &Server{
		alias:             alias,
		srv:               srv,
		helloWorldHandler: http_handlers.NewPreviewerHandler(),
	}
	server.initRoting()

	return server
}

func (s *Server) initRoting() {
	s.router = mux.NewRouter()
	s.router.Use(middleware.Logging)
	s.router.Host("cut-service.com")
	s.router.StrictSlash(true)
	s.router.HandleFunc("/fill/{width:[0-9]+}/{height:[0-9]+}/{url:.*?(?:\\s|$)}", middleware.Serve(s.helloWorldHandler.GetPreview)).Methods(http.MethodGet)

	s.srv.Handler = s.router
}

func (s *Server) Start() error {
	logger.UseLogger().Info(fmt.Sprintf("http: Server %s started on ", s.alias), s.srv.Addr)
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
