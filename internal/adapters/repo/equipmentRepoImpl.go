package repo

import (
	"context"
	"log"
	"pi/internal/app/entities"
	"pi/pkg/db"

	"github.com/jackc/pgx/v5/pgconn"
)

type EquipmentRepoImpl struct {
	Client db.Client
}

func NewEquipmentImpl(client db.Client) EquipmentRepoImpl {
	return EquipmentRepoImpl{Client: client}
}

func (r *EquipmentRepoImpl) Create(ctx context.Context, equipment *entities.Equipment) error {
	query := `
	INSERT INTO equipment (name, type, serial_number, inventory_number, status, location)
	RETURNING id
	`

	if err := r.Client.QueryRow(ctx, query,
		&equipment.Name, &equipment.Type, &equipment.SerialNumber, &equipment.InventoryNumber,
		&equipment.Status, &equipment.Location).Scan(&equipment.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			log.Printf("SQL Error: %s\nDetail: %s\nWhere: %s\nCode: %s\nSQL state: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())

			return err
		}

		return err

	}

	return nil
}

func (r *EquipmentRepoImpl) GetAll(ctx context.Context) ([]entities.Equipment, error) {
	query := `
	SELECT id, name, type, serial_number, inventory_number, status, location
	FROM equipment LIMIT 1000
	`
	rows, err := r.Client.Query(ctx, query)
	if err != nil {
		log.Printf("error selecting all equipment: %v\n", err)

		return nil, err
	}

	defer rows.Close()

	equipments := make([]entities.Equipment, 0)

	for rows.Next() {
		var equipment entities.Equipment
		if err := rows.Scan(
			&equipment.ID, &equipment.Name, &equipment.SerialNumber,
			&equipment.InventoryNumber, &equipment.Status, &equipment.Location); err != nil {
			log.Printf("error scanning equipment: %v\n", err)

			return nil, err
		}

		equipments = append(equipments, equipment)
	}

	return equipments, nil
}

func (r *EquipmentRepoImpl) DeleteByUUID(ctx context.Context, id string) error {
	query := `
	DELETE FROM equipment WHERE id = $1
	`

	if _, err := r.Client.Exec(ctx, query, id); err != nil {
		log.Printf("error deleting equipment: %v\n", err)

		return err
	}

	return nil
}
