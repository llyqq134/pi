package repointerfaces

import (
	"context"
	"pi/internal/app/entities"
	"time"
)

type RecordsRepo interface {
	Create(ctx context.Context, records *entities.Records) error
	GetByUUID(ctx context.Context, uuid string) (entities.Records, error)
	GetRecordsUpTo(ctx context.Context, startDate, endDate time.Time) ([]entities.Records, error)
	GetAll(ctx context.Context) ([]entities.Records, error)
	DeleteByUUID(ctx context.Context, uuid string) error
}
