package web

import (
	"log"
	"net/http"
	"pi/internal/app/interfaces/services"

	"github.com/gin-gonic/gin"
)

type equipmentHandler struct {
	equipmentService services.EquipmentService
}

func NewEquipmentHandler (equipmentService services.EquipmentService) Handler {
	return &equipmentHandler{equipmentService: equipmentService}
}

func (h *equipmentHandler) Register(router *gin.Engine) {
	equipment := router.Group("/equipment")
	{
		equipment.GET("/list", h.ListAllEquipment)
	}
}

func (h *equipmentHandler)ListAllEquipment(c *gin.Context) {
	equipments, err := h.equipmentService.GetAllEquipments()
	if err != nil {
		log.Printf("error listing all equipment: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H {
			"message": "internal server error",
		})

		return
	}

	c.JSON(http.StatusOK, equipments)
}

