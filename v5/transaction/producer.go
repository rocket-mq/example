package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/apache/rocketmq-clients/golang/v5"
	"github.com/rocket-mq/rocket-mq/v5"
	"github.com/rocket-mq/rocket-mq/v5/producer"

	c "example/v5"
)

var n int32

func checker(mv *golang.MessageView) golang.TransactionResolution {
	atomic.AddInt32(&n, 1)

	if string(mv.GetBody()) == "wb-4" {
		fmt.Println(mv.GetMessageId(), " ================ ", string(mv.GetBody()), " ================ ", golang.ROLLBACK)
		return golang.ROLLBACK
	}

	if _n := atomic.LoadInt32(&n); _n > 20 {
		fmt.Println(mv.GetMessageId(), " ================ ", string(mv.GetBody()), " ================ ", _n, " ================ ", golang.COMMIT)
		return golang.COMMIT
	}

	fmt.Println(mv.GetMessageId(), " ================ ", string(mv.GetBody()), " ================ ", golang.UNKNOWN)
	return golang.UNKNOWN
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

	topics := []string{data.TopicNormal, data.TopicFifo, data.TopicDelay, data.TopicTransaction}

	p, err := v5.New(
		data.Endpoint,
		data.AccessKey,
		data.SecretKey,
		topics,
		v5.WithNameSpace(data.NameSpace), // 外网必填,内网选填
		v5.WithDebug(false),
		v5.WithTransactionChecker(checker),
	).Producer()
	if err != nil {
		panic(err)
		return
	}

	defer p.Close()

	for i := 0; i < 10; i++ {
		body := fmt.Sprintf("%s-%d", "wb", i)
		message := producer.NewMessage(
			data.TopicTransaction,
			[]byte(body),
			producer.WithTag("test-4-tag"),
			producer.WithKeys("test-4-key", time.Now().Format(time.DateTime)),
		)
		var resp []*golang.SendReceipt
		var t golang.Transaction
		if resp, t, err = p.SendWithTransaction(context.Background(), message); err != nil {
			panic(err)
			return
		}

		status := golang.UNKNOWN

		if i == 0 {
			if err = t.Commit(); err != nil {
				panic(err)
				return
			}
			status = golang.COMMIT
		}

		if i == 7 {
			if err = t.RollBack(); err != nil {
				panic(err)
				return
			}
			status = golang.ROLLBACK
		}

		for _, r := range resp {
			fmt.Println(r.MessageID, " -------- ", body, " -------- ", status)
		}
	}

	select {}
}
