package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kbnchk/contmon/internal/container"
)

func GetData(c *gin.Context) {
	cont := container.Container1()
	data := cont.GetData()
	c.IndentedJSON(http.StatusOK, data)

}
