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

func sendMail(playlistName string, newSongs []newSongData) error {
	msg := []byte(fmt.Sprintf(
		"From: stalkerify@archiez.io\r\n"+
			"To: %v\r\n"+
			"Subject: New song(s) in %v!\r\n"+
			"\r\n"+
			formMessage(newSongs), mailTo, playlistName))

	auth := smtp.PlainAuth("", mailUser, mailPw, mailSmtpHost)

	if err := smtp.SendMail(mailSmtpHost+":"+mailSmtpPort, auth, mailFrom, []string{mailTo}, msg); err != nil {
		return fmt.Errorf("error sending mail: %v", err)
	}

	return nil
}

func formMessage(newSongs []newSongData) string {
	finalMessage := "Hey there! Here is a list of the new songs:\r\n" + "\r\n"

	for i, ns := range newSongs {
		finalMessage += fmt.Sprintf("%d. %s - %s added by %v at %v\r\n", i+1, ns.artist, ns.title, ns.addedBy, ns.addedAt) + "\r\n"
	}

	return finalMessage + "Regards: Stalkerify created by Archie"
}
