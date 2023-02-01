package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kbnchk/contmon/server"
)

func main() {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	router.GET("/data", server.GetData)
	router.Run(":1588")
}
