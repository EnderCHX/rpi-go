package common

import (
	"image"
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func AddLabel(img *image.Gray, x, y int, label string) {
	col := color.White

	// 创建字体绘制器
	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 80),
		Y: fixed.Int26_6(y * 64),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  &image.Uniform{col},
		Face: MyMyFace.Face,
		Dot:  point,
	}
	d.DrawString(label)
}

func AddLabelColor(img *image.RGBA, x, y int, label string, col color.Color) {

	point := fixed.P(x, y)
	d := &font.Drawer{
		Dst:  img,
		Src:  &image.Uniform{col},
		Face: MyMyFace.FaceSizeAndDPI(20, 60),
		Dot:  point,
	}

	d.DrawString(label)
}
