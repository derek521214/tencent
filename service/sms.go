package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"tencent/api"
	"tencent/config"
	"tencent/models"
	"tencent/email"
	"time"
)

type KafkaData struct {
	Type 		string			`json:"type"`
	Data 		api.SmsInfo		`json:"data"`
}

//http发送短信
func sendHttp(smsContent api.SmsInfo)(api.SmsResultInfo,  error)  {
	//发送请求
	var param = new(api.SmsRequestParam)
	param.Sign = smsContent.Sign
	param.Tel.Mobile = smsContent.Tel.Mobile
	param.Tel.NationCode = smsContent.Tel.NationCode
	param.TplId = smsContent.TplId
	param.Params = smsContent.Params
	param.Time = time.Now().Unix()
	//清空随机数
	param.ClearRandom()
	data, err := param.GetUrlParams()
	if err == nil {
		req, err := http.NewRequest("POST", param.GetUrl(), bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err == nil {
			body, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				var result api.SmsResultInfo
				err := json.Unmarshal(body, &result)
				if err==nil {
					return result, nil
				} else {
					return result, err
				}
			}
		}
	}
	return api.SmsResultInfo{}, err
}

//获取发送短信模板信息
func getSmsContent(tplId int64) (*models.SmsContent, bool)  {
	sms := new(models.SmsContent)
	res := sms.GetOne(tplId)
	return sms, res
}

//添加发送日志
func addSmsLog(smsInfo *api.SmsInfo, result *api.SmsResultInfo) {
	go func(smsInfo *api.SmsInfo, result *api.SmsResultInfo) {
		//获取模板信息内容
		smsContent, res := getSmsContent(smsInfo.TplId)
		if res {
			var logInfo models.SmsSendLog
			//实际发送短信内容
			logInfo.Content = getSmsSendContent(smsInfo.Params, smsContent)
			logInfo.Sid = result.Sid
			logInfo.Fee = result.Fee
			logInfo.ErrMsg = result.ErrMsg
			logInfo.Result = result.Result
			logInfo.Sign = smsContent.SmsSign
			mobile, err:= strconv.ParseUint(smsInfo.Tel.Mobile, 10, 64)
			if err== nil {
				logInfo.Mobile = mobile
			}
			logInfo.NationCode = smsInfo.Tel.NationCode
			logInfo.SendTime = time.Now().Unix()
			logInfo.Ctime = time.Now().Unix()
			logInfo.TplId = smsInfo.TplId
			logInfo.AddOne()
		} else {
			logInfo("获取不到短信模板" + " " + NowTime())
		}
	}(smsInfo, result)
}

//获取实际发送的短信内容
func getSmsSendContent(params []string, smsContent *models.SmsContent) string {
	contentStr := smsContent.Content
	smsParamName := strings.Split(smsContent.SmsParam, ",")
	for i, val := range params {
		contentStr = strings.Replace(contentStr, smsParamName[i], val, -1)
	}
	return contentStr
}

//处理kafka消息
func Run(message string){
	msg := []byte(message)
	var data KafkaData
	err:=json.Unmarshal(msg, &data)
	if err==nil {
		//判断是否短信
		if data.Type == config.Conf.Kafka.Key {
			go send(&data)
		}
	} else {
		logInfo("解析kafka数据失败数据信息:" + message+ " " + NowTime())
		logInfo("解析数据失败原因:" + err.Error()+ " " + NowTime())
	}
}

//发送短信
func send(data *KafkaData) {
	var smsInfo api.SmsInfo
	smsInfo = data.Data
	result, err := sendHttp(smsInfo)
	if err == nil {
		resultStr := fmt.Sprintf("接口返回信息%v", result)
		logInfo(resultStr)
		smsStr := fmt.Sprintf("短信内容%v", smsInfo)
		logInfo(smsStr)
		if result.Result != 0 {
			email.Warn("短信发送失败","号码：" + smsInfo.Tel.Mobile + "短信模板ID "+ fmt.Sprintf("%v",smsInfo.TplId) +" 发送失败错误信息："+ result.ErrMsg+ " " + NowTime())
		} else {
			addSmsLog(&smsInfo, &result)
			logInfo("发送完成")
			return
		}
	} else {
		logInfo(fmt.Sprintf("发送失败错误信息%v", err) + " " + NowTime())
	}
}

//添加日志信息
func logInfo(msg string)  {
	go func(msg string) {
		log := new(models.KafkaLog)
		log.Type = models.TYPE
		log.Msg = msg
		log.AddOne()
	}(msg)
}

func NowTime() string {
	return time.Now().Format("2006-01-02 15:04")
}
