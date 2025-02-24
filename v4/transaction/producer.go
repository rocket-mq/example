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

func checker(cme mq_http_sdk.ConsumeMessageEntry) producer.TransactionResolution {

	if cme.MessageBody == "wb-4" {
		fmt.Println(cme.MessageId, " ================ ", cme.MessageBody, " ================ ", cme.ConsumedTimes, " ================ ", producer.ROLLBACK)
		return producer.ROLLBACK
	}

	if cme.ConsumedTimes > 2 {
		fmt.Println(cme.MessageId, " ================ ", cme.MessageBody, " ================ ", cme.ConsumedTimes, " ================ ", producer.COMMIT)
		return producer.COMMIT
	}

	fmt.Println(cme.MessageId, " ================ ", cme.MessageBody, " ================ ", cme.ConsumedTimes, " ================ ", producer.UNKNOWN)
	return producer.UNKNOWN
}

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

	p := producer.New(
		data.Endpoint,
		data.AccessKey,
		data.SecretKey,
		data.TopicTransaction,
		data.InstanceId,
		producer.WithTrans(data.GroupId, checker),
	)

	for i := 0; i < 10; i++ {
		body := fmt.Sprintf("%s-%d", "wb", i)
		msg := message.New(
			body,
			message.WithTag("test-4-tag"),
			message.WithKey("test-4-key"),
			message.WithProperties(map[string]string{"time": time.Now().Format("20060102150405")}),
			message.WithTransCheckImmunityTime(15),
		)
		var pmr mq_http_sdk.PublishMessageResponse
		if pmr, err = p.PublishTransMessage(msg); err != nil {
			panic(err)
			return
		}

		status := producer.UNKNOWN

		if i == 0 {
			if err = p.Commit(pmr.ReceiptHandle); err != nil {
				panic(err)
				return
			}
			status = producer.COMMIT
		}

		if i == 7 {
			if err = p.Rollback(pmr.ReceiptHandle); err != nil {
				panic(err)
				return
			}
			status = producer.ROLLBACK
		}

		fmt.Println(pmr.MessageId, " -------- ", body, " -------- ", status)
	}

	select {}
}
