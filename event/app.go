package event

type App struct {
	EventID      string
	ReleaseStage string `json:"releaseStage"`
}
