package event

type Exception struct {
	ID         string
	EventID    string
	ErrorClass string        `json:"errorClass"`
	Message    string        `json:"message"`
	Stacktrace []*Stacktrace `json:"stacktrace"`
}
