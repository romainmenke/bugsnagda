package bugsnagda

import "time"

type Breadcrumb struct {
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	MetaData  interface{} `json:"metaData"`
}
