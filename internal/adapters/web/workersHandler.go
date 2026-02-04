package web

import (
	"log"
	"net/http"
	"pi/internal/adapters/payloads"
	"pi/internal/app/interfaces/services"
	"pi/pkg/utils/auth"
	"pi/pkg/utils/middleware"

	"github.com/gin-gonic/gin"
)

const (
	loginURL    = "/login"      //POST
	listWorkers = "/admin/list" //GROUP
)

type workerHandler struct {
	workerService services.WorkerService
}

func NewWorkerHandler(workerService services.WorkerService) Handler {
	return &workerHandler{workerService: workerService}
}

func (h *workerHandler) Register(router *gin.Engine) {
	public := router.Group("/login")
	{
		public.POST("", h.LoginHandler)
	}

	protected := router.Group("/listworkers")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/department/:department", h.ListDepartmentHandler)
		protected.POST("", h.Addworker)
		protected.PUT("/:id")
		protected.DELETE("/:id", h.DeleteWorker)
	}
}

func (h *workerHandler) LoginHandler(c *gin.Context) {
	var loginPayload payloads.LoginPayload

	if err := c.ShouldBindJSON(&loginPayload); err != nil {
		log.Printf("error binding login request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "incorrect login",
		})
		return
	}

	worker, err := h.workerService.GetWorkerByName(loginPayload.Name)
	if err != nil {
		log.Printf("error getting worker: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "incorrect name or password",
		})
		return
	}

	if worker.Password != loginPayload.Password {
		log.Printf("invalid password for user: %s", loginPayload.Name)
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "incorrect name or password",
		})
		return
	}

	token, err := auth.GenerateToken(
		worker.UUID,
		worker.Name,
		worker.JobTitle,
		worker.Department_id,
		worker.Department_name,
		worker.AccessLevel,
	)

	if err != nil {
		log.Printf("error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "token generation error",
		})
		return
	}

	c.SetCookie(
		"auth_token",
		token,
		86400,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "auth success",
		"token":   token,
		"worker": gin.H{
			"id":              worker.UUID,
			"name":            worker.Name,
			"jobtitle":        worker.JobTitle,
			"department_id":   worker.Department_id,
			"department_name": worker.Department_name,
			"accesslevel":     worker.AccessLevel,
		},
	})
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
	var worker payloads.WorkerPayload

	if err := c.BindJSON(&worker); err != nil {
		log.Printf("error while binding worker payload: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad json body",
		})

		return
	}

	requestingUser, err := h.workerService.GetWorkerByName(c.GetString("worker_name"))

	if err != nil {
		log.Printf("error getting requesting user: %v\n", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})

		return
	}
	if requestingUser.AccessLevel < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":     "your request was rejected",
			"accesslevel": requestingUser.AccessLevel,
		})

		return
	}

	newWorker, err := h.workerService.CreateWorker(worker.Name, worker.JobTitle, requestingUser.Department_id, requestingUser.Department_name, worker.Password)
	if err != nil {
		log.Printf("error creating worker: %v", err)
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

	requestingUser, err := h.workerService.GetWorkerByName(c.GetString("worker_name"))
	if err != nil {
		log.Printf("error getting requesting user: %v\n", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})

		return
	}
	if requestingUser.AccessLevel < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":     "your request was rejected",
			"accesslevel": requestingUser.AccessLevel,
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
