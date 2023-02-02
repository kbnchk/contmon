package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kbnchk/contmon/internal/container"
	"github.com/kbnchk/contmon/internal/server"
)

func main() {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	cont := container.Container1()
	router.GET("/data", server.GetData(cont))
	router.Run(":1588")
}
