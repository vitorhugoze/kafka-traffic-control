package producer

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	dataChannel   chan ProducerMessage
	cancelContext context.CancelFunc
}

type ProducerMessage struct {
	Key, Value []byte
}

func NewProducer(config kafka.WriterConfig) *Producer {

	dataChan := make(chan ProducerMessage)
	ctx, cancel := context.WithCancel(context.Background())

	w := kafka.NewWriter(config)

	go func() {
		for {

			select {
			case <-ctx.Done():
				fmt.Println("Terminating producer...")
				w.Close()
				return
			case msg := <-dataChan:
				err := w.WriteMessages(ctx, kafka.Message{
					Key:   msg.Key,
					Value: msg.Value,
				})
				if err != nil {
					panic(err)
				}
			}

		}
	}()

	return &Producer{
		dataChannel:   dataChan,
		cancelContext: cancel,
	}
}

/*Sends message to broker, must first initialize producer*/
func (p *Producer) SendMessage(key, value []byte) {

	p.dataChannel <- ProducerMessage{
		Key:   key,
		Value: value,
	}
}

func (p *Producer) Close() {
	p.cancelContext()
}
