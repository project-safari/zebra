package noifications

var Notified bool //nolint:gochecknoglobals

func NotifyLogin() (string, bool) {
	Notified = true

	return "This is a notification to let you know that there has been a new login into your account.", Notified
}
