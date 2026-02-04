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

func (s *workerService) CreateWorker(name, jobTitle, departament, password string) (*entities.Worker, error) {
	var accesslevel int

	switch jobTitle {
	case "manager":
		accesslevel = 2
	case "admin":
		accesslevel = 3
	default:
		accesslevel = 1
	}

	worker := entities.NewWorker(name, jobTitle, departament, password, accesslevel)

	return worker, s.repo.Create(context.Background(), worker)
}

func (s *workerService) GetWokerByUUID(uuid string) (entities.Worker, error) {
	return s.repo.GetByUUID(context.Background(), uuid)
}

func (s *workerService) GetAllWorkers() ([]entities.Worker, error) {
	return s.repo.GetAll(context.Background())
}

func (s *workerService) GetAllWorkersByDepartment(department string) ([]entities.Worker, error) {
	return s.repo.GetAllByDepartment(context.Background(), department)
}
