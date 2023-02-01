package container

type Container interface {
	GetData() Data
}

type Data struct {
	Meteo Meteo
	Fan   Fan
	Meter Meter
}

type Meteo struct {
	Ok       bool
	Temp     float64
	Humidity float64
}

type Fan struct {
	Ok        bool
	Frequency float64
	State     string
	ErrorCode uint8
}

type Meter struct {
	Ok      bool
	Voltage Voltage
	Power   Power
}

type Voltage struct {
	Phase1, Phase2, Phase3 uint16
}

type Power struct {
	Phase1, Phase2, Phase3 uint16
}
