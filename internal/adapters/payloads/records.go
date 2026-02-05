package payloads

type RecordsPayload struct {
	EquipmentId         string `json:"equipment_id"`
	ExpectedReturnDate  string `json:"expected_return_date"` // формат: "2006-01-02" (YYYY-MM-DD)
	Status              string `json:"status"`
	DepartmentId        string `json:"department_id"`   // целевой отдел (опционально)
	DepartmentName      string `json:"department_name"` // целевой отдел (опционально)
	WorkerId            string `json:"worker_id"`       // сотрудник-получатель (опционально)
	WorkerName          string `json:"worker_name"`     // сотрудник-получатель (опционально)
}
