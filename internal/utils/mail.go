package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

var (
	fromEmail         string
	fromEmailPassword string
	emailSmtpHost     string
	emailSmtpPort     string
)

func InitMail() {
	fromEmail = os.Getenv("FROM_EMAIL")
	fromEmailPassword = os.Getenv("FROM_EMAIL_PASSWORD")
	emailSmtpHost = os.Getenv("FROM_EMAIL_SMTP_HOST")
	emailSmtpPort = os.Getenv("FROM_EMAIL_SMTP_PORT")
	if fromEmail == "" || fromEmailPassword == "" || emailSmtpHost == "" || emailSmtpPort == "" {
		log.Fatal("Missing email variable environment")
	}
}

func SendEmail(to string, subject string, htmlBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fromEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	port, err := strconv.Atoi(emailSmtpPort)
	if err != nil {
		log.Println("Invalid SMTP port:", err)
		return err
	}
	d := gomail.NewDialer(emailSmtpHost, port, fromEmail, fromEmailPassword)

	if err := d.DialAndSend(m); err != nil {
		log.Println("Send email failed:", err)
		return err
	}
	return err
}

func SendVerificationEmail(userEmail string, token string) error {
	verificationLink := fmt.Sprintf("http://localhost:8080/mail/verify-account?token=%s", token)

	subject := "Verify your account"
	html := fmt.Sprintf(`
		<h2>Welcome!</h2>
		<p>Please verify your email by clicking the link below:</p>
		<a href="%s">Verify Email</a>
	`, verificationLink)

	return SendEmail(userEmail, subject, html)
}

func SendResetPassword(userEmail string, resetPassword string) error {
	subject := "Reset Your Password"
	html := fmt.Sprintf(`
		<h2>Hello,</h2>
		<p>You requested a password reset. Here is your new password:</p>
		<p style="font-size: 18px; font-weight: bold;">%s</p>
		<br>
		<p>Regards,<br>Hotel Management Team</p>
	`, resetPassword)
	return SendEmail(userEmail, subject, html)
}

func SendStaffPassword(staffEmail string, password string) error {
	subject := "Your Password"
	html := fmt.Sprintf(`
		<h2>Hello,</h2>
		<p>Here is your password:</p>
		<p style="font-size: 18px; font-weight: bold;">%s</p>
		<br>
		<p>Regards,<br>Hotel Management Team</p>
	`, password)
	return SendEmail(staffEmail, subject, html)
}
