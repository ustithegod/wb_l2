package main

import (
	"calendar/event"
	h "calendar/http"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	repo := event.NewRepository()
	router := h.NewRouter(repo)

	log.Info().Msg("loading config from .env file")
	srvConfig, err := NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("error while loading config")
	}
	

	srv := &http.Server{
		Addr:        srvConfig.HttpAddress,
		Handler:     router,
		ReadTimeout: time.Duration(srvConfig.ReadTimeout),
		IdleTimeout: time.Duration(srvConfig.IdleTimeout),
	}

	fmt.Printf("running http server on address: %s\n", srv.Addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Str("server address", srv.Addr).Msg("error while starting server")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info().Msg("closing server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(srvConfig.ShutdownTimeout))
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("cant gracefully shutdown http server")
	}
}
