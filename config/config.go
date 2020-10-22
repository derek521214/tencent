package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type (

	Config struct {
		Mysql 		MysqlConfig		`yaml:"mysql"`
		Kafka 		KafkaConfig		`yaml:"kafka"`
		Api 		apiData			`yaml:"api"`
		Email 		EmailInfo		`yaml:"email"`
	}

	apiData struct {
		Url 		string			`yaml:"url"`
		Id 			string			`yaml:"id"`
		Key 		string			`yaml:"key"`
	}

	MysqlConfig struct {
		UserName			string			`yaml:"user_name"`
		Password 			string			`yaml:"password"`
		Host 				string			`yaml:"host"`
		Port				string			`yaml:"port"`
		Db 					string			`yaml:"db"`
		Charset				string			`yaml:"charset"`
		Debug 				bool			`yaml:"debug"`
		MaxIdleConn			int				`yaml:"max_idle_conn"`
		MaxOpenConn  		int				`yaml:"max_open_conn"`
		ConnMaxLifetime		time.Duration	`yaml:"conn_max_lifetime"`
		TablePrefix			string			`yaml:"table_prefix"`
	}


	KafkaConfig struct {
		Address		[]string		`yaml:"address"`	 //kafka 地址端口列表
		Topic		string			`yaml:"topic"`	//主题
		Group 		string			`yaml:"group"`	//消费组别
		Key 		string		 	`yaml:"key"`	// key
		Partition	int32			`yaml:"partition"`	//消费分区
	}

	EmailInfo struct {
		Smtp 			SmtpConfig 		`yaml:"smtp"`
		SendConfig		SendConfig		`yaml:"send_config"`
	}

	SmtpConfig struct {
		Host 		string 			`yaml:"host"`
		Port 		int 			`yaml:"port"`
		UserName 	string 			`yaml:"username"`
		Password    string 			`yaml:"password"`
		FromList 	[]string 		`yaml:"from"`
	}

	SendConfig struct {
		From 		[]FromList		`yaml:"from"`
		To 			[]string		`yaml:"to"`
		Title 		string			`yaml:"title"`
	}

	FromList struct {
		UserName 	string 			`yaml:"username"`
		Password    string 			`yaml:"password"`
	}

)

var (
	Conf *Config
	AppPath string
	AppConfigPath string
)

func init()  {
	getConfigPath()
	Conf = new(Config).initConfig()
}

func(c *Config) initConfig() *Config  {
	yamlFile, err := ioutil.ReadFile(AppConfigPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		panic(err)
	}
	return c
}

func getConfigPath() {
	var err error
	if AppPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		panic(err)
	}
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var filename = "conf.yaml"
	appConfigPath := filepath.Join(workPath, "config", filename)
	if !fileExists(appConfigPath) {
		appConfigPath = filepath.Join(AppPath, "config", filename)
		if !fileExists(appConfigPath) {
			panic("未找到配置文件")
		}
	}
	AppConfigPath = appConfigPath
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}





