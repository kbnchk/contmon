package container

import (
	"fmt"

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
func (c *container1) GetData() (Data, error) {
	meteo, err := device.CWS19New(c.serial, c.meteoAddr, c.baudrate, c.databits, c.stopbits, c.parity)
	if err != nil {
		return Data{}, fmt.Errorf("error initializing meteo device")
	}
	temp, err := meteo.GetTemperature()
	if err != nil {
		return Data{}, fmt.Errorf("error getting temperature from meteo device")
	}
	hum, err := meteo.GetHumidity()
	if err != nil {
		return Data{}, fmt.Errorf("error getting humidity from meteo device")
	}
	meteo.Close()

	fan, err := device.ESQ760New(c.serial, c.fanAddr, c.baudrate, c.databits, c.stopbits, c.parity)
	if err != nil {
		return Data{}, fmt.Errorf("error initializing fan device")
	}
	state, err := fan.GetStatus()
	if err != nil {
		return Data{}, fmt.Errorf("error getting state from fan device")
	}
	errorcode, err := fan.GetError()
	if err != nil {
		return Data{}, fmt.Errorf("error getting error code from fan device")
	}

	freq, err := fan.GetFreq()
	if err != nil {
		return Data{}, fmt.Errorf("error getting frequency from fan device")
	}
	fan.Close()

	meter, err := device.ZM194New(c.serial, c.meterAddr, c.baudrate, c.databits, c.stopbits, c.parity)
	if err != nil {
		return Data{}, fmt.Errorf("error initializing meter device")
	}
	defer meter.Close()

	p1v, p2v, p3v, err := meter.GetVoltage()
	if err != nil {
		return Data{}, fmt.Errorf("error getting voltage from meter device")
	}

	p1p, p2p, p3p, err := meter.GetPower()
	if err != nil {
		return Data{}, fmt.Errorf("error getting voltage from meter device")
	}
	meter.Close()

	return Data{
		Meteo: Meteo{
			Temp:     temp,
			Humidity: hum,
		},
		Fan: Fan{
			Frequency: freq,
			State:     state,
			ErrorCode: errorcode,
		},
		Meter: Meter{
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
	}, nil
}
