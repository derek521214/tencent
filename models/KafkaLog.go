package models

type KafkaLog struct {
	Id               	uint64			`gorm:"primary_key"`
	Type 				uint8
	Msg 				string
}

const TYPE = 2

func (log *KafkaLog) AddOne() bool  {
	res := db.Create(&log)
	if res.Error == nil && res.RowsAffected >0 {
		return  true
	}
	return false
}

