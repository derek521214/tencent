package models

import (
	"time"
	"tencent/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"database/sql"
)

var (
	db *gorm.DB
	sqlDB *sql.DB
)

func init()  {
	var err error
	db, err = gorm.Open(mysql.Open(getMysqlConnectString(config.Conf.Mysql)), &gorm.Config{
		NamingStrategy : schema.NamingStrategy{
			TablePrefix:   config.Conf.Mysql.TablePrefix,
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(config.Conf.Mysql.MaxIdleConn)
	sqlDB.SetMaxOpenConns(config.Conf.Mysql.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Second * config.Conf.Mysql.ConnMaxLifetime)
}

func getMysqlConnectString(conf config.MysqlConfig) string  {
	var linkString string
	linkString = conf.UserName + ":" + conf.Password + "@tcp("+ conf.Host + ":" + conf.Port + ")/" + conf.Db + "?charset=" + conf.Charset
	return linkString
}

