package event

type Thread struct {
	EventID              string
	NumericID            int           `json:"id"`
	Name                 string        `json:"name"`
	ErrorReportingThread bool          `json:"error_reporting_thread"`
	Stacktrace           []*Stacktrace `json:"stacktrace"`
}
