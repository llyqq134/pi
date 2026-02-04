package main

import (
	"context"
	"log"
	"pi/internal/adapters/repo"
	"pi/internal/adapters/web"
	"pi/internal/app/usecases"
	"pi/pkg/db"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	log.Println("Connecting to database...")
	client, err := db.NewClient(context.Background(), 3)
	if err != nil {
		panic(err)
	}

	workerRepo := repo.NewWorkerImpl(client)
	workerCase := usecases.NewWorkerService(&workerRepo)
	workerHandler := web.NewWorkerHandler(workerCase)

	workerHandler.Register(router)

	log.Println("Starting server...")
	router.Use()
	router.Run("localhost:8080")
}
