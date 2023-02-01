package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kbnchk/contmon/internal/container"
)

func GetData(c *gin.Context) {

	cont := container.Container1()
	data, err := cont.GetData()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())

	} else {
		c.IndentedJSON(http.StatusOK, data)
	}
}
