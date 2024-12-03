package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

func newRouter() *httprouter.Router {

	mux := httprouter.New()
	ytAPIKey := os.Getenv("YOUTUBE_API_KEY")
	ytChannelID := os.Getenv("YOUTUBE_CHANNEL_ID")

	if ytAPIKey == "" {
		log.Fatal("Key not provided")
	}

	if ytChannelID == "" {
		log.Fatal("youtube channel ID not provided")
	}

	mux.GET("/youtube/channel/stats", getChannelStats(ytAPIKey, ytChannelID))

	return mux
}

func main() {
	srv := &http.Server{
		Addr:    ":10101",
		Handler: newRouter(),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		log.Println("Interrupt recieved")

		log.Println("http server shutting down")
		time.Sleep(5 * time.Second)

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("http server shutdown error: %v", err)
		}

		log.Println("Shutdown completed")

		close(idleConnsClosed)
	}()

	log.Printf("Starting server on port 10101")
	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start: %v", err)
		}
	}

	<-idleConnsClosed
	log.Println("Server Stop")
}
