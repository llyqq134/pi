package repo

import (
	"context"
	"log"
	"pi/internal/app/entities"
	"pi/internal/app/interfaces/repoInterfaces"
	"pi/pkg/db"

	"github.com/jackc/pgx/v5/pgconn"
)

type WorkerRepoImpl struct {
	Client db.Client
}

func NewWorkerImpl(client db.Client) repointerfaces.WorkerRepo {
	return WorkerRepoImpl{
	Client: client,
	}
}

func (r *WorkerRepoImpl) Create (ctx context.Context, worker *entities.Worker) error {
	query := `
	INSERT INTO
	workers (name, jobtitle, departament, hashpass, accesslevel)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`

	if err := r.Client.QueryRow(ctx, query, 
		&worker.Name, &worker.JobTitle, &worker.Departament, &worker.HashPass, &worker.AcessLevel).Scan(&worker.UUID);
		err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				log.Printf("SQL Error: %s\nDetail: %s\nWhere: %s\nCode: %s\nSQL state: %s",
					pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())

					return nil
			}

			return err
		}

		return nil
}

func (r *WorkerRepoImpl) GetByUUID(ctx context.Context, uuid string) (entities.Worker, error) {
	
} 
