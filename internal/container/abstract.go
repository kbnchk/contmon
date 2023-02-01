package container

type Container interface {
	GetData() (Data, error)
}

type Data struct {
	Meteo Meteo
	Fan   Fan
	Meter Meter
}

type Meteo struct {
	Temp     float64
	Humidity float64
}

type Fan struct {
	Frequency float64
	State     string
	ErrorCode uint8
}

type Meter struct {
	Voltage Voltage
	Power   Power
}

type Voltage struct {
	Phase1, Phase2, Phase3 uint16
}

type Power struct {
	Phase1, Phase2, Phase3 uint16
}
