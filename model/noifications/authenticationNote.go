package noifications

var NoteType string //nolint:gochecknoglobals

func NotifyLogin() (string, string) {
	NoteType = "l"
	return "This is a notification to let you know that there has been a new login into your account.", NoteType
}
