package main

import (
	"flag"
	"fmt"
	"time"
	"tencent/consumer"
)


var (
	topic 		string
	group 		string
)

func main()  {
	flag.StringVar(&topic, "topic", "","设置kafka主题")
	flag.StringVar(&group, "group", "","设置kafka消费组")
	flag.Parse()
	fmt.Println("开始运行kafka短信消费任务，时间：" + time.Now().Format("2006-01-02 15:04"))
	consumer.Run(topic, group)
}