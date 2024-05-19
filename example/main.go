package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/vitorhugoze/kafka-traffic-control/pkg/traffic-control"
	trafficdto "github.com/vitorhugoze/kafka-traffic-control/pkg/traffic-control/dto"

	"github.com/segmentio/kafka-go"
)

const (
	broker = "localhost:9092"
	topic  = "datatopic"
)

var (
	producer *traffic.TrafficProducer
	consumer *traffic.TrafficConsumer
)

func main() {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		<-sigChan

		producer.Close()
		consumer.Close()

		fmt.Println("Terminating kafka service...")
		os.Exit(0)
	}()

	go sendTrafficData()
	receiveTrafficData()
}

/*Send some mock traffic data every second*/
func sendTrafficData() {

	producer = traffic.NewTrafficLightProducer(kafka.WriterConfig{
		Brokers:      []string{broker},
		Topic:        topic,
		BatchTimeout: time.Millisecond * 10,
	})

	go func() {
		i := 0
		for {
			i++

			time.Sleep(time.Millisecond * 1000)

			producer.SendTrafficData(trafficdto.TrafficLightData{
				Id:             15792,
				CurrentState:   trafficdto.TRAFFIC_LIGHT_PROCEED,
				TrafficFlow:    0.6,
				PedestrianFlow: 0,
				Vehicles: []trafficdto.VehicleData{
					{
						Plate: "MHI-7297",
						Speed: 57.0,
					},
				},
			})
		}
	}()

}

/*Process traffic data received using custom handler*/
func receiveTrafficData() {

	consumer = traffic.NewTrafficLightConsumer(kafka.ReaderConfig{
		Brokers:     []string{broker},
		Topic:       topic,
		GroupID:     "group-01",
		StartOffset: kafka.FirstOffset,
	})

	consumer.ConsumeTrafficData(func(t trafficdto.TrafficLightData) error {
		fmt.Printf("\nTraffic data received from traffic light: %v", t.Id)

		if len(t.Vehicles) > 0 {

			for _, v := range t.Vehicles {

				if t.CurrentState == trafficdto.TRAFFIC_LIGHT_STOP {
					fmt.Printf("\nVehicle with plate: %v didn't stop on red light!", v.Plate)
				}

				if v.Speed > 50 {
					fmt.Printf("\nVehicle with plate: %v over speed limit!", v.Plate)
				}
			}

		}

		return nil
	})
}
