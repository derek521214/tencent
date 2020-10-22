package email

import "tencent/config"

var (
	start = 0
)

func Warn( subject string, content string) {
	email := new(EmailInfo)
	fromInfo := getFromEmail(config.Conf.Email.SendConfig)
	email.From = fromInfo.UserName
	email.To = config.Conf.Email.SendConfig.To
	email.Subject = config.Conf.Email.SendConfig.Title + subject
	email.Content = content
	config.Conf.Email.Smtp.UserName = fromInfo.UserName
	config.Conf.Email.Smtp.Password = fromInfo.Password
	email.Send(config.Conf.Email.Smtp)
}

//获取发送地址 轮询方式
func getFromEmail(config config.SendConfig) config.FromList {
	if (start+1) >= len(config.From) {
		start = 0
	}
	start++
	return config.From[start-1]
}
