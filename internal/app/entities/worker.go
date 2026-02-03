package entities

type Worker struct {
	UUID string
	Name string `json:"name"`
	JobTitle string `json:"jobtitle"`
	Department string `json:"department"`
	HashPass string `json:"password"`
	AcessLevel int `json:"acesslevel"`
}

func (w *Worker) NewWorker (name, jobTitle, department, hashPass string, accessLevel int) *Worker {
	return &Worker {
		Name: name,
		JobTitle: jobTitle,
		Department: department,
		HashPass: hashPass,
		AcessLevel: accessLevel,
	}
}
