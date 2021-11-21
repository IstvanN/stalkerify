package main

import (
	"fmt"
	"net/smtp"
	"os"
)

var (
	mailFrom     = os.Getenv("MAIL_FROM")
	mailTo       = os.Getenv("MAIL_TO")
	mailUser     = os.Getenv("MAIL_USER")
	mailPw       = os.Getenv("MAIL_PW")
	mailSmtpHost = os.Getenv("MAIL_SMTP_HOST")
	mailSmtpPort = os.Getenv("MAIL_SMTP_PORT")
)

func sendMail(playlistName string) error {
	msg := []byte(fmt.Sprintf("To: %v\r\n"+
		"Subject: New song in %v!\r\n"+
		"\r\n"+
		"Hey there! Maybe you should know: there's a new song in %v. Go check that out!\r\n", mailTo, playlistName, playlistName))
	auth := smtp.PlainAuth("", mailUser, mailPw, mailSmtpHost)

	if err := smtp.SendMail(mailSmtpHost+":"+mailSmtpPort, auth, mailFrom, []string{mailTo}, msg); err != nil {
		return fmt.Errorf("error sending mail: %v", err)
	}

	return nil
}
