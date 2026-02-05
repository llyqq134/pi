package repo

import (
	"context"
	"log"
	"pi/internal/app/entities"
	"pi/pkg/db"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type RecordsRepoImpl struct {
	Client db.Client
}

func NewRecordsImpl(client db.Client) RecordsRepoImpl {
	return RecordsRepoImpl{Client: client}
}

func (r *RecordsRepoImpl) Create(ctx context.Context, record *entities.Records) error {
	query := `
	INSERT INTO records (equipment_id, worker_id, department_id, issued_at, returned_at, expected_return_date, status)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`

	if err := r.Client.QueryRow(
		ctx, query, &record.EquipmentId, &record.WorkerId, &record.DepartmentID, &record.IssuedAt, &record.ReturnedAt, &record.ExpectedReturnDate, &record.Status).
		Scan(&record.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			log.Printf("SQL Error: %s\nDetail: %s\nWhere: %s\nCode: %s\nSQL state: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())

			return err
		}

		return err
	}
	return nil
}

func (r *RecordsRepoImpl) GetByUUID(ctx context.Context, uuid string) (entities.Records, error) {
	query := `
	SELECT id, equipment_id, worker_id, department_id, issued_at, returned_at, expected_return_date, status 
	FROM records WHERE id = $1
	`
	var record entities.Records

	if err := r.Client.QueryRow(
		ctx, query, &record.ID, &record.EquipmentId, &record.WorkerId, &record.DepartmentID, &record.IssuedAt, &record.ReturnedAt, &record.ExpectedReturnDate, &record.Status).
		Scan(&record.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			log.Printf("SQL Error: %s\nDetail: %s\nWhere: %s\nCode: %s\nSQL state: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())

			return entities.Records{}, err
		}

		return entities.Records{}, err
	}

	return record, nil
}

func (r *RecordsRepoImpl) GetRecordsUpTo(ctx context.Context, startDate, endDate time.Time) ([]entities.Records, error) {
	query := `
	SELECT id, equipment_id, worker_id, department_id, issued_at, returned_at, 
               expected_return_date, status
    FROM records
    WHERE issued_at >= $1 and issued_at < $2
    ORDER BY issued_at DESC
	`

	rows, err := r.Client.Query(ctx, query, startDate, endDate)
	if err != nil {
		log.Printf("error getting records in time stamp: %v\n", err)
		return nil, err
	}

	defer rows.Close()

	records := make([]entities.Records, 0)
	for rows.Next() {
		var record entities.Records
		if err := rows.Scan(
			&record.ID, &record.EquipmentId, &record.WorkerId, &record.DepartmentID, &record.IssuedAt, &record.ReturnedAt, &record.ExpectedReturnDate, &record.Status); err != nil {
			log.Printf("error scanning record: %v\n", err)

			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func (r *RecordsRepoImpl) GetAll(ctx context.Context) ([]entities.Records, error) {
	query := `
	SELECT equipment_id, worker_id, department_id, issued_at, returned_at, expected_return_date, status
	FROM records LIMIT 1000
	`

	rows, err := r.Client.Query(ctx, query)
	if err != nil {
		log.Printf("error querying all records: %v\n", err)

		return nil, err
	}

	defer rows.Close()

	records := make([]entities.Records, 0)
	for rows.Next() {
		var record entities.Records
		if err := rows.Scan(
			&record.ID, &record.EquipmentId, &record.WorkerId, &record.DepartmentID, &record.IssuedAt, &record.ReturnedAt, &record.ExpectedReturnDate, &record.Status); err != nil {
			log.Printf("error scanning record: %v\n", err)
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func (r *RecordsRepoImpl) DeleteByUUID(ctx context.Context, uuid string) error {
	query := `
	DELETE FROM records WHERE id = $1
	`

	if _, err := r.Client.Exec(ctx, query, uuid); err != nil {
		log.Printf("error deleting record: %v\n", err)
		return err
	}

	return nil
}
