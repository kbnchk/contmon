**Contmon** is a minimalistic web server providing json-formatted data from cryptocurrency mining container equipped with modbus-rtu controlled sensors.

At this moment it supports only my own sensors setup:

* cws19 thermal/humidity sensor
* esq760 fan frequency modulator
* zm194-d9y 3-phase electric meter

Sensors device functions implemented through [go modbus](https://github.com/goburrow/modbus).

Web server listens port 1588 and supports only http-get request to /data endpoint. It's enough for now.

Response example:

```json
{
   "Farm": "Myfarm",
   "FarmData": {
      "Container": "Container1",
      "ContainerData": {
         "Fan": {
            "ErrorCode": 0,
            "Frequency": 50,
            "Ok": true,
            "State": "Прямое вращение"
         },
         "Meteo": {
            "Humidity": 28.8,
            "Ok": true,
            "Temp": 28
         },
         "Meter": {
            "Ok": true,
            "Power": [
               826,
               747,
               657
            ],
            "Voltage": [
               228,
               227,
               228
            ]
         }
      }
   }
}
```
