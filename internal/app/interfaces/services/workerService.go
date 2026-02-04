package services

import "pi/internal/app/entities"

type WorkerService interface {
	CreateWorker(name, jobTitle, department_id, department_name, password string) (*entities.Worker, error)
	GetWokerByUUID(uuid string) (entities.Worker, error)
	GetWorkerByName(name string) (entities.Worker, error)
	GetAllWorkers() ([]entities.Worker, error)
	GetAllWorkersByDepartment(department string) ([]entities.Worker, error)
	DeleteWorkerByUUID(uuid string) error
}
