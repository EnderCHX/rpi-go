package dht11

import (
	"fmt"
	"log"
	"rpi-go/adaptor"
	"time"

	"gobot.io/x/gobot/drivers/gpio"
)

var DhtPin *gpio.DirectPinDriver

func Init() error {
	r := adaptor.Adaptor
	DhtPin = gpio.NewDirectPinDriver(r, "7")
	if DhtPin == nil {
		return fmt.Errorf("failed to open DHT11 pin")
	}
	return nil
}

func Read() (float64, float64, error) {
	return readDHT11(DhtPin)
}

func readDHT11(gpioPin *gpio.DirectPinDriver) (float64, float64, error) {
	// 初始化数据读取引脚
	gpioPin.DigitalWrite(1) // 发送启动信号
	time.Sleep(18 * time.Millisecond)
	gpioPin.DigitalWrite(0)

	// 延时接收数据
	time.Sleep(40 * time.Millisecond)

	// 等待返回信号
	read, err := gpioPin.DigitalRead()
	if err != nil {
		log.Printf("failed to read DHT11 response: %v", err)
		return 0, 0, err
	}
	if read != 1 {
		return 0, 0, fmt.Errorf("failed to get DHT11 response")
	}

	// 读取数据 - DHT11 发送 40 位数据（5字节）
	var data [5]byte
	for i := 0; i < 5; i++ {
		for j := 7; j >= 0; j-- {
			read, err = gpioPin.DigitalRead()
			if err != nil {
				log.Printf("failed to read DHT11 data: %v", err)
				return 0, 0, err
			}
			if read == 1 {
				data[i] |= (1 << j)
			}
			time.Sleep(1 * time.Microsecond)
		}
	}

	// 校验数据 - DHT11 校验和
	checksum := data[0] + data[1] + data[2] + data[3]
	if data[4] != checksum {
		return 0, 0, fmt.Errorf("DHT11 checksum error")
	}

	// 提取湿度和温度
	humidity := float64(data[0])
	temperature := float64(data[2])

	return humidity, temperature, nil
}
