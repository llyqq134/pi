package entities

type Department struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func NewDepartment(name string) *Department {
	return &Department{Name: name}
}
