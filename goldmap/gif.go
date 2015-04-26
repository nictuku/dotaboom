// Mostly copied from:
// github.com/SimonWaldherr/GolangSortingVisualization/blob/master/gsv.go
// Copyright 2014 Simon Waldherr
package main

/*
for f in *.gif ; do
  gifsicle --resize 320x320 -O --careful -d 5 -o sort_$f $f
  rm $f
done
*/

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
	"sort"
)

//var Max = 100
var ImgSize = 600
var Fps = 10
var Count = 10
var Mode = 1

type Visualizer interface {
	Setup(string)
	AddFrame([]int)
	Complete()
}

type GifVisualizer struct {
	name string
	g    *gif.GIF
}

func (gv *GifVisualizer) Setup(name string) {
	gv.g = &gif.GIF{
		LoopCount: 1,
	}
	gv.name = name
}

func (gv *GifVisualizer) AddFrame(heroState map[string]*heroState) {
	frame := buildImage(heroState)
	gv.g.Image = append(gv.g.Image, frame)
	gv.g.Delay = append(gv.g.Delay, 50)
}

func (gv *GifVisualizer) Complete() {
	WriteGif(gv.name, gv.g)
}

type FrameGen func(map[string]*heroState)

func (fg FrameGen) Setup(name string) {
}

func (fg FrameGen) AddFrame(heroState map[string]*heroState) {
	fg(heroState)
}

func (fg FrameGen) Complete() {
}

func buildImage(heroState map[string]*heroState) *image.Paletted {
	/*max := uint(0)
	for _, gold := range heroState {
		if gold > max {
			max = gold
		}
	}
	*/
	max := 1000
	goldPerPixel := float64(max) / float64(ImgSize)

	var frame = image.NewPaletted(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{ImgSize, ImgSize},
		},
		color.Palette{
			color.Gray{uint8(255)},
			color.Gray{uint8(0)},
		},
	)

	playerWidth := ImgSize / len(heroState)

	players := []string{}
	for k := range heroState {
		players = append(players, k)
	}
	sort.Strings(players)
	i := 0
	for _, k := range players {
		v := heroState[k]
		g := int(float64(v.gold) / goldPerPixel)
		r := image.Rect(i*playerWidth, 5, (i+1)*playerWidth, g)
		draw.Draw(frame, r, &image.Uniform{image.Black}, image.ZP, draw.Src)
		i++
		// k has player name
		fmt.Println("x, y", v.posX, v.posY)
	}

	return frame
}

func WriteGif(name string, g *gif.GIF) {
	w, err := os.Create(name + ".gif")
	if err != nil {
		fmt.Println("os.Create")
		panic(err)
	}
	fmt.Println(name+".gif created?", g)
	defer func() {
		if err := w.Close(); err != nil {
			fmt.Println("w.Close")
			panic(err)
		}
	}()
	err = gif.EncodeAll(w, g)
	if err != nil {
		fmt.Println("gif.EncodeAll")
		panic(err)
	}
}

/*
var InsertionSort = func(arr []int, frameGen FrameGen) {
	var i int
	var j int

	for i = 0; i < len(arr); i++ {
		j = i
		for j > 0 && arr[j-1] > arr[j] {
			arr[j], arr[j-1] = arr[j-1], arr[j]
			j = j - 1
			frameGen(arr)
		}
		frameGen(arr)
	}
}
/*

/*
func goldMap() {
	visualizer := &GifVisualizer{}
	visualizer.Setup("..")
	arr := gsv.RandomArray(Count, Max)
	sortFunc(arr, visualizer.AddFrame)
	visualizer.Complete()
}
*/

//func main() {
//	runSort("insertion", "insertion", InsertionSort)
//}
