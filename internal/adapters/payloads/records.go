package payloads

import "time"

type RecordsPayload struct {
	EquipmentId        string    `json:"equipment_id"`
	ExpectedReturnDate time.Time `json:"expected_return_date"`
	Status             string    `json:"status"`
}
