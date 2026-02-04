package repo

import (
	"context"
	"log"
	"pi/internal/app/entities"
	"pi/pkg/db"

	"github.com/jackc/pgx/v5/pgconn"
)

type WorkerRepoImpl struct {
	Client db.Client
}

func NewWorkerImpl(client db.Client) WorkerRepoImpl {
	return WorkerRepoImpl{
		Client: client,
	}
}

func (r *WorkerRepoImpl) Create(ctx context.Context, worker *entities.Worker) error {
	query := `
	INSERT INTO
	workers (name, jobtitle, departament, password, accesslevel)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`

	if err := r.Client.QueryRow(ctx, query,
		&worker.Name, &worker.JobTitle, &worker.Department, &worker.Password, &worker.AcessLevel).
		Scan(&worker.UUID); err != nil {
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
	query := `
		SELECT id, name, jobtitle, department, password, accesslevel FROM workers WHERE id = $1
	`
	var worker entities.Worker

	if err := r.Client.QueryRow(ctx, query, uuid).Scan(
		&worker.UUID, &worker.Name, &worker.JobTitle, &worker.Password, &worker.AcessLevel); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			log.Printf("SQL Error: %s\nDetail: %s\nWhere: %s\nCode: %s\nSQL state: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())

			return entities.Worker{}, nil
		}
		return entities.Worker{}, err
	}

	return worker, nil
}

func (r *WorkerRepoImpl) GetAll(ctx context.Context) ([]entities.Worker, error) {
	query := `
	SELECT id, name, jobtitle, department, password, accesslevel FROM workers LIMIT 1000
	`
	rows, err := r.Client.Query(ctx, query)

	if err != nil {
		log.Printf("Error querying workers: %v\n", err)
		return nil, err
	}

	defer rows.Close()

	workers := make([]entities.Worker, 0)

	for rows.Next() {
		var worker entities.Worker
		if err := rows.Scan(&worker.UUID, &worker.Name, &worker.JobTitle, &worker.Department, &worker.AcessLevel); err != nil {
			log.Printf("err scanning workers: %v\n", err)
			return nil, err
		}
		workers = append(workers, worker)
	}

	return workers, err
}

func (r *WorkerRepoImpl) GetAllByDepartment(ctx context.Context, department string) ([]entities.Worker, error) {
	query := `
	SELECT id, name, jobtitle, department, password, accesslevel FROM workers WHERE department = $1
	`
	rows, err := r.Client.Query(ctx, query, department)

	if err != nil {
		log.Printf("error quering workers by department: %v\n", err)
		return nil, err
	}

	workers := make([]entities.Worker, 0)

	for rows.Next() {
		var worker entities.Worker
		if err := rows.Scan(
			&worker.UUID, &worker.Name, &worker.JobTitle, &worker.Department, &worker.Password, &worker.AcessLevel,
		); err != nil {
			log.Printf("error scanning workers: %v\n", err)
			return workers, err
		}

		workers = append(workers, worker)
	}

	return workers, err
}

func (r *WorkerRepoImpl) Update(ctx context.Context, worker *entities.Worker) error {
	query := `
		UPDATE workers SET name = $2, jobtitle = $3, department = $4, password = $5, accesslevel = $6
		WHERE id = $1
	`

	if _, err := r.Client.Exec(ctx, query, &worker.UUID, &worker.Name, &worker.JobTitle, &worker.Department, &worker.Password, &worker.AcessLevel); err != nil {
		log.Printf("error updating worker: %v\n", err)

		return err
	}

	return nil
}

func (r *WorkerRepoImpl) DeleteByUUID(ctx context.Context, uuid string) error {
	query := `
		DELETE FROM workers
		WHERE id = $1
	`

	if _, err := r.Client.Exec(ctx, query, uuid); err != nil {
		log.Printf("error deleting worker: %v\n", err)

		return err
	}

	return nil
}
