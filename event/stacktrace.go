package event

type Stacktrace struct {
	ID                string
	EventID           string
	ThreadID          int
	ExceptionID       string
	Chain             int
	LineNumber        int               `json:"lineNumber"`
	ColumnNumber      int               `json:"columnNumber"`
	File              string            `json:"file"`
	InProject         bool              `json:"inProject"`
	Method            string            `json:"method"`
	Code              []*StacktraceCode `json:"code"`
	MachoUUID         string            `json:"machoUUID"`
	SourceControlLink string            `json:"sourceControlLink"`
	SourceControlName string            `json:"sourceControlName"`
}
