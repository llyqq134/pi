package usecases

import (
	"context"
	"pi/internal/app/entities"
	repointerfaces "pi/internal/app/interfaces/repoInterfaces"
	"pi/internal/app/interfaces/services"
)

type workerService struct {
	repo repointerfaces.WorkerRepo
}

func NewWorkerService(repo repointerfaces.WorkerRepo) services.WorkerService {
	return &workerService{repo: repo}
}

func (s *workerService) CreateWorker(name, jobTitle, department_id, department_name, password string) (*entities.Worker, error) {
	var accesslevel int

	switch jobTitle {
	case "manager":
		accesslevel = 2
	case "admin":
		accesslevel = 3
	default:
		accesslevel = 1
	}

	worker := entities.NewWorker(name, jobTitle, department_id, department_name, password, accesslevel)

	if err := s.repo.Create(context.Background(), worker); err != nil {
		return &entities.Worker{}, err
	}

	return worker, nil
}

func (s *workerService) GetWokerByUUID(uuid string) (entities.Worker, error) {
	return s.repo.GetByUUID(context.Background(), uuid)
}

func (s *workerService) GetWorkerByName(name string) (entities.Worker, error) {
	return s.repo.GetByName(context.Background(), name)
}

func (s *workerService) GetAllWorkers() ([]entities.Worker, error) {
	return s.repo.GetAll(context.Background())
}

func (s *workerService) GetAllWorkersByDepartment(department string) ([]entities.Worker, error) {
	return s.repo.GetAllByDepartment(context.Background(), department)
}

func (s *workerService) DeleteWorkerByUUID(uuid string) error {
	return s.repo.DeleteByUUID(context.Background(), uuid)
}
