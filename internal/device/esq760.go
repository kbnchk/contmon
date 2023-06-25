package device

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"

	"github.com/goburrow/modbus"
	"github.com/kbnchk/contmon/internal/proto"
)

// ESQ-760 Frequency modulator
type esq760 struct {
	handler *modbus.RTUClientHandler
	client  modbus.Client
}

// Creates new ESQ-760 Frequency modulator
func ESQ760(c proto.RTUConfig, address byte) Device {
	handler := proto.NewRTUHadler(c)
	handler.SlaveId = address
	return &esq760{
		handler: handler,
		client:  modbus.NewClient(handler),
	}
}

func (d esq760) GetData() map[string]interface{} {
	result := make(map[string]interface{})
	data := make(map[string]interface{})

	if err := d.handler.Connect(); err != nil {
		data["Errors"] = fmt.Sprintf("ESQ760 error connecting device = %v", err)
		data["Ok"] = false
	} else {
		errors := make([]string, 0, 3)

		errcode, err := d.GetError()
		if err != nil {
			errors = append(errors, fmt.Sprintf("ESQ760 error getting errcode = %v", err))
		} else {
			data["ErrorCode"] = errcode
		}

		state, err := d.GetStatus()
		if err != nil {
			errors = append(errors, fmt.Sprintf("ESQ760 error getting status = %v", err))
		} else {
			data["State"] = state
		}

		frequency, err := d.GetFreq()
		if err != nil {
			errors = append(errors, fmt.Sprintf("ESQ760 error getting frequency = %v", err))
		} else {
			data["Frequency"] = frequency
		}

		if len(errors) == 0 {
			data["Ok"] = true
		} else {
			data["Errors"] = strings.Join(errors, ",\n")
			data["Ok"] = false
		}
	}

	result["Fan"] = data
	return result
}

// Sets frequency control mode to RTU disabling manual frequency control with potentiometer or buttons
func (d *esq760) SetRTUControl() error {
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

// Gets current output frequency
func (d *esq760) GetFreq() (float64, error) {
	result, err := d.client.ReadInputRegisters(uint16(0x1100), 1)
	if err != nil {
		return 0, err
	}
	freq := float64(binary.BigEndian.Uint16(result)) / 100
	freq = math.Round(freq*100) / 100
	return freq, nil

}

// Gets current state
func (d *esq760) GetStatus() (string, error) {
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
func (d *esq760) GetError() (uint8, error) {
	result, err := d.client.ReadHoldingRegisters(uint16(0x6002), 1)
	if err != nil {
		return 0, err
	}
	return uint8(binary.BigEndian.Uint16(result)), nil
}

// Sets frequency control mode to potentiometer
func (d *esq760) SetManualControl() error {
	_, err := d.client.WriteSingleRegister(uint16(0x0006), 1)
	if err != nil {
		return err
	}
	return nil
}

// Sets output frequency. You need to SetRTUControl() first
func (d *esq760) SetFreq(v uint16) error {
	if v > 500 {
		return fmt.Errorf("value must be 0<500")
	}
	_, err := d.client.WriteSingleRegister(uint16(0x3000), v*10)
	if err != nil {
		return err
	}
	return nil
}
