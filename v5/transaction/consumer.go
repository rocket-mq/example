package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/apache/rocketmq-clients/golang/v5"
	"github.com/rocket-mq/rocket-mq/v5"

	c "example/v5"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "./v5/config.json", "local config file path")

	flag.Parse()

	bytes, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
		return
	}

	var data c.Config
	if err = json.Unmarshal(bytes, &data); err != nil {
		panic(err)
		return
	}

	consumer, err := v5.New(
		data.Endpoint,
		data.AccessKey,
		data.SecretKey,
		[]string{data.TopicTransaction},
		v5.WithConsumerGroup(data.ConsumeGroup),
		v5.WithAwaitDuration(5*time.Second),
		v5.WithNameSpace(data.NameSpace), // 外网必填,内网选填
		v5.WithDebug(data.Debug),
	).SimpleConsumer()
	if err != nil {
		panic(err)
		return
	}

	defer consumer.Close()

	for {
		var list []*golang.MessageView
		if list, err = consumer.Receive(context.Background(), 3, 20*time.Second); err != nil {
			if !v5.IsMessageNotFoundErr(err) {
				panic(err)
				return
			}
			continue
		}
		for _, msg := range list {
			fmt.Println("message：", string(msg.GetBody()))
			fmt.Println("message id：", msg.GetMessageId())
			fmt.Println("message keys：", msg.GetKeys())
			fmt.Println("message tag：", msg.GetTag())
			fmt.Println("message topic：", msg.GetTopic())
			fmt.Println("--------------------------------------")
			_ = consumer.Ack(context.Background(), msg)
		}
	}
}
