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

	from      *mail.Address
	auth      smtp.Auth
	tlsconfig *tls.Config
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
	satisfactionStatus := l.IsSatisfied()

	strOne := "This is a notification to let you know that your lease request for resource "
	strTwo := " is satisfied.\nLog back in to check it out!"

	for each := 0; each < len(l.Request); each++ {
		message := strOne + l.Request[each].Name + " " + l.Request[each].Type + strTwo
		user := l.GetEmail()

		if satisfactionStatus {
			l.Request[each].SendNotification("Zebra Lease Request Satisfied", message, user)
		}
	}
}

func (l *Lease) NotifyActive() {
	strOne := "This is a notification to let you know that your lease request for resource "
	strTwo := " has been activated.\nCheck back later to see if it's satisfited."

	for each := 0; each < len(l.Request); each++ {
		message := strOne + l.Request[each].Name + " " + l.Request[each].Type + strTwo
		user := l.GetEmail()

		l.Request[each].SendNotification("Zebra Lease Request Placed", message, user)
	}
}

// function to email notification.
//
//nolint:gomnd
func (r *ResourceReq) SendNotification(subject, msg string, recipient string) {
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
}
