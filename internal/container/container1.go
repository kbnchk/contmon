package container

import (
	"github.com/kbnchk/contmon/internal/device"
)

type container1 struct {
	serial, parity                string
	baudrate, databits, stopbits  int
	meteoAddr, fanAddr, meterAddr byte
}

func Container1() Container {
	return &container1{
		serial:    "/dev/ttyUSB0",
		baudrate:  9600,
		databits:  8,
		stopbits:  1,
		parity:    "N",
		meteoAddr: 1,
		fanAddr:   2,
		meterAddr: 101,
	}
}

// GetData gets data from sensors and represent it in Data struct.
// Modbus-RTU id a serial communication standert so I close each device handler before opening new one
// instead of using defer statement because it makes all handler work simultaneously causing timeouts.
func (c *container1) GetData() Data {
	meteoOk := true
	fanOk := true
	meterOk := true

	meteo, err := device.CWS19New(c.serial, c.meteoAddr, c.baudrate, c.databits, c.stopbits, c.parity)
	if err != nil {
		meteoOk = false
	}
	temp, err := meteo.GetTemperature()
	if err != nil {
		meteoOk = false
	}
	hum, err := meteo.GetHumidity()
	if err != nil {
		meteoOk = false
	}
	meteo.Close()

	fan, err := device.ESQ760New(c.serial, c.fanAddr, c.baudrate, c.databits, c.stopbits, c.parity)
	if err != nil {
		fanOk = false
	}
	state, err := fan.GetStatus()
	if err != nil {
		fanOk = false
	}
	errorcode, err := fan.GetError()
	if err != nil {
		fanOk = false
	}
	freq, err := fan.GetFreq()
	if err != nil {
		fanOk = false
	}
	fan.Close()

	meter, err := device.ZM194New(c.serial, c.meterAddr, c.baudrate, c.databits, c.stopbits, c.parity)
	if err != nil {
		meterOk = false
	}
	defer meter.Close()

	p1v, p2v, p3v, err := meter.GetVoltage()
	if err != nil {
		meterOk = false
	}

	p1p, p2p, p3p, err := meter.GetPower()
	if err != nil {
		meterOk = false
	}
	meter.Close()

	return Data{
		Meteo: Meteo{
			Ok:       meteoOk,
			Temp:     temp,
			Humidity: hum,
		},
		Fan: Fan{
			Ok:        fanOk,
			Frequency: freq,
			State:     state,
			ErrorCode: errorcode,
		},
		Meter: Meter{
			Ok: meterOk,
			Voltage: Voltage{
				Phase1: p1v,
				Phase2: p2v,
				Phase3: p3v,
			},
			Power: Power{
				Phase1: p1p,
				Phase2: p2p,
				Phase3: p3p,
			},
		},
	}
}
