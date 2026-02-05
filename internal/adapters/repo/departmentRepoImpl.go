package repo

import (
	"context"
	"log"
	"pi/internal/app/entities"
	"pi/pkg/db"

	"github.com/jackc/pgx/v5/pgconn"
)

type DepartmentRepoImpl struct {
	client db.Client
}

func NewDepartmentImpl(client db.Client) DepartmentRepoImpl {
	return DepartmentRepoImpl{client: client}
}

func (r *DepartmentRepoImpl) Create(ctx context.Context, departament *entities.Department) error {
	query := `
	INSERT INTO departments (name) values ($1)
	`

	if err := r.client.QueryRow(ctx, query, &departament.Name).Scan(&departament.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			log.Printf("SQL Error: %s\nDetail: %s\nWhere: %s\nCode: %s\nSQL state: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())

			return err
		}

		return err
	}

	return nil
}

func (r *DepartmentRepoImpl) GetAll(ctx context.Context) ([]entities.Department, error) {
	query := `
	SELECT id, name FROM departments LIMIT 1000
	`
	rows, err := r.client.Query(ctx, query)
	if err != nil {
		log.Printf("erorr querying departments: %v\n", err)
		return nil, err
	}

	defer rows.Close()

	departaments := make([]entities.Department, 0)

	for rows.Next() {
		var departament entities.Department
		if err := rows.Scan(&departament.ID, &departament.Name); err != nil {
			log.Printf("error scanning departments: %v\n", err)

			return nil, err
		}

		departaments = append(departaments, departament)
	}

	return departaments, nil
}

func (r *DepartmentRepoImpl) Update(ctx context.Context, departament *entities.Department) error {
	query := `
	UPDATE departments SET name = $2 WHERE id = $1
	`

	if _, err := r.client.Exec(ctx, query, &departament.ID, &departament.Name); err != nil {
		log.Printf("error updating department: %v\n", err)
		return err
	}

	return nil
}

func (r *DepartmentRepoImpl) DeleteByName(ctx context.Context, name string) error {
	query := `
	DELETE FROM departments WHERE name = $1
	`

	if _, err := r.client.Exec(ctx, query, name); err != nil {
		log.Printf("error deleting department: %v\n", err)
		return err
	}

	return nil
}
