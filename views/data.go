package views

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"

	// AlertMsgGeneric is displayed when any random error
	// is encountered by our backend
	AlertMsgGeneric = "Something went wrong. Please try again, and contact us if the problem persists!"
)

type Alert struct {
	Level   string
	Message string
}

type Data struct {
	Alert *Alert
	Body  interface{}
}
