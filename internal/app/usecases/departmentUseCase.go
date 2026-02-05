package usecases

import (
	"context"
	"pi/internal/app/entities"
	repointerfaces "pi/internal/app/interfaces/repoInterfaces"
	"pi/internal/app/interfaces/services"
)

type departamentService struct {
	repo       repointerfaces.DepartemntRepo
	workerRepo repointerfaces.WorkerRepo
}

func NewDepartmentService(repo repointerfaces.DepartemntRepo, workerRepo repointerfaces.WorkerRepo) services.DepartmentService {
	return &departamentService{repo: repo, workerRepo: workerRepo}
}

func (s *departamentService) CreateDepartment(name string) (*entities.Department, error) {
	department := entities.NewDepartment(name)
	return department, s.repo.Create(context.Background(), department)
}

func (s *departamentService) GetAllDepartments() ([]entities.Department, error) {
	return s.repo.GetAll(context.Background())
}

func (s *departamentService) DeleteDepartmentByName(name string) error {
	// Каскадное удаление: сначала удаляем сотрудников отдела
	if err := s.workerRepo.DeleteByDepartmentName(context.Background(), name); err != nil {
		return err
	}
	return s.repo.DeleteByName(context.Background(), name)
}
