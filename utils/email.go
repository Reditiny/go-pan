package utils

import (
	"fmt"
	"github.com/beego/beego/v2/server/web"
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
)

/**
邮件工具
*/

var (
	adminEmail, _ = web.AppConfig.String("email")
	password, _   = web.AppConfig.String("emailpass")
	emailsmtp, _  = web.AppConfig.String("emailsmtp")
)

func SendEmail(targetEmail, captcha string) {
	e := email.NewEmail()

	e.From = fmt.Sprintf("go-pan 注册 <%s>", adminEmail)
	e.To = []string{targetEmail}
	e.Subject = "验证码"
	e.Text = []byte(captcha)
	err := e.Send(fmt.Sprintf("%s:25", emailsmtp), smtp.PlainAuth("", adminEmail, password, emailsmtp))
	if err != nil {
		log.Fatal(err)
	}
}
