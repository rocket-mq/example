package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/apache/rocketmq-clients/golang/v5"
	"github.com/rocket-mq/rocket-mq/v5"

	c "example/v5"
)

func consume(messageView *golang.MessageView) golang.ConsumerResult {
	fmt.Println("message：", string(messageView.GetBody()))
	fmt.Println("message id：", messageView.GetMessageId())
	fmt.Println("message keys：", messageView.GetKeys())
	fmt.Println("message tag：", messageView.GetTag())
	fmt.Println("message topic：", messageView.GetTopic())
	return golang.SUCCESS
}

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

	ch := make(chan struct{})

	err = v5.New(
		data.Endpoint,
		data.AccessKey,
		data.SecretKey,
		[]string{data.TopicNormal},
		v5.WithConsumerGroup(data.ConsumeGroup),
		v5.WithAwaitDuration(5*time.Second),
		v5.WithNameSpace(data.NameSpace), // 外网必填,内网选填
		v5.WithDebug(data.Debug),
	).PushConsumer(20, 1024, consume, ch)
	if err != nil {
		panic(err)
		return
	}

	time.Sleep(time.Hour * 72)

	ch <- struct{}{}
}
