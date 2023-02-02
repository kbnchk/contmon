package container

import (
	"time"

	"github.com/kbnchk/contmon/internal/device"
	"github.com/kbnchk/contmon/internal/entity"
)

type Container interface {
	GetData() entity.Data
}

type container1 struct {
	serial, parity                string
	baudrate, databits, stopbits  int
	meteoAddr, fanAddr, meterAddr byte
	timeout                       time.Duration
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
		timeout:   2 * time.Second,
	}
}

// GetData gets data from sensors and represent it in Data struct.
// Modbus-RTU is a serial communication standart so I close each device's handler before opening new one
// instead of using defer statement because it makes all handler work simultaneously causing timeouts.
func (c *container1) GetData() entity.Data {
	meteoOk := true
	fanOk := true
	meterOk := true

	meteo, err := device.CWS19New(c.serial, c.meteoAddr, c.baudrate, c.databits, c.stopbits, c.parity, c.timeout)
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

	fan, err := device.ESQ760New(c.serial, c.fanAddr, c.baudrate, c.databits, c.stopbits, c.parity, c.timeout)
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

	meter, err := device.ZM194New(c.serial, c.meterAddr, c.baudrate, c.databits, c.stopbits, c.parity, c.timeout)
	if err != nil {
		meterOk = false
	}

	p1v, p2v, p3v, err := meter.GetVoltage()
	if err != nil {
		meterOk = false
	}

	p1p, p2p, p3p, err := meter.GetPower()
	if err != nil {
		meterOk = false
	}
	meter.Close()

	return entity.Data{
		Meteo: entity.Meteo{
			Ok:       meteoOk,
			Temp:     temp,
			Humidity: hum,
		},
		Fan: entity.Fan{
			Ok:        fanOk,
			Frequency: freq,
			State:     state,
			ErrorCode: errorcode,
		},
		Meter: entity.Meter{
			Ok: meterOk,
			Voltage: entity.Voltage{
				Phase1: p1v,
				Phase2: p2v,
				Phase3: p3v,
			},
			Power: entity.Power{
				Phase1: p1p,
				Phase2: p2p,
				Phase3: p3p,
			},
		},
	}
}
