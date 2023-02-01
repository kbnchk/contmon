package device

import (
	"encoding/binary"
	"math"
	"time"

	"github.com/goburrow/modbus"
)

// Star Meter CWS19 temperature and humidity sensor
type CWS19 struct {
	handler *modbus.RTUClientHandler
	client  modbus.Client
	address byte
}

func CWS19New(serial string, address byte, baudrate, databits, stopbits int, parity string) (CWS19, error) {
	handler := modbus.NewRTUClientHandler(serial)
	handler.BaudRate = baudrate
	handler.DataBits = databits
	handler.Parity = parity
	handler.StopBits = stopbits
	handler.SlaveId = address
	handler.Timeout = 2 * time.Second

	if err := handler.Connect(); err != nil {
		return CWS19{}, err
	}
	return CWS19{
		handler: handler,
		client:  modbus.NewClient(handler),
		address: address,
	}, nil
}

// Gets current temperature in Celsium
func (d *CWS19) GetTemperature() (float64, error) {
	result, err := d.client.ReadHoldingRegisters(0, 2)
	if err != nil {
		return 125, err
	}
	tempraw := binary.BigEndian.Uint16(result[:2])
	temp := float64(tempraw)/10 - 40
	temp = math.Round(temp*10) / 10
	return temp, nil
}

// Gets current humidity in %
func (d *CWS19) GetHumidity() (float64, error) {
	result, err := d.client.ReadHoldingRegisters(0, 2)
	if err != nil {
		return 0, err
	}
	humraw := binary.BigEndian.Uint16(result[2:])
	hum := float64(humraw) / 10
	hum = math.Round(hum*10) / 10
	return hum, nil
}

// Closes client handler
func (d *CWS19) Close() {
	d.handler.Close()
}
