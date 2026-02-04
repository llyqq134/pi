package entities

type Worker struct {
	UUID            string
	Name            string `json:"name" db:"name"`
	JobTitle        string `json:"jobtitle" db:"jobtitle"`
	Department_id   string `json:"department_id" db:"department"`
	Department_name string `json:"department_name" db:"department"`
	Password        string `json:"password" db:"password"`
	AccessLevel     int    `json:"acesslevel" db:"accesslevel"`
}

func NewWorker(name, jobTitle, department_id, department_name, password string, accessLevel int) *Worker {
	return &Worker{
		Name:            name,
		JobTitle:        jobTitle,
		Department_id:   department_id,
		Department_name: department_name,
		Password:        password,
		AccessLevel:     accessLevel,
	}
}
