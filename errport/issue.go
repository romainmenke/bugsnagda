package errport

type Issue struct {
	ReportID string
	IssueID  string `json:"id"`
	Key      string `json:"key"`
	Number   int    `json:"number"`
	Type     string `json:"type"`
	URL      string `json:"url"`
}
