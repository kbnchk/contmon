**Contmon** is a minimalistic web server providing json-formatted data from cryptocurrency mining container equipped with modbus-rtu controlled sensors.

At this moment it supports only my own sensors setup:

* cws19 thermal/humidity sensor
* esq760 fan frequency modulator
* zm194-d9y 3-phase electric meter

Device module implements only some basic functions necessary for my data structure. ESQ760 device implementation have some additional methods unused in this code (such as switching frequency control mode and remote frequncy control through modbus-rtu).

Sensors device functions implemented through [go modbus](https://github.com/goburrow/modbus) and [Gin](https://github.com/gin-gonic/gin) is used as web server framework.

Web server is listening port 1588 and supports only http-get request to /data endpoint. It's enough for now.

Response example:

`{ "Meteo": { "Temp": 5.9, "Humidity": 53.5 }, "Fan": { "Frequency": 34.51, "State": "Прямое вращение", "ErrorCode": 0 }, "Meter": { "Voltage": { "Phase1": 226, "Phase2": 225, "Phase3": 225 }, "Power": { "Phase1": 885, "Phase2": 749, "Phase3": 765 } } }`

There is no any config and any configurable parametrs at this time, all settings are hardcoded in container1 struct implementing Container interface.
