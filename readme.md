<h1 align="center">kafka-traffic-control</h1>
<img src="https://img.shields.io/badge/go-%2300ADD8.svg?&style=for-the-badge&logo=go&logoColor=white"  align="right" position="absolute">

```
go get github.com/vitorhugoze/kafka-traffic-control
```

## Description

**kafka-traffic-control** is a library that has consumer and producer abstractions on top of **github.com/segmentio/kafka-go** for sending inputs from traffic signals and it's IoT sensors and receiving and handling that data on the other end.

### Project Structure

```bash
├───example
│   │   go.mod
│   │   go.sum
│   │   main.go
│   │
│   └───broker - compose file to create a broker appropriate to that example
│           docker-compose.yaml 
│
└───pkg
    ├───kafka
    │   ├───consumer
    │   │       consumer.go
    │   │
    │   └───producer
    │           producer.go
    │
    └───traffic-control - really useful stuff is here
        │   traffic-light-consumer.go
        │   traffic-light-producer.go
        │
        └───dto - objects used to send and receive traffic data
                traffic-control-dtos.go
```

## Usage

#### On the producer side:

```go
	producer := traffic.NewTrafficLightProducer(kafka.WriterConfig{
		Brokers:      []string{"localhost:9092"},
		Topic:        "datatopic",
		BatchTimeout: time.Millisecond * 10,
	})
	
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
```

#### On the consumer side:

```go
	consumer = traffic.NewTrafficLightConsumer(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		Topic:       "datatopic",
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
```