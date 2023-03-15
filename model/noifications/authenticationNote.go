package noifications

var noteType string //nolint:gochecknoglobals

func NotifyLogin() (string, string) {
	noteType = "l"
	return "This is a notification to let you know that there has been a new login into your account.", noteType
}
