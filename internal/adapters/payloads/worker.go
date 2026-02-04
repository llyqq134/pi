package payloads

type WorkerPayload struct {
	Name       string `json:"name"`
	JobTitle   string `json:"jobtitle"`
	Department string `json:"department"`
	Password   string `json:"password"`
}
