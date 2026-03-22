package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"taskmanager/internal/framework/http/handler"
	"taskmanager/internal/framework/http/router"
	"taskmanager/internal/infra/memory"
	"taskmanager/internal/usecase"
)

func main() {
	// Ports/Adapters wiring
	repo := memory.NewTaskRepository()
	uc := usecase.NewTaskUsecase(repo)

	taskHandler := handler.NewTaskHandler(uc)
	rt := &router.Router{Task: taskHandler}

	addr := envOr("TASKMANAGER_ADDR", ":8080")
	srv := &http.Server{
		Addr:              addr,
		Handler:           rt.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("taskmanager listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("server error: %v", err)
		os.Exit(1)
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
