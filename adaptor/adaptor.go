package adaptor

import (
	"gobot.io/x/gobot/platforms/raspi"
)

var Adaptor *raspi.Adaptor

func Init() {
	Adaptor = raspi.NewAdaptor()
}
