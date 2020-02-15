package bugsnagda

type Exception struct {
	ErrorClass string       `json:"errorClass"`
	Message    string       `json:"message"`
	Stacktrace []Stacktrace `json:"stacktrace"`
}
