package device

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/goburrow/modbus"
)

// ESQ-760 Frequency modulator
type ESQ760 struct {
	handler *modbus.RTUClientHandler
	client  modbus.Client
	address byte
}

func ESQ760New(serial string, address byte, baudrate, databits, stopbits int, parity string) (ESQ760, error) {
	handler := modbus.NewRTUClientHandler(serial)
	handler.BaudRate = baudrate
	handler.DataBits = databits
	handler.Parity = parity
	handler.StopBits = stopbits
	handler.SlaveId = address
	handler.Timeout = 2 * time.Second

	if err := handler.Connect(); err != nil {
		return ESQ760{}, err
	}
	return ESQ760{
		handler: handler,
		client:  modbus.NewClient(handler),
		address: address,
	}, nil
}

// Sets frequency control mode to RTU disabling manual frequency control with potentiometer or buttons
func (d *ESQ760) SetRTUControl() error {
	_, err := d.client.WriteSingleRegister(uint16(0x0002), 0)
	if err != nil {
		return err
	}
	_, err = d.client.WriteSingleRegister(uint16(0x0006), 9)
	if err != nil {
		return err
	}
	return nil
}

// Sets frequency control mode to potentiometer
func (d *ESQ760) SetManualControl() error {
	_, err := d.client.WriteSingleRegister(uint16(0x0006), 1)
	if err != nil {
		return err
	}
	return nil
}

// Gets current output frequency
func (d *ESQ760) GetFreq() (float64, error) {
	result, err := d.client.ReadInputRegisters(uint16(0x1100), 1)
	if err != nil {
		return 0, err
	}
	freq := float64(binary.BigEndian.Uint16(result)) / 100
	freq = math.Round(freq*100) / 100
	return freq, nil

}

// Gets current state
func (d *ESQ760) GetStatus() (string, error) {
	result, err := d.client.ReadHoldingRegisters(uint16(0x6000), 1)
	if err != nil {
		return "", err
	}
	value := binary.BigEndian.Uint16(result)
	switch value {
	case 1:
		return "Прямое вращение", nil
	case 2:
		return "Обратное вращение", nil
	case 3:
		return "Остановка", nil
	case 4:
		return "Авария", nil
	case 5:
		return "Пониженное напряжение", nil
	}
	return "", err
}

// Gets last error code
func (d *ESQ760) GetError() (uint8, error) {
	result, err := d.client.ReadHoldingRegisters(uint16(0x6002), 1)
	if err != nil {
		return 0, err
	}
	return uint8(binary.BigEndian.Uint16(result)), nil
}

// Sets output frequency. You need to SetRTUControl() first
func (d *ESQ760) SetFreq(v uint16) error {
	if v > 500 {
		return fmt.Errorf("value must be 0<500")
	}
	_, err := d.client.WriteSingleRegister(uint16(0x3000), v*10)
	if err != nil {
		return err
	}
	return nil
}

// Closes client handler
func (d *ESQ760) Close() {
	d.handler.Close()
}
