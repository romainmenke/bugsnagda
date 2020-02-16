package event

import "time"

type Breadcrumb struct {
	EventID   string
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	MetaData  interface{} `json:"metaData"`
}
