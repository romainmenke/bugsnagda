package errport

type ReopenRules struct {
	ReportID              string
	ReopenIf              string `json:"reopen_if"`
	ReopenAfter           string `json:"reopen_after"`
	Seconds               int    `json:"seconds"`
	Occurences            int    `json:"occurences"`
	Hours                 int    `json:"hours"`
	OccurrenceThreshold   int    `json:"occurrence_threshold"`
	AdditionalOccurrences int    `json:"additional_occurrences"`
}
