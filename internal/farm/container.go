package farm

import (
	"github.com/kbnchk/contmon/internal/device"
)

type Container struct {
	Name    string
	Devices []device.Device
}

func ContainerNew(name string, devs ...device.Device) Container {
	return Container{
		Name:    name,
		Devices: devs,
	}
}

func (c *Container) Data() map[string]interface{} {
	result := make(map[string]interface{})
	containerData := make(map[string]interface{})
	for _, device := range c.Devices {
		data := device.GetData()
		for key, value := range data {
			containerData[key] = value
		}
	}
	result["ContainerData"] = containerData
	result["Container"] = c.Name
	return result
}
