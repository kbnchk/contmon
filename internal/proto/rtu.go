package proto

import (
	"time"

	"github.com/goburrow/modbus"
)

type RTUConfig struct {
	Serial                       string
	Baudrate, Databits, Stopbits int
	Parity                       string
	Timeout                      time.Duration
}

func NewRTUHadler(c RTUConfig) *modbus.RTUClientHandler {
	handler := modbus.NewRTUClientHandler(c.Serial)
	handler.BaudRate = c.Baudrate
	handler.DataBits = c.Databits
	handler.Parity = c.Parity
	handler.StopBits = c.Stopbits
	handler.Timeout = c.Timeout
	return handler
}
