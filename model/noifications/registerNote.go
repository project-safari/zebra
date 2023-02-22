package noifications

var noteType string //nolint:gochecknoglobals

func NotifyRegister() (string, string) {
	noteType = "r"
	return "This is a notification to let you know that a new account was created for you with this email address.", noteType
}
