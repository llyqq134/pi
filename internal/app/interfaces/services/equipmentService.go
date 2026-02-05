package services

import "pi/internal/app/entities"

type EquipmentService interface {
	CreateEquipment(
		name, equipmentType, serialNum, inventoryNum, status, location string) (*entities.Equipment, error)
	GetAllEquipments() ([]entities.Equipment, error)
	DeleteEquipmentByUUID(id string) error
}
