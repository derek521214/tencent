package consumer

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"tencent/config"
	"tencent/email"
	"tencent/service"
)

func Run(topic string, group string)  {
	if topic != "" {
		config.Conf.Kafka.Topic = topic
	}
	if group != "" {
		config.Conf.Kafka.Group = group
	}
	clusterConsumer()
}

// 支持brokers cluster的消费者
/*func clusterConsumer(kafkaConfig config.KafkaConfig)  {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   kafkaConfig.Address,
		GroupID:   kafkaConfig.Group,
		Topic:     kafkaConfig.Topic,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		CommitInterval: time.Second,
	})
	defer r.Close()
	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			email.Warn("kafka短信消费任务错误", "错误信息：" + fmt.Sprintf("%v", err))
			break
		}
		consumerMsg(msg, kafkaConfig)
	}
}*/

func clusterConsumer()  {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   config.Conf.Kafka.Address,
		GroupID:   config.Conf.Kafka.Group,
		Topic:     config.Conf.Kafka.Topic,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	defer r.Close()
	ctx := context.Background()
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			email.Warn("kafka短信消费任务错误", "错误信息：" + fmt.Sprintf("%v", err) + " " + service.NowTime())
			break
		}
		go consumerMsg(msg)
		r.CommitMessages(ctx, msg)
	}
}

func clusterConsumerNew()  {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   config.Conf.Kafka.Address,
		GroupID:   config.Conf.Kafka.Group,
		Topic:     config.Conf.Kafka.Topic,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	defer r.Close()
	ctx := context.Background()
	var msgChan = make(chan kafka.Message, 20)
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			email.Warn("kafka短信消费任务错误", "错误信息：" + fmt.Sprintf("%v", err) + " " + service.NowTime())
			break
		}
		go func() {
			msgChan <- msg
		}()
		select {
			case info, ok := <- msgChan :
				if ok {
					go consumerMsg(info)
				}
			default:

		}
		r.CommitMessages(ctx, msg)
	}
}

//处理kafka信息
func consumerMsg(msg kafka.Message)  {
	//fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
	if string(msg.Key) == config.Conf.Kafka.Key {
		//time.Sleep(2 * time.Second)
		service.Run(string(msg.Value))
	}
}