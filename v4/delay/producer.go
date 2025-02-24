package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aliyunmq/mq-http-go-sdk"
	"github.com/rocket-mq/rocket-mq/v4/producer"
	"github.com/rocket-mq/rocket-mq/v4/producer/message"

	"example/v4"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "./v4/config.json", "local config file path")

	flag.Parse()

	bytes, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
		return
	}

	var data v4.Config
	if err = json.Unmarshal(bytes, &data); err != nil {
		panic(err)
		return
	}

	p := producer.New(data.Endpoint, data.AccessKey, data.SecretKey, data.TopicDelay, data.InstanceId)

	for i := 0; i < 10; i++ {
		body := fmt.Sprintf("%s-%d", "wb", i)
		msg := message.New(
			body,
			message.WithTag("test-2-tag"),
			message.WithKey("test-2-key"),
			message.WithProperties(map[string]string{"time": time.Now().Format("20060102150405")}),
			message.WithDelayTime(time.Now().Add(time.Minute)), // 延时消息必传，推荐使用该方式
			//message.WithStartDeliverTime(time.Now().Unix()*1000+60*1000), // 以上二选一
		)
		var pmr mq_http_sdk.PublishMessageResponse
		if pmr, err = p.PublishMessage(msg); err != nil {
			panic(err)
			return
		}
		fmt.Println(pmr.MessageId, " -------- ", body)
	}
}
