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
	"github.com/rocket-mq/rocket-mq/v5/producer"

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

	topics := []string{data.TopicNormal, data.TopicFifo, data.TopicDelay, data.TopicTransaction}

	p, err := v5.New(
		data.Endpoint,
		data.AccessKey,
		data.SecretKey,
		topics,
		v5.WithNameSpace(data.NameSpace), // 外网必填,内网选填
		v5.WithDebug(true),
	).Producer()
	if err != nil {
		panic(err)
		return
	}

	defer p.Close()

	for i := 0; i < 10; i++ {
		body := fmt.Sprintf("%s-%d", "wb", i)
		message := producer.NewMessage(
			data.TopicNormal,
			[]byte(body),
			producer.WithTag("test-1-tag"),
			producer.WithKeys("test-1-key", time.Now().Format(time.DateTime)),
		)
		var resp []*golang.SendReceipt
		if resp, err = p.Send(context.Background(), message); err != nil {
			panic(err)
			return
		}
		for _, r := range resp {
			fmt.Println(r.MessageID, " -------- ", body)
		}
	}
}
