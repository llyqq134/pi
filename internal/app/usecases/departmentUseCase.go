package usecases

import (
	"context"
	"pi/internal/app/entities"
	repointerfaces "pi/internal/app/interfaces/repoInterfaces"
	"pi/internal/app/interfaces/services"
)

type departamentService struct {
	repo repointerfaces.DepartemntRepo
}

func NewDepartmentService(repo repointerfaces.DepartemntRepo) services.DepartmentService {
	return &departamentService{repo: repo}
}

func (s *departamentService) CreateDepartment(name string) (*entities.Department, error) {
	department := entities.NewDepartment(name)
	return department, s.repo.Create(context.Background(), department)
}

func (s *departamentService) GetAllDepartments() ([]entities.Department, error) {
	return s.repo.GetAll(context.Background())
}

func (s *departamentService) DeleteDepartmentByName(name string) error {
	return s.repo.DeleteByName(context.Background(), name)
}
