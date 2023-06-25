package farm

type Farm struct {
	Name       string
	Containers []Container
}

func FarmNew(name string, containers ...Container) Farm {
	return Farm{
		Name:       name,
		Containers: containers,
	}
}

func (f *Farm) GetInfo() map[string]interface{} {
	result := make(map[string]interface{})
	farmData := make(map[string]interface{})
	for _, cont := range f.Containers {
		contdata := cont.Data()
		for key, value := range contdata {
			farmData[key] = value
		}
	}
	result["FarmData"] = farmData
	result["Farm"] = f.Name
	return result
}
