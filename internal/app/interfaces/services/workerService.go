package services

import "pi/internal/app/entities"

type WorkerService interface {
	CreateWorker(name, jobTitle, departament, password string) (*entities.Worker, error)
	GetWokerByUUID(uuid string) (*entities.Worker, error)
	GetAllWorkers() ([]entities.Worker, error)
	GetAllWorkersByDepartment(department string) ([]entities.Worker, error)
}
