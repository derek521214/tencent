package email

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"tencent/config"
)

type (
	EmailInfo struct {
		From 		string			`yaml:"from"`
		To	 		[]string		`yaml:"to"`
		Cc	 		Cc				`yaml:"cc"`
		Subject	 	string			`yaml:"subject"`
		Attach   	string			`yaml:"attach"`
		Content 	string			`yaml:"content"`
	}

	Cc struct {
		Address 	string		`yaml:"address"`
		Name	 	string		`yaml:"name"`
	}
)

func (conf *EmailInfo) Send (smtp config.SmtpConfig) (bool, error)  {
	m := gomail.NewMessage()
	if conf.From == "" {
		return false, fmt.Errorf("发件人地址必须填写")
	}
	m.SetHeader("From", conf.From)
	if len(conf.To) == 0 {
		return false, fmt.Errorf("收件人地址必须填写")
	}
	m.SetHeader("To", conf.To...)
	if conf.Cc.Address != "" {
		m.SetAddressHeader("Cc", conf.Cc.Address, conf.Cc.Name)
	}
	m.SetHeader("Subject", conf.Subject)
	if conf.Content == "" {
		return false,fmt.Errorf("邮件内容不能为空")
	}
	m.SetBody("text/html", conf.Content)
	if conf.Attach != "" {
		m.Attach(conf.Attach)
	}
	d := gomail.NewDialer(smtp.Host, smtp.Port, smtp.UserName, smtp.Password)
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return false, err
	} else {
		return true, nil
	}
}
