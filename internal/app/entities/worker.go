package entities

type Worker struct {
	UUID string
	Name string `json:"name"`
	JobTitle string `json:"jobtitle"`
	Departament string `json:"departament"`
	HashPass string `json:"password"`
	AcessLevel int `json:"acesslevel"`
}

func (w *Worker) NewWorker (name, jobTitle, departament, hashPass string, accessLevel int) *Worker {
	return &Worker {
		Name: name,
		JobTitle: jobTitle,
		Departament: departament,
		HashPass: hashPass,
		AcessLevel: accessLevel,
	}
}
