package repointerfaces

import (
	"context"
	"pi/internal/app/entities"
)

type EquipmentRepo interface {
	Create (ctx context.Context, equipment *entities.Equipment) error 
	GetAll(ctx context.Context) ([]entities.Equipment, error)
	DeleteByUUID(ctx context.Context, id string) error
}

