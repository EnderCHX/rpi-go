package camera

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os/exec"
	"rpi-go/common"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	frame image.Image
	// frameLock sync.Mutex
)

func Init() {
	log.Println("Starting camera")
	go GetFrame()
}

func GetFrame() {
	cmd := exec.Command("libcamera-vid",
		"-t", "0",
		"--codec", "mjpeg",
		"--denoise", "cdn_hq",
		"--sharpness", "1.0",
		"--inline",
		"-o", "-")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("failed to create stdout pipe: ", err)
	}

	err = cmd.Start()
	if err != nil {
		log.Println("Error starting command:", err)
		panic(err)
	}

	buffer := make([]byte, 1024)
	var frameBuffer bytes.Buffer

	for {
		// 从 stdout 中读取数据
		n, err := stdout.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("Error reading stdout:", err)
			break
		}

		// 将数据写入帧缓冲区
		frameBuffer.Write(buffer[:n])

		// 查找 JPEG 图像结束标志 (0xFFD9)
		data := frameBuffer.Bytes()
		endIdx := bytes.Index(data, []byte{0xFF, 0xD9})
		if endIdx != -1 {
			// 提取完整的 JPEG 数据
			jpegData := data[:endIdx+2]

			// 解码 JPEG 数据为 image.Image
			img, err := jpeg.Decode(bytes.NewReader(jpegData))
			if err != nil {
				log.Println("Error decoding JPEG data:", err)
				continue
			}
			img = imaging.Rotate180(img)
			rgba := image.NewRGBA(image.Rect(0, 0, 640, 480))
			draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
			common.AddLabelColor(rgba, 640-140, 20, "https://chxc.cc", color.White)
			common.AddLabelColor(rgba, 10, 480-20, "Date: "+time.Now().Format("2006-01-02 15:04:05"), color.White)

			// 设置帧
			// frameLock.Lock()
			frame = rgba
			// frameLock.Unlock()

			// 移除已处理的 JPEG 数据
			frameBuffer.Next(endIdx + 2)
		}
	}
	cmd.Wait()
	log.Println("Camera exited")
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源的 WebSocket 连接
	},
}

func WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()
	jpegData := new(bytes.Buffer)
	for {
		// frameLock.Lock()
		img := frame
		// frameLock.Unlock()

		err = jpeg.Encode(jpegData, img, nil)
		if err != nil {
			log.Println("Error encoding JPEG data:", err)
			continue
		}
		err := conn.WriteMessage(websocket.BinaryMessage, jpegData.Bytes())
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
		jpegData.Reset()
		time.Sleep(10 * time.Millisecond)
	}
}

func MjpegStreamHander(c *gin.Context) {
	c.Header("Content-Type", "multipart/x-mixed-replace; boundary=--frame")
	timeout := time.Duration(time.Minute)
	jpegData := new(bytes.Buffer)
LOOP:
	for {

		select {
		case <-time.After(timeout):
			break LOOP
		default:
		}

		img := frame

		err := jpeg.Encode(jpegData, img, nil)
		if err != nil {
			log.Println("Error encoding JPEG data:", err)
			continue
		}

		// 发送JPEG图像作为MJPEG的一帧
		c.Writer.Write([]byte("--frame\r\n"))
		c.Writer.Write([]byte("Content-Type: image/jpeg\r\n"))
		c.Writer.Write([]byte("Content-Length: "))
		c.Writer.Write([]byte(strconv.Itoa(len(jpegData.Bytes()))))
		c.Writer.Write([]byte("\r\n\r\n"))
		c.Writer.Write(jpegData.Bytes())
		c.Writer.Write([]byte("\r\n"))

		// 控制帧率，避免过于频繁发送
		jpegData.Reset()
		time.Sleep(10 * time.Millisecond)
	}
}
