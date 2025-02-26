package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/aliyunmq/mq-http-go-sdk"
	"github.com/rocket-mq/rocket-mq/v4/consumer"

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

	c := consumer.New(data.Endpoint, data.AccessKey, data.SecretKey, data.TopicNormal, data.InstanceId, data.GroupId)

	var n int32

	ch := make(chan struct{}, 1000)

	for {
		ch <- struct{}{}
		go func() {
			defer func() {
				<-ch
			}()

			if err = c.Receive(3, 3, func(cme mq_http_sdk.ConsumeMessageEntry) {
				atomic.AddInt32(&n, 1)
				time.Sleep(3 * time.Second)
				fmt.Println(
					cme.MessageId, " ---- ",
					cme.MessageBody, " ---- ",
					cme.MessageTag, " ---- ",
					cme.MessageKey, " ---- ",
					cme.Properties, " ----  ",
					c.Ack(cme), " ---- ",
					atomic.LoadInt32(&n))
			}); err != nil {
				panic(err)
				return
			}
		}()

	}

}
