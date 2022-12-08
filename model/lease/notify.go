package lease

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

var SMTP_Host = os.Getenv("HOST")

// function to notify user once lease is satisfied.
func (r *ResourceReq) Notify() {
	satisfactionStatus := r.IsSatisfied()

	if satisfactionStatus == true {
		r.SendNotification()
	}
}

// function to email notification.
func (r *ResourceReq) SendNotification() {
	// probably need a custom email address for the zebra server.
	from_mail := "admin@zebra.project-safari.io"
	from_pass := "Riddikulus"

	// need a function to get specific's user's email address.
	recipient := []string{}
	recipient = append(recipient, "user@zebra.project-safari.io")

	// Message.
	message := []byte("Your lease request for resource " + r.Name + " of type " + r.Type + "has been satisfied. \nYou can now check it out on Zebra.") //nolint:lll

	// Authentication.
	auth := smtp.PlainAuth("", from_mail, from_pass, SMTP_Host)

	fmt.Println(auth)
	// Sending email.
	if err := smtp.SendMail(fmt.Sprintf("%s:%d", SMTP_Host, 587), auth, from_mail, recipient, message); err != nil {
		log.Printf("Error sending mail %v", err)
		return
	}

	log.Println("Email Sent Successfully!")
}
