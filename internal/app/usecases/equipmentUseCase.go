package usecases

import (
	"context"
	"pi/internal/app/entities"
	repointerfaces "pi/internal/app/interfaces/repoInterfaces"
	"pi/internal/app/interfaces/services"
)

type equipmentService struct {
	repo repointerfaces.EquipmentRepo
}

func NewEquipmentService(repo repointerfaces.EquipmentRepo) services.EquipmentService {
	return &equipmentService {
		repo: repo,
	}
}

func (us *equipmentService) CreateEquipment(
	name, equipmentType, serialNum, inventoryNum, status, location string) (*entities.Equipment, error) {
		equipment := entities.NewEquipment(name, equipmentType, serialNum, inventoryNum, status, location)

		return equipment, us.repo.Create(context.Background(), equipment)
	}

	func (us *equipmentService) GetAllEquipments() ([]entities.Equipment, error) {
		return us.repo.GetAll(context.Background())
	}

	func (us *equipmentService) DeleteEquipmentByUUID (id string) error {
		return us.repo.DeleteByUUID(context.Background(), id)
	}
