package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kikudesuyo/arumonogohan-app/api/route"
)

func main() {
	_ = godotenv.Load(".env")
	port := flag.String("port", "8081", "port to run HTTP server on")
	flag.Parse()
	if envPort := os.Getenv("PORT"); envPort != "" {
		*port = envPort
	}

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", *port),
		Handler:           route.NewRouter(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("API server is running on port %s", *port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
