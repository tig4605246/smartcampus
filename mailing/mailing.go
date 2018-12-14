package mailing

import (
	"gopkg.in/gomail.v2"
)

const (
	kevin = "tig4605246@gmail.com"
)

//SendMail : Wrap a send-mail function for error report
func SendMail(content string, gwName string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "ti4605246@gmail.com")
	m.SetHeader("To", kevin)
	//m.SetAddressHeader("Cc", "avbee.lab@gmail.com")
	header := "Gateway[" + gwName + "]" + " Status Report Email"
	m.SetHeader("Subject", header)
	m.SetBody("text/html", content)

	d := gomail.NewDialer("smtp.gmail.com", 587, "example@gmail.com", "password")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
