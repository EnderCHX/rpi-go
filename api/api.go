package api

import (
	"log"
	"rpi-go/device/camera"

	"github.com/gin-gonic/gin"
)

func ApiRun() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/ws/video", func(c *gin.Context) {
		log.Println("/ws/video ", "from ", c.Request.RemoteAddr)
		camera.WebSocketHandler(c)
	})

	r.GET("/video", func(c *gin.Context) {
		log.Println("/video ", "from ", c.Request.RemoteAddr)
		camera.MjpegStreamHander(c)
	})

	r.Run(":8080")
}
