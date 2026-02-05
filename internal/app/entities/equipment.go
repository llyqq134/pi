package entities

type Equipment struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	SerialNumber string `json:"serialnumber"`
	InventoryNumber string `json:"inventorynumber"`
	Status string `json:"status"`
	Location string `json:"location"`
}

func NewEquipment(name, equipmentType, serialNumber, inventoryNumber, status, location string) *Equipment {
	return &Equipment{
		Name: name,
		Type: equipmentType,
		SerialNumber: serialNumber,
		InventoryNumber: inventoryNumber,
		Status: status,
		Location: location,
	}
}
