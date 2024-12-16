package oled1306

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"log"
	"net"
	"os"
	"rpi-go/adaptor"
	"rpi-go/api/Amap"
	"rpi-go/common"
	"sync"
	"time"

	"github.com/nfnt/resize"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"

	"github.com/disintegration/imaging"
)

var r *raspi.Adaptor
var Dev *i2c.SSD1306Driver

func Init() {

	r = adaptor.Adaptor

	Dev = i2c.NewSSD1306Driver(r, i2c.WithSSD1306DisplayWidth(128), i2c.WithSSD1306DisplayHeight(64))

	if Dev == nil {
		log.Fatal("failed to initialize OLED display")
	}

	Dev.Clear()

}

func StartDisplay(gifpath string) {
	wg := &sync.WaitGroup{}
	lock := &sync.Mutex{}

	wg.Add(2)
	//log.Println("run")
	go displayTimeAndWeatherAndSystemInfo(lock, wg)
	time.Sleep(time.Second * 1)
	go displayImage(lock, wg, gifpath)

	// wg.Wait()
}

func convertAndResizeAndCenter(w, h int, src image.Image) *image.Gray {
	src = resize.Thumbnail(uint(w), uint(h), src, resize.Bicubic)
	img := image.NewGray(image.Rect(0, 0, w, h))
	r := src.Bounds()
	r = r.Add(image.Point{(w - r.Max.X) / 2, (h - r.Max.Y) / 2})
	draw.Draw(img, r, src, image.Point{}, draw.Src)
	return img
}

func turnGif(path string) ([]*image.Gray, *gif.GIF) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	g, err := gif.DecodeAll(f)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	imgs := make([]*image.Gray, len(g.Image))
	for i := range g.Image {
		imgs[i] = convertAndResizeAndCenter(128, 64, g.Image[i])
	}
	return imgs, g
}

func displayTimeAndWeatherAndSystemInfo(lock *sync.Mutex, wg *sync.WaitGroup) {
	//log.Println("displayTimeAndWeatherAndSystemInfo1")
	defer log.Println("Exiting displayTimeAndWeatherAndSystemInfo")
	defer wg.Done()
	meminfo := common.MemoryInfo()
	loadavg := common.LoadAvg()
	addripv4 := "127.0.0.1"
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				addripv4 = ipnet.IP.String()
				break
			}
		}
	}
	go common.UpdateSystemInfo(&meminfo, &loadavg)

	weather := Amap.Weather{
		Key:      "b62c6f8c8ffbe0723c1cf4424ad78737",
		CityCode: "230100",
	}
	weatherStatus := weather.GetWeather("230100")
	// log.Println(weatherStatus)
	go func(weather *Amap.Weather, weatherStatus *Amap.WeatherResponse) {
		for {
			*weatherStatus = weather.GetWeather("230100")
			time.Sleep(time.Minute * 10)
		}
	}(&weather, &weatherStatus)

	textImg := image.NewGray(image.Rect(0, 0, 256, 128))

	for {
		lock.Lock()
		// log.Println("Text获得锁")
		timeout := time.After(time.Second * 10)

	displayloop:
		for {
			select {
			case <-timeout:
				// log.Println("Text超时")
				break displayloop
			default:
				shift := func() int {
					if time.Now().Second()%10 == 0 {
						return 1
					} else {
						return 0
					}
				}
				draw.Draw(textImg, textImg.Bounds(), image.Black, image.Point{}, draw.Src)
				common.AddLabel(textImg, shift(), 22+shift(), fmt.Sprintf("%s    %s", time.Now().Format("2006-01-02"), weatherStatus.Lives[0].Weather))
				common.AddLabel(textImg, shift(), 50+shift(), fmt.Sprintf("%s        %s°C", time.Now().Format("15:04:05"), weatherStatus.Lives[0].Temperature))
				common.AddLabel(textImg, shift(), 74+shift(), fmt.Sprintf("Mem:%d/%dMB", meminfo["MemTotal"]-meminfo["MemFree"], meminfo["MemTotal"]))
				common.AddLabel(textImg, shift(), 98+shift(), fmt.Sprintf("Load:%s/%s/%s", loadavg[0], loadavg[1], loadavg[2]))
				common.AddLabel(textImg, shift(), 122+shift(), fmt.Sprintf("IP:%s", addripv4))
				rotatedimg := imaging.Rotate180(textImg)
				rotatedimg = imaging.Resize(rotatedimg, 128, 64, imaging.NearestNeighbor)
				rotatedimg = imaging.Sharpen(rotatedimg, 3.5)
				if err := Dev.ShowImage(rotatedimg); err != nil {
					log.Fatalf("failed to draw image to OLED: %v", err)
				}
				time.Sleep(time.Millisecond * 100)
			}

		}

		lock.Unlock()
		//log.Println("Text释放锁")
		time.Sleep(time.Second * 1)
	}

}

func displayImage(lock *sync.Mutex, wg *sync.WaitGroup, gifPath string) {
	defer wg.Done()
	imgs, g := turnGif(gifPath)
	i := 0

	for {
		lock.Lock()
		//log.Println("Image获得锁")

		timeout := time.After(time.Second * 10)

	displayloop:
		for ; ; i++ {
			select {
			case <-timeout:
				//log.Println("Image超时")
				break displayloop
			default:
				index := i % len(imgs)
				c := time.After(time.Duration(10*g.Delay[index]) * time.Millisecond)
				img := (imgs)[index]
				rotatedimg := imaging.Rotate180(img)
				Dev.ShowImage(rotatedimg)
				<-c
			}
		}
		lock.Unlock()
		//log.Println("Image释放锁")
		time.Sleep(time.Second * 1)
	}
}
