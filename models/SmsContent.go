package models


type SmsContent struct {
	Title      			string
	SmsSign         	string
	SmsType 			uint8
	SmsParam	        string
	Content           	string
	Status        		uint8
	Ctime 				int64
}

func (smsContent *SmsContent) GetOne(tplId int64) bool  {
	res :=db.Where("sms_tpl=?", tplId).First(smsContent)
	if res.Error == nil && smsContent.Title != "" {
		return  true
	}
	return false
}
