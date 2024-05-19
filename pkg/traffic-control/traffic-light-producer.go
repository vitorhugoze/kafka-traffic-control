package traffic

import (
	"encoding/json"

	trafficdto "github.com/vitorhugoze/kafka-traffic-control/pkg/traffic-control/dto"

	"github.com/vitorhugoze/kafka-traffic-control/pkg/kafka/producer"

	"github.com/segmentio/kafka-go"
)

type TrafficProducer struct {
	*producer.Producer
}

/*Create a producer specifically to send traffic data*/
func NewTrafficLightProducer(config kafka.WriterConfig) *TrafficProducer {
	return &TrafficProducer{
		producer.NewProducer(config),
	}
}

/*Sends traffic data as json using kafka*/
func (t *TrafficProducer) SendTrafficData(data trafficdto.TrafficLightData) error {

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	t.SendMessage([]byte("traffic-light-data"), jsonBytes)
	return nil
}
