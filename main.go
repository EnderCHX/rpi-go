package main

import (
	"log"

	"rpi-go/adaptor"
	"rpi-go/api"
	"rpi-go/common"
	"rpi-go/device/camera"
	"rpi-go/device/oled1306"

	"gobot.io/x/gobot"
)

func main() {
	log.Println("启动服务器...")

	// if len(os.Args) != 2 {
	// 	log.Fatal("please provide the path to an animated GIF")
	// }

	adaptor.Init()
	oled1306.Init()
	camera.Init()
	common.Init()

	work := func() {
		oled1306.StartDisplay("/home/c/code/go/rpi-go/ballerine.gif")
	}

	bot := gobot.NewRobot("rpi-go",
		[]gobot.Connection{adaptor.Adaptor},
		[]gobot.Device{
			oled1306.Dev,
		},
		work,
	)

	go bot.Start()

	api.ApiRun()

}
