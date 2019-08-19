package views

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"
)

type Alert struct {
	Level   string
	Message string
}

type Data struct {
	Alert *Alert
	Body  interface{}
}
