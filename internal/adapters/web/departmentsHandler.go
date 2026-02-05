package web

import (
	"log"
	"net/http"
	"pi/internal/adapters/payloads"
	"pi/internal/app/interfaces/services"
	"pi/pkg/utils/middleware"

	"github.com/gin-gonic/gin"
)

type departmentHandler struct {
	departmentService services.DepartmentService
}

func NewDepartmentHandler(departmentService services.DepartmentService) Handler {
	return &departmentHandler{departmentService: departmentService}
}

func (h *departmentHandler) Register(router *gin.Engine) {
	departments := router.Group("/departments")
	departments.Use(middleware.AuthMiddleware())
	{
		departments.GET("/", h.ListAllDepartmentsHandler)
		departments.POST("/new", h.AddDepartmentHandler)
		departments.PUT("/:department")
		departments.DELETE("/:department", h.DeleteDepartmentHandler)
	}
}

func (h *departmentHandler) ListAllDepartmentsHandler(c *gin.Context) {
	departments, err := h.departmentService.GetAllDepartments()
	if err != nil {
		log.Printf("Failed to get all departments: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})

		return
	}

	c.JSON(http.StatusOK, departments)
}

func (h *departmentHandler) DeleteDepartmentHandler(c *gin.Context) {
	accessLevel, _ := c.Get("access_level")
	if level, ok := accessLevel.(int); !ok || level < 2 {
		c.JSON(http.StatusForbidden, gin.H{
			"message":     "доступ запрещён: требуется уровень 2 или выше",
			"accesslevel": accessLevel,
		})
		return
	}

	departmentName := c.Param("department")
	if len(departmentName) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})

		return
	}

	if err := h.departmentService.DeleteDepartmentByName(departmentName); err != nil {
		log.Printf("error deleting department %v: %v\n", departmentName, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "department was successfully deleted",
	})
}

func (h *departmentHandler) AddDepartmentHandler(c *gin.Context) {
	accessLevel, _ := c.Get("access_level")
	if level, ok := accessLevel.(int); !ok || level < 2 {
		c.JSON(http.StatusForbidden, gin.H{
			"message":     "доступ запрещён: требуется уровень 2 или выше",
			"accesslevel": accessLevel,
		})
		return
	}

	var departament payloads.DepartmentPayload

	if err := c.BindJSON(&departament); err != nil {
		log.Printf("error while binding department payload: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad json body",
		})

		return
	}

	newDepartment, err := h.departmentService.CreateDepartment(departament.Name)
	if err != nil {
		log.Printf("error creating department: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "new department was created",
		"department": newDepartment,
	})
}
