package event

type Thread struct {
	EventID              string
	ID                   int           `json:"id"`
	Name                 string        `json:"name"`
	ErrorReportingThread bool          `json:"error_reporting_thread"`
	Stacktrace           []*Stacktrace `json:"stacktrace"`
}