package device

import (
	"encoding/binary"
	"fmt"

	"github.com/goburrow/modbus"
	"github.com/kbnchk/contmon/internal/proto"
)

// Star Meter WS19 temperature and humidity sensor
type cws19 struct {
	handler *modbus.RTUClientHandler
	client  modbus.Client
}

// creades CWS19 temperature and humidity sensor
func CWS19(c proto.RTUConfig, address byte) Device {
	handler := proto.NewRTUHadler(c)
	handler.SlaveId = address
	return &cws19{
		handler: handler,
		client:  modbus.NewClient(handler),
	}
}

// Gets all data
func (d cws19) GetData() map[string]interface{} {
	result := make(map[string]interface{})
	data := make(map[string]interface{})

	if err := d.handler.Connect(); err != nil {
		data["Errors"] = fmt.Sprintf("CWS19 error conntcting device = %v", err)
		data["Ok"] = false
	} else {
		defer d.handler.Close()
		sensordata, err := d.client.ReadHoldingRegisters(0, 2)
		if err != nil {
			data["Errors"] = fmt.Sprintf("CWS19 error reading register = %v", err)
			data["Ok"] = false
		} else {
			tempdata := binary.BigEndian.Uint16(sensordata[:2])
			//temperature sensor returns uint 0-1650 which matchs temp from -40 to +125 C
			temperature := (tempdata - 400) / 10
			humdata := binary.BigEndian.Uint16(sensordata[2:])
			//humidity sensor returns uint 0-1000 which matchs 0-100 %RH
			humidity := float64(humdata) / 10

			data["Humidity"] = humidity
			data["Temp"] = temperature
			data["Ok"] = true
		}
	}
	result["Meteo"] = data
	return result
}
