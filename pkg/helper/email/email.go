package email

import (
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"strings"
)

var (
	username string = os.Getenv("EmailUsername")
	password string = os.Getenv("EmailPassword")
	host     string = os.Getenv("EmailType")
)

func SendCode(context string, email string) error {
	to := email
	subject := `Oldgeneral`
	body := `<html><body><h1>您的验证码为</h1><h3>` + context + `</h3><br/>验证码有效期为1小时，请在1小时内完成验证<br/>如果不是您本人操作，请忽略本条邮件</body></html>`
	err := sendToMail(username, password, host, to, subject, body)
	if err != nil {
		return err
	}
	return nil
}

func sendToMail(user, password, host, to, subject, body string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var contentType string
	contentType = "Content-Type: text/html; charset=UTF-8"
	msg := []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	sendTo := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, sendTo, msg)
	log.Println("sent")
	return err
}

func GenerateVerificationCode() string {
	var result string
	directory := `0123456789ABCDEFGHJKLMNPQRSTUVWXYZ`
	for i := 0; i < 6; i++ {
		a := rand.Int() % 34
		result += directory[a : a+1]
	}
	return result
}
