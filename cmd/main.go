package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cyberdr0id/astro/internal/handler"
	"github.com/cyberdr0id/astro/internal/repository"
	"github.com/cyberdr0id/astro/internal/service"
	"github.com/cyberdr0id/astro/internal/storage"
)

func main() {
	s3, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}

	db, err := repository.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.New(db)
	service := service.New(repo, s3)

	s := handler.NewServer(service)

	if os.Getenv("APP_PORT") == "" {
		log.Fatal("empty application port")
	}

	if err := http.ListenAndServe(":"+os.Getenv("APP_PORT"), s); err != nil {
		log.Fatal("unable to listen server:", err)
	}
}
