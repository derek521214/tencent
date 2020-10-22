package models

type SmsSendLog struct {
	Id               	uint64				`gorm:"primary_key"`
	Sid      			string
	Fee         		int
	ErrMsg 				string
	Result	         	int64
	State           	int16
	Sign        		string
	Mobile          	uint64
	NationCode          string
	SendTime			int64
	TplId				int64
	Content 			string
	Ctime 				int64
}

func (log SmsSendLog) AddOne() bool {
	resDb := db.Create(&log)
	if resDb.Error ==nil && log.Id>0 {
		return true
	}
	return false
}


