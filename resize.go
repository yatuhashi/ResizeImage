package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	root := "/Users/jatuha/Develop/go/MyCode/Image/"

	listFiles(root, root)
}

func listFiles(rootPath, searchPath string) {
	fis, err := ioutil.ReadDir(searchPath)

	if err != nil {
		panic(err)
	}

	for _, fi := range fis {
		fullPath := filepath.Join(searchPath, fi.Name())

		if fi.IsDir() {
			listFiles(rootPath, fullPath)
		} else {
			rel, err := filepath.Rel(rootPath, fullPath)

			if err != nil {
				panic(err)
			}

			pos := strings.LastIndex(rel, ".")
			if pos == -1 {
				continue
			}
			kaku := rel[pos:]
			if kaku == ".jpg" || kaku == ".png" {
				fmt.Println(rel, rel[:pos])
				MyResize(rel, rel[:pos]+"-s"+rel[pos:])
			}
		}
	}
}

func MyResize(filename, outname string) {
	// open "test.jpg"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	//最大、横1280縦960にしてリサイズ
	m := resize.Thumbnail(250, 250, img, resize.Lanczos3)
	middle := m.Bounds()
	middleX := middle.Size().X
	middleY := middle.Size().Y
	//白い部分を作成する
	wm := image.NewRGBA(image.Rect(0, 0, 250, 250))
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(wm, wm.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)

	if middleX < 250 {
		//余白部分の左右を片方の長さを求める
		addWhite := int((250 - middleX) / 2)
		p := image.Pt(addWhite, 0)
		//余白分のずらして合成する
		draw.Draw(wm, wm.Bounds().Add(p), m, m.Bounds().Min, draw.Over)
	}

	if middleY < 250 {
		//余白部分の左右を片方の長さを求める
		addWhite := int((250 - middleY) / 2)
		p := image.Pt(0, addWhite)
		//余白分のずらして合成する
		draw.Draw(wm, wm.Bounds().Add(p), m, m.Bounds().Min, draw.Over)
	}

	out, err := os.Create(outname)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	// write new image to file
	jpeg.Encode(out, wm, nil)
}
