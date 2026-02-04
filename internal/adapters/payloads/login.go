package payloads

type LoginPayload struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
