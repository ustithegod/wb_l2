package main

import (
	"calendar/event"
	h "calendar/http"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	repo := event.NewRepository()
	router := h.NewRouter(repo)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("running http server")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
