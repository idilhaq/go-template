package http

import (
	"github.com/go-chi/chi"

	"github.com/idilhaq/go-template/src/app"
	healthcheck "github.com/idilhaq/go-template/src/app/delivery/http/handler/healthcheck"
	home "github.com/idilhaq/go-template/src/app/delivery/http/handler/home"
	"github.com/idilhaq/go-template/src/app/delivery/http/middleware"
)

type Server struct {
	handlers []HandlerItf
}

//go:generate mockgen -source=http.go -destination=http_mock.go -package=http
type HandlerItf interface {
	RegisterHandler(router *chi.Mux)
}

func NewServer(usecase *app.AppUsecaseDepedency) Server {
	return Server{
		handlers: []HandlerItf{
			healthcheck.InitHandler(usecase),
			home.InitHandler(usecase),
		},
	}
}

func (s Server) RegisterHandler() *chi.Mux {
	router := chi.NewRouter()

	// router.Use(panics.CaptureHandler)
	router.Use(middleware.RequestLoggerMiddleware)

	for _, handler := range s.handlers {
		handler.RegisterHandler(router)
	}

	return router
}
