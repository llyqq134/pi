package repo

import (
	"context"
	"database/sql"
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
	INSERT INTO equipment_records (equipment_id, worker_id, worker_name, department_id, department_name, issued_at, returned_at, expected_return_date, status)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id
	`
	var returnedAt interface{}
	if record.ReturnedAt.IsZero() {
		returnedAt = nil // NULL for "not yet returned"
	} else {
		returnedAt = &record.ReturnedAt
	}

	if err := r.Client.QueryRow(
		ctx, query, &record.EquipmentId, &record.WorkerId, &record.WorkerName, &record.DepartmentID, &record.DepartmentName, &record.IssuedAt, returnedAt, &record.ExpectedReturnDate, &record.Status).
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
	SELECT id, equipment_id, worker_id, worker_name, department_id, department_name, issued_at, returned_at, expected_return_date, status 
	FROM equipment_records WHERE id = $1
	`
	var record entities.Records
	var returnedAt sql.NullTime

	if err := r.Client.QueryRow(ctx, query, uuid).Scan(
		&record.ID, &record.EquipmentId, &record.WorkerId, &record.WorkerName, &record.DepartmentID, &record.DepartmentName,
		&record.IssuedAt, &returnedAt, &record.ExpectedReturnDate, &record.Status); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			log.Printf("SQL Error: %s\nDetail: %s\nWhere: %s\nCode: %s\nSQL state: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())

			return entities.Records{}, err
		}

		return entities.Records{}, err
	}
	if returnedAt.Valid {
		record.ReturnedAt = returnedAt.Time
	}

	return record, nil
}

func (r *RecordsRepoImpl) GetRecordsUpTo(ctx context.Context, startDate, endDate time.Time) ([]entities.Records, error) {
	query := `
	SELECT r.id, r.equipment_id, COALESCE(e.name, r.equipment_id::text) as equipment_name,
	       r.worker_id, r.worker_name, r.department_id, r.department_name, r.issued_at, r.returned_at, 
	       r.expected_return_date, r.status
	FROM equipment_records r
	LEFT JOIN equipment e ON r.equipment_id = e.id
	WHERE r.issued_at >= $1 AND r.issued_at <= $2
	ORDER BY r.issued_at DESC
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
		var returnedAt sql.NullTime
		if err := rows.Scan(
			&record.ID, &record.EquipmentId, &record.EquipmentName, &record.WorkerId, &record.WorkerName, &record.DepartmentID, &record.DepartmentName, &record.IssuedAt, &returnedAt, &record.ExpectedReturnDate, &record.Status); err != nil {
			log.Printf("error scanning record: %v\n", err)

			return nil, err
		}
		if returnedAt.Valid {
			record.ReturnedAt = returnedAt.Time
		}

		records = append(records, record)
	}

	return records, nil
}

func (r *RecordsRepoImpl) GetAll(ctx context.Context) ([]entities.Records, error) {
	query := `
	SELECT id, equipment_id, worker_id, worker_name, department_id, department_name, issued_at, returned_at, expected_return_date, status
	FROM equipment_records LIMIT 1000
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
		var returnedAt sql.NullTime
		if err := rows.Scan(
			&record.ID, &record.EquipmentId, &record.WorkerId, &record.WorkerName, &record.DepartmentID, &record.DepartmentName, &record.IssuedAt, &returnedAt, &record.ExpectedReturnDate, &record.Status); err != nil {
			log.Printf("error scanning record: %v\n", err)
			return nil, err
		}
		if returnedAt.Valid {
			record.ReturnedAt = returnedAt.Time
		}

		records = append(records, record)
	}

	return records, nil
}

func (r *RecordsRepoImpl) DeleteByUUID(ctx context.Context, uuid string) error {
	query := `
	DELETE FROM equipment_records WHERE id = $1
	`

	if _, err := r.Client.Exec(ctx, query, uuid); err != nil {
		log.Printf("error deleting record: %v\n", err)
		return err
	}

	return nil
}
