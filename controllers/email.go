package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
)

type EmailRequest struct {
	Target  string `json:"target"`
	Content string `json:"content"`
	Secret  string `json:"secret"`
}

func SendEmail(res http.ResponseWriter, req *http.Request) {
	var emailReq EmailRequest
	err := json.NewDecoder(req.Body).Decode(&emailReq)
	if err != nil {
		http.Error(res, "Invalid request format", http.StatusBadRequest)
		return
	}

	SECRET := os.Getenv("API_SECRET")
	if emailReq.Secret != SECRET {
		http.Error(res, "Invalid secret", http.StatusUnauthorized)
		return
	}

	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SMTP_EMAIL")
	senderPassword := os.Getenv("SMTP_PASSWORD")

	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: Test Email\nContent-Type: text/html\n\n%s",
		senderEmail, emailReq.Target, emailReq.Content)

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)

	err = smtp.SendMail(smtpServer+":"+smtpPort, auth, senderEmail, []string{emailReq.Target}, []byte(message))
	if err != nil {
		http.Error(res, "Failed to send email", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Email sent successfully"))
}
