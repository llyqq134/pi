package payloads

type EquipmentPayload struct {
	Name string `json:"name"`
	Type string `json:"type"`
	SerialNumber string `json:"serialnumber"`
	InventoryNumber string `json:"inventorynumber"`
	Status string `json:"status"`
	Location string `json:"location"`
}
