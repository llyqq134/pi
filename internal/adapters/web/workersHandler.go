package web

import (
	"log"
	"net/http"
	"pi/internal/adapters/payloads"
	"pi/internal/app/interfaces/services"

	"github.com/gin-gonic/gin"
)

const (
	loginURL    = "/login"      //POST
	listWorkers = "/admin/list" //GROUP
)

type workerHandler struct {
	workerService services.WorkerService
}

func NewHandler(workerService services.WorkerService) Handler {
	return &workerHandler{workerService: workerService}
}

func (h *workerHandler) Register(router *gin.Engine) {
	workers := router.Group("listworkers")
	{
		workers.POST(loginURL)
		workers.GET("/:department", h.ListDepartmentHandler)
		workers.POST("/:department", h.Addworker)
		workers.PUT("/:id")
		workers.DELETE("/:id", h.DeleteWorker)
	}
}

func (h *workerHandler) ListDepartmentHandler(c *gin.Context) {
	departemnt := c.Param("department")

	if len(departemnt) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})

		return
	}

	workers, err := h.workerService.GetAllWorkersByDepartment(departemnt)
	if err != nil {
		log.Printf("Failed to get list of workers for %v with error: %v\n", departemnt, err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})

		return
	}

	c.JSON(http.StatusOK, workers)
}

func (h *workerHandler) Addworker(c *gin.Context) {
	departmentName := c.Param("department")
	var worker payloads.WorkerPayload

	if err := c.BindJSON(&worker); err != nil {
		log.Printf("server error: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad json body",
		})

		return
	}

	newWorker, err := h.workerService.CreateWorker(worker.Name, worker.JobTitle, departmentName, worker.Password)
	if err != nil {
		log.Printf("server error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "new worker was added",
		"worker":  newWorker,
	})
}

func (h *workerHandler) DeleteWorker(c *gin.Context) {
	workerUUID := c.Param("id")

	if len(workerUUID) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})

		return
	}

	if err := h.workerService.DeleteWorkerByUUID(workerUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "worker was successfully deleted",
	})
}
