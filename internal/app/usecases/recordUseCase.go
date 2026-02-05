package usecases

import (
	"context"
	"pi/internal/app/entities"
	repointerfaces "pi/internal/app/interfaces/repoInterfaces"
	"pi/internal/app/interfaces/services"
	"time"
)

type recordService struct {
	repo repointerfaces.RecordsRepo
}

func NewRecordService(repo repointerfaces.RecordsRepo) services.RecordService {
	return &recordService{repo: repo}
}

func (us *recordService) CreateRecord(
	equipmentId, workerId, workerName, departmentId, departmentName string, expectedReturnDate time.Time, status string) (*entities.Records, error) {
	record := entities.NewRecords(equipmentId, workerId, workerName, departmentId, departmentName, expectedReturnDate, status)
	return record, us.repo.Create(context.Background(), record)
}

func (us *recordService) GetRecordByUUID(uuid string) (entities.Records, error) {
	return us.repo.GetByUUID(context.Background(), uuid)
}

func (us *recordService) GetRecordsUpTo(startDate, endDate time.Time) ([]entities.Records, error) {
	return us.repo.GetRecordsUpTo(context.Background(), startDate, endDate)
}

func (us *recordService) GetAllRecords() ([]entities.Records, error) {
	return us.repo.GetAll(context.Background())
}

func (us *recordService) DeleteRecordByUUID(uuid string) error {
	return us.repo.DeleteByUUID(context.Background(), uuid)
}
