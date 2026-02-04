package payloads

type WorkerPayload struct {
	Name     string `json:"name"`
	JobTitle string `json:"jobtitle"`
	Password string `json:"password"`
}
