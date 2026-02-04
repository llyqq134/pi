package services

import "pi/internal/app/entities"

type DepartmentService interface {
	CreateDepartment(name string) (*entities.Department, error)
	GetAllDepartments() ([]entities.Department, error)
	DeleteDepartmentByName(name string) error
}
