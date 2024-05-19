package traffic

import (
	"encoding/json"

	"github.com/vitorhugoze/kafka-traffic-control/pkg/kafka/consumer"
	trafficdto "github.com/vitorhugoze/kafka-traffic-control/pkg/traffic-control/dto"

	"github.com/segmentio/kafka-go"
)

type TrafficConsumer struct {
	*consumer.Consumer
}

/*Create a producer specifically to receive traffic data*/
func NewTrafficLightConsumer(config kafka.ReaderConfig) *TrafficConsumer {
	return &TrafficConsumer{
		consumer.NewConsumer(config),
	}
}

/*Process the incoming traffic data with the appropriate handlers*/
func (c *TrafficConsumer) ConsumeTrafficData(trafficHandler ...func(t trafficdto.TrafficLightData) error) {

	c.Listen(func(m kafka.Message) error {

		if string(m.Key) == "traffic-light-data" {

			var data trafficdto.TrafficLightData

			if err := json.Unmarshal(m.Value, &data); err != nil {
				return err
			}

			for _, handler := range trafficHandler {
				if err := handler(data); err != nil {
					return err
				}
			}

		}

		return nil
	})
}
