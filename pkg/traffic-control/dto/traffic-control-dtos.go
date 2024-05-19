package trafficdto

const (
	TRAFFIC_LIGHT_STOP    TrafficLightState = 0
	TRAFFIC_LIGHT_WARNING TrafficLightState = 1
	TRAFFIC_LIGHT_PROCEED TrafficLightState = 2
)

type TrafficLightState int

/*Individual vehicle data that can be optionally sent*/
type VehicleData struct {
	Plate string
	Speed float32
}

/*Structured data to send traffic information to consumer*/
type TrafficLightData struct {
	Id             int
	CurrentState   TrafficLightState
	TrafficFlow    float32
	PedestrianFlow float32
	Vehicles       []VehicleData
}
