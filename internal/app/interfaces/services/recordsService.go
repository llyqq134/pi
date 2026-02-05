package services

import (
	"pi/internal/app/entities"
	"time"
)

type RecordService interface {
	CreateRecord(equipmentId, workerId, workerName, departemntId, departmentName string, expectedReturnDate time.Time, status string) (*entities.Records, error)
	GetRecordByUUID(uuid string) (entities.Records, error)
	GetRecordsUpTo(startDate, endDate time.Time) ([]entities.Records, error)
	GetAllRecords() ([]entities.Records, error)
	DeleteRecordByUUID(uuid string) error
}
