package noifications

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"sync"
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

//nolint:gochecknoglobals
var (
	FromMail     = "admin@zebra.project-safari.io"
	Mailpassword = "Riddikulus"
	SMTPhost     = os.Getenv("HOST")
	MailSubject  string
	MailBody     string

	from      *mail.Address
	auth      smtp.Auth
	tlsconfig *tls.Config
)

type NoteActions struct {
	Notified bool
	NoteType string
}

func (note *NoteActions) Sent() {
	note.Notified = true
}

func (note *NoteActions) Type(this string) {
	note.NoteType = this
}

// nolint: funlen, gomnd
func SendAccountNotification(subject, msg string, recipient string, kind string) {
	notif := new(NoteActions)
	to := mail.Address{Name: "", Address: recipient}

	MailSubject = subject
	MailBody = msg

	// initialize new container object
	container := NewContainer()
	// call mutex.lock to avoid multiple writes to
	// one header instance from running goroutines
	container.m.Lock()
	container.Headers["From"] = from.String()
	container.Headers["To"] = to.String()
	container.Headers["Subject"] = MailSubject
	// unlock mutex after function returns
	defer container.m.Unlock()

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

	notif.Sent()
	notif.Type(kind)
}
