package bugsnagda

type Stacktrace struct {
	LineNumber        int            `json:"lineNumber"`
	ColumnNumber      int            `json:"columnNumber"`
	File              string         `json:"file"`
	InProject         bool           `json:"inProject"`
	Method            string         `json:"method"`
	Code              map[int]string `json:"code"`
	MachoUUID         string         `json:"machoUUID"`
	SourceControlLink string         `json:"sourceControlLink"`
	SourceControlName string         `json:"sourceControlName"`
}
