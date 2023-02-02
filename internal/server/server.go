package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kbnchk/contmon/internal/container"
)

func GetData(cont container.Container) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		data := cont.GetData()
		c.IndentedJSON(http.StatusOK, data)
	})
}
