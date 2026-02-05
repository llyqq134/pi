package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"pi/internal/adapters/repo"
	"pi/internal/adapters/web"
	"pi/internal/app/usecases"
	"pi/pkg/db"

	"github.com/gin-gonic/gin"
)

func findBaseDir() string {
	// Ищем директорию с index.html (рядом с go.mod)
	dir, _ := os.Getwd()
	for i := 0; i < 5; i++ {
		if _, err := os.Stat(filepath.Join(dir, "index.html")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
		if dir == filepath.Dir(dir) {
			break
		}
	}
	return "." // fallback
}

func main() {
	baseDir := findBaseDir()
	staticDir := filepath.Join(baseDir, "static")
	log.Printf("Serving static from: %s", baseDir)

	router := gin.Default()

	// CORS middleware для работы фронтенда
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Раздача статических файлов (HTML, CSS, JS)
	router.Static("/static", staticDir)
	router.StaticFile("/", filepath.Join(baseDir, "index.html"))
	router.StaticFile("/dashboard.html", filepath.Join(baseDir, "dashboard.html"))
	router.StaticFile("/workers.html", filepath.Join(baseDir, "workers.html"))
	router.StaticFile("/departments.html", filepath.Join(baseDir, "departments.html"))
	router.StaticFile("/equipment.html", filepath.Join(baseDir, "equipment.html"))
	router.StaticFile("/records.html", filepath.Join(baseDir, "records.html"))
	router.StaticFile("/reports.html", filepath.Join(baseDir, "reports.html"))

	log.Println("Connecting to database...")
	client, err := db.NewClient(context.Background(), 3)
	if err != nil {
		panic(err)
	}

	// Workers
	workerRepo := repo.NewWorkerImpl(client)
	workerCase := usecases.NewWorkerService(&workerRepo)
	workerHandler := web.NewWorkerHandler(workerCase)
	workerHandler.Register(router)

	// Departments
	departmentRepo := repo.NewDepartmentImpl(client)
	departmentCase := usecases.NewDepartmentService(&departmentRepo, &workerRepo)
	departmentHandler := web.NewDepartmentHandler(departmentCase)
	departmentHandler.Register(router)

	// Equipment
	equipmentRepo := repo.NewEquipmentImpl(client)
	equipmentCase := usecases.NewEquipmentService(&equipmentRepo)
	equipmentHandler := web.NewEquipmentHandler(equipmentCase)
	equipmentHandler.Register(router)

	// Records
	recordRepo := repo.NewRecordsImpl(client)
	recordCase := usecases.NewRecordService(&recordRepo)
	recordHandler := web.NewRecordsHandler(recordCase)
	recordHandler.Register(router)

	log.Println("Starting server on localhost:8080...")
	router.Run("localhost:8080")
}
