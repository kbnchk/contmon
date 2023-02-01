package device

import (
	"encoding/binary"
	"time"

	"github.com/goburrow/modbus"
)

// ZIZM ZM194-D9Y 3-phase electric meter
type ZM194 struct {
	handler *modbus.RTUClientHandler
	client  modbus.Client
	address byte
}

func ZM194New(serial string, address byte, baudrate, databits, stopbits int, parity string) (ZM194, error) {
	handler := modbus.NewRTUClientHandler(serial)
	handler.BaudRate = baudrate
	handler.DataBits = databits
	handler.Parity = parity
	handler.StopBits = stopbits
	handler.SlaveId = address
	handler.Timeout = 2 * time.Second

	if err := handler.Connect(); err != nil {
		return ZM194{}, err
	}
	return ZM194{
		handler: handler,
		client:  modbus.NewClient(handler),
		address: address,
	}, nil
}

// Gets current voltage
func (d *ZM194) GetVoltage() (phase1, phase2, phase3 uint16, err error) {
	result, err := d.client.ReadHoldingRegisters(uint16(0x0000), 2)
	if err != nil {
		return 0, 0, 0, err
	}
	phase1 = uint16(binary.BigEndian.Uint32(result) / 1000)
	result, err = d.client.ReadHoldingRegisters(uint16(0x0002), 2)
	if err != nil {
		return 0, 0, 0, err
	}
	phase2 = uint16(binary.BigEndian.Uint32(result) / 1000)
	result, err = d.client.ReadHoldingRegisters(uint16(0x0004), 2)
	if err != nil {
		return 0, 0, 0, err
	}
	phase3 = uint16(binary.BigEndian.Uint32(result) / 1000)
	return
}

// Gets power
func (d *ZM194) GetPower() (phase1, phase2, phase3 uint16, err error) {
	result, err := d.client.ReadHoldingRegisters(uint16(0x0012), 2)
	if err != nil {
		return 0, 0, 0, err
	}
	phase1 = uint16(binary.BigEndian.Uint32(result))
	result, err = d.client.ReadHoldingRegisters(uint16(0x0014), 2)
	if err != nil {
		return 0, 0, 0, err
	}
	phase2 = uint16(binary.BigEndian.Uint32(result))
	result, err = d.client.ReadHoldingRegisters(uint16(0x0016), 2)
	if err != nil {
		return 0, 0, 0, err
	}
	phase3 = uint16(binary.BigEndian.Uint32(result))
	return
}

// Closes client handler
func (d *ZM194) Close() {
	d.handler.Close()
}
