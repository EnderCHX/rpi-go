package common

import (
	"log"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

type MyFace struct {
	Face         font.Face
	ParseFont    *sfnt.Font
	FaceFilePATH string
	FontSize     float64
	DPI          float64
	Hinting      font.Hinting
}

func (f *MyFace) GetFace() {
	fontBytes, err := os.ReadFile(f.FaceFilePATH)
	if err != nil {
		log.Println(err)
	}
	f.ParseFont, err = opentype.Parse(fontBytes)
	if err != nil {
		log.Println(err)
	}
	f.Face, err = opentype.NewFace(f.ParseFont, &opentype.FaceOptions{
		Size:    24,
		DPI:     80,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Println(err)
	}
}

func (f *MyFace) FaceSizeAndDPI(size float64, dpi float64) font.Face {
	face, err := opentype.NewFace(f.ParseFont, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: f.Hinting,
	})
	if err != nil {
		log.Println(err)
	}

	return face
}

var MyMyFace MyFace

func Init() {

	MyMyFace.FaceFilePATH = "./fonts/Minecraft/类像素字体_俐方体11号.ttf"

	MyMyFace.GetFace()
}
