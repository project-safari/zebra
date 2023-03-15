package noifications

func NotifyRegister() (string, string) {
	NoteType = "r"

	return "This is a notification to let you know that a new account was created for you with this email address.", NoteType
}
