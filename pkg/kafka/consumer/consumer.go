package consumer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	readerConfig  kafka.ReaderConfig
	cancelContext context.CancelFunc
}

func NewConsumer(config kafka.ReaderConfig) *Consumer {
	return &Consumer{
		readerConfig: config,
	}
}

/*Listen for messages and use the handlers to process them*/
func (c *Consumer) Listen(msgHandlers ...func(m kafka.Message) error) {

	ctx, cancel := context.WithCancel(context.Background())
	c.cancelContext = cancel

	r := kafka.NewReader(c.readerConfig)
	defer r.Close()

	for {

		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic(err)
		}

		for _, handler := range msgHandlers {
			if err := handler(msg); err != nil {
				panic(err)
			}
		}

		if c.readerConfig.GroupID != "" {
			r.CommitMessages(ctx, msg)
		}
	}
}

func (c *Consumer) Close() {
	c.cancelContext()
}
