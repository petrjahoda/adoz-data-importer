package main

import (
	"gopkg.in/gomail.v2"
	"os"
)

func SendMail(subject string, message string) {
	name, err := os.Hostname()
	if err != nil {
		LogError("MAIN", "Problem getting name of the computer, "+err.Error())
		name = ""
	}
	m := gomail.NewMessage()
	m.SetHeader("From", "jahoda@zapsi.eu")
	m.SetHeader("To", "jahoda@zapsi.eu")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", name+": "+message)
	d := gomail.NewDialer("smtp.forpsi.com", 587, "jahoda@zapsi.eu", "password") // PETRzpsMAIL79..
	if emailSentError := d.DialAndSend(m); emailSentError != nil {
		LogError("MAIN", "Email not sent: "+emailSentError.Error())
	} else {
		LogInfo("MAIN", "Email sent: "+subject)
	}
}
