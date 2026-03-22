package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"taskmanager/internal/framework/http/handler"
)

type Router struct {
	Task *handler.TaskHandler
}

func (rt *Router) Handler() http.Handler {
	r := chi.NewRouter()
	// Keep dependencies minimal and Go 1.20 compatible.
	// Add middlewares later as needed (logging, request-id, recover, etc.).

	// basic health
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	if rt.Task != nil {
		rt.Task.Register(r)
	}

	return r
}
