package lease

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"sync"
)

//nolint:gochecknoglobals
var (
	FromMail     = "admin@zebra.project-safari.io"
	Mailpassword = "Riddikulus"
	SMTPhost     = os.Getenv("HOST")
	MailSubject  string
	MailBody     string
	to           *mail.Address
	from         *mail.Address
	auth         smtp.Auth
	tlsconfig    *tls.Config
)

type Container struct {
	m       sync.Mutex
	Headers map[string]string
}

func NewContainer() *Container {
	return &Container{
		Headers: make(map[string]string),
	}
}

func (l *Lease) GetEmail() string {
	// return "user@zebra.project-safari.io"
	return l.Status.UsedBy
}

// function to notify user once lease is satisfied.
func (l *Lease) Notify() {
	strOne := "This is a notification to let you know that your lease request has been satisfied."
	strThree := "for resource "
	strTwo := " is satisfied.\nLog back in to check it out!"

	message := strOne

	for _, r := range l.Request {
		message += strTwo + r.Name + " " + r.Type + strThree
		user := l.GetEmail()

		if r.IsSatisfied() {
			l.SendNotification("Zebra Lease Request Satisfied", message, user)
		}
	}
}

func (l *Lease) NotifyActive() {
	if l != nil {
		strOne := "This is a notification to let you know that your lease request "
		strTwo := "has been activated.\nCheck back later to see if it's satisfited."

		message := strOne + strTwo
		// user := l.GetEmail()

		l.SendNotification("Zebra Lease Request Placed", message, "")
	}
}

// function to email notification.
//
//nolint:gomnd, funlen, cyclop
func (l *Lease) SendNotification(subject, msg string, recipient string) {
	to := mail.Address{Name: "", Address: recipient}

	MailSubject = subject
	MailBody = msg

	container := setHeaders()

	// Setup message
	message := ""
	for k, v := range container.Headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += "\r\n" + MailBody

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", SMTPhost, 465), tlsconfig)
	if err != nil {
		log.Printf("Error sending mail %v", err)
	}

	c, err := smtp.NewClient(conn, SMTPhost)
	if err != nil {
		log.Printf("Error sending mail %v", err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Printf("Error sending mail %v", err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Printf("Error sending mail %v", err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Printf("Error sending mail %v", err)
	}

	if err != nil {
		return
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Printf("Error sending mail %v", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Printf("Error sending mail %v", err)
	}

	err = w.Close()
	if err != nil {
		log.Printf("Error sending mail %v", err)
	}

	if err = c.Quit(); err != nil {
		return
	}

	if err != nil {
		return
	}
}
