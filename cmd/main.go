package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-clean-arch/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Info("No .env file found", log.Ferror(err))
	}

	var addr string
	flag.StringVar(&addr, "addr", ":8083", "tcp host:port to connect")
	flag.Parse()

	mainCtx, cancelMain := context.WithCancel(context.Background())
	defer cancelMain()

	// This is framework-agnostic and can be replaced with any HTTP framework like net/http, Gin, or Echo
	// The following example uses net/http.
	container, err := HTTPBuildContainer(mainCtx)
	if err != nil {
		log.Critical("Failed to build container", log.Ferror(err))
		return
	}

	err = container.Invoke(func(router *chi.Mux, config *config.ServerConfig) {
		srv := &http.Server{
			Addr:         addr,
			Handler:      router,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			IdleTimeout:  config.IdleTimeout,
		}
		log.Info("Server running...")

		signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
		defer stop()

		go func() {
			if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Error("Server failed", log.Ferror(err))
				return
			}
		}()

		<-signalCtx.Done()
		log.Info("Server stopping...")

		tctx, cancelShutdown := context.WithTimeout(context.Background(), config.GracefulShutdownTimeout)
		defer cancelShutdown()

		if err = srv.Shutdown(tctx); err != nil {
			log.Error("Failed to shutdown http server", log.Ferror(err))
		}
		log.Info("Server exited")
	})
	if err != nil {
		log.Critical("Failed to start server", log.Ferror(err))
		return
	}
}
