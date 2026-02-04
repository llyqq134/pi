package usecases

import (
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
	if jobTitle == 
	worker := entities.NewWorker()
}
func (s *workerService) GetWokerByUUID(uuid string) (*entities.Worker, error)
func (s *workerService) GetAllWorkers() ([]entities.Worker, error)
func (s *workerService) GetAllWorkersByDepartment(department string) ([]entities.Worker, error)
