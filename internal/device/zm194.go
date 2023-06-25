package device

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/goburrow/modbus"
	"github.com/kbnchk/contmon/internal/proto"
)

// ZIZM zm194-D9Y 3-phase electric meter
type zm194 struct {
	handler *modbus.RTUClientHandler
	client  modbus.Client
}

// Creates new ZIZM zm194-D9Y 3-phase electric meter
func ZM194(c proto.RTUConfig, address byte) Device {
	handler := proto.NewRTUHadler(c)
	handler.SlaveId = address
	return zm194{
		handler: handler,
		client:  modbus.NewClient(handler),
	}
}

func (d zm194) GetData() map[string]interface{} {
	result := make(map[string]interface{})
	data := make(map[string]interface{})

	if err := d.handler.Connect(); err != nil {
		data["Errors"] = fmt.Sprintf("ZM194 error connecting device = %v", err)
		data["Ok"] = false
	} else {
		defer d.handler.Close()

		errors := make([]string, 0, 3)

		power, err := d.GetPower()
		if err != nil {
			errors = append(errors, fmt.Sprintf("ZM194 error getting power data = %v", err))
		} else {
			data["Power"] = power
		}

		voltage, err := d.GetVoltage()
		if err != nil {
			errors = append(errors, fmt.Sprintf("ZM194 error getting voltage data = %v", err))
		} else {
			data["Voltage"] = voltage
		}

		if len(errors) == 0 {
			data["Ok"] = true
		} else {
			data["Errors"] = strings.Join(errors, ",\n")
			data["Ok"] = false
		}

	}

	result["Meter"] = data
	return result
}

// Gets current voltage
func (d *zm194) GetVoltage() (voltage [3]uint16, err error) {
	result, err := d.client.ReadHoldingRegisters(uint16(0x0000), 2)
	if err != nil {
		return
	}
	voltage[0] = uint16(binary.BigEndian.Uint32(result) / 1000)

	result, err = d.client.ReadHoldingRegisters(uint16(0x0002), 2)
	if err != nil {
		return
	}
	voltage[1] = uint16(binary.BigEndian.Uint32(result) / 1000)

	result, err = d.client.ReadHoldingRegisters(uint16(0x0004), 2)
	if err != nil {
		return
	}
	voltage[2] = uint16(binary.BigEndian.Uint32(result) / 1000)
	return
}

// Gets power
func (d *zm194) GetPower() (power [3]uint16, err error) {
	result, err := d.client.ReadHoldingRegisters(uint16(0x0012), 2)
	if err != nil {
		return
	}
	power[0] = uint16(binary.BigEndian.Uint32(result))

	result, err = d.client.ReadHoldingRegisters(uint16(0x0014), 2)
	if err != nil {
		return
	}
	power[1] = uint16(binary.BigEndian.Uint32(result))
	result, err = d.client.ReadHoldingRegisters(uint16(0x0016), 2)
	if err != nil {
		return
	}
	power[2] = uint16(binary.BigEndian.Uint32(result))
	return
}

// Closes client handler
func (d *zm194) Close() {
	d.handler.Close()
}
