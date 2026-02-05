package entities

import "time"

type Records struct {
	ID                 string    `json:"id"`
	EquipmentId        string    `json:"equipment_id"`
	EquipmentName      string    `json:"equipment_name"` // human-readable name for report
	WorkerId           string    `json:"worker_id"`
	WorkerName         string    `json:"worker_name"`
	DepartmentID       string    `json:"department_id"`
	DepartmentName     string    `json:"department_name"`
	IssuedAt           time.Time `json:"issued_at"`
	ReturnedAt         time.Time `json:"returned_at"`
	ExpectedReturnDate time.Time `json:"expected_returned_date"`
	Status             string    `json:"status"`
}

func NewRecords(equipmentId, workerId, workerName, departmentId, departmentName string, ExpectedReturnDate time.Time, status string) *Records {
	return &Records{
		EquipmentId:        equipmentId,
		WorkerId:           workerId,
		WorkerName:         workerName,
		DepartmentID:       departmentId,
		DepartmentName:     departmentName,
		ExpectedReturnDate: ExpectedReturnDate,
		Status:             status,
	}
}
