package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
	"tencent/config"
	"time"
)

var (
	//URL 中的 random 字段的值
	random string
)

type (
	//请求参数
	SmsRequestParam struct {
		SmsInfo
		Time    int64			`json:"time"`
		Sig     string			`json:"sig"`
	}
	//发送的短信内容
	SmsInfo struct {
		Ext		string			`json:"ext"`
		Extend  string			`json:"extend"`
		Sign    string			`json:"sign"`
		Params  []string 		`json:"params"`
		Tel     mobileParam		`json:"tel"`
		TplId   int64			`json:"tpl_id"`
	}

	//手机号参数
	mobileParam struct {
		Mobile	 		string		`json:"mobile"`
		NationCode	 	string		`json:"nationcode"`
	}

	// 返回参数数据结构
	SmsResultInfo struct {
		Result 	int64 		`json:"result"`
		ErrMsg 	string		`json:"errmsg"`
		Ext 	string		`json:"ext"`
		Fee 	int			`json:"fee"`
		Sid		string		`json:"sid"`
	}

)

// 随机字符
func randStr() int32  {
	rand1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand1.Int31()
}
//sha256加密
func sha(str string) string  {
	h := sha256.New()
	h.Write([]byte(str))
	hash := h.Sum(nil)
	return fmt.Sprintf("%x", hash)
}
//获取 sig 参数值
func(param *SmsRequestParam) getSig () string {
	str := fmt.Sprintf("appkey=%v&random=%v&time=%v&mobile=%v", config.Conf.Api.Key, random,param.Time, param.Tel.Mobile)
	return sha(str)
}
//设置随机数
func (param *SmsRequestParam) setRandom() {
	if random == "" {
		random = fmt.Sprintf("%v", randStr())
	}
}

//清空随机数
func (param *SmsRequestParam)ClearRandom()  {
	random = ""
}

// 获取参数地址
func (param *SmsRequestParam) GetUrl() string  {
	param.setRandom()
	return config.Conf.Api.Url + "?"+"sdkappid="+ config.Conf.Api.Id + "&random="+ random
}

// 获取api post请求参数
func (param *SmsRequestParam) GetUrlParams() ([]byte, error)  {
	param.setRandom()
	param.Sig = param.getSig()
	str, err := json.Marshal(param)
	if err == nil {
		return str, nil
	} else {
		return []byte(""), err
	}
}



