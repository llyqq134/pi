package repointerfaces

import (
	"context"
	"pi/internal/app/entities"
)

type WorkerRepo interface {
	Create(ctx context.Context, worker *entities.Worker) error
	GetByUUID(ctx context.Context, uuid string) (entities.Worker, error)
	GetByName(ctx context.Context, name string) (entities.Worker, error)
	GetAll(ctx context.Context) ([]entities.Worker, error)
	GetAllByDepartment(ctx context.Context, department string) ([]entities.Worker, error)
	Update(ctx context.Context, worker *entities.Worker) error
	DeleteByUUID(ctx context.Context, uuid string) error
	DeleteByDepartmentName(ctx context.Context, departmentName string) error
}
