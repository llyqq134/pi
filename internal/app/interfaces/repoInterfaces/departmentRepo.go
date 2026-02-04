package repointerfaces

import (
	"context"
	"pi/internal/app/entities"
)

type DepartemntRepo interface {
	Create(ctx context.Context, department *entities.Department) error
	GetAll(ctx context.Context) ([]entities.Department, error)
	Update(ctx context.Context, department *entities.Department) error
	DeleteByName(ctx context.Context, name string) error
}
