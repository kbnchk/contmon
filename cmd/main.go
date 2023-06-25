package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/kbnchk/contmon/internal/device"
	"github.com/kbnchk/contmon/internal/farm"
	"github.com/kbnchk/contmon/internal/proto"
)

func main() {

	rtuConf := proto.RTUConfig{
		Serial:   "/dev/ttyUSB0",
		Baudrate: 9600,
		Databits: 8,
		Stopbits: 1,
		Parity:   "N",
		Timeout:  2 * time.Second,
	}
	meteo := device.CWS19(rtuConf, 1)
	fan := device.ESQ760(rtuConf, 2)
	meter := device.ZM194(rtuConf, 101)

	container := farm.ContainerNew("Container1", meteo, fan, meter)

	myfarm := farm.FarmNew("Myfarm", container)

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		data, err := json.MarshalIndent(myfarm.GetInfo(), "", "   ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(data)
		}
	})

	err := http.ListenAndServe(":1588", nil)
	if err != nil {
		log.Fatal(err)
	}

}
