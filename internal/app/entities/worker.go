package entities

type Worker struct {
	UUID       string
	Name       string `json:"name" db:"name"`
	JobTitle   string `json:"jobtitle" db:"jobtitle"`
	Department string `json:"department" db:"department"`
	Password   string `json:"password" db:"password"`
	AcessLevel int    `json:"acesslevel" db:"accesslevel"`
}

func NewWorker(name, jobTitle, department, password string, accessLevel int) *Worker {
	return &Worker{
		Name:       name,
		JobTitle:   jobTitle,
		Department: department,
		Password:   password,
		AcessLevel: accessLevel,
	}
}
