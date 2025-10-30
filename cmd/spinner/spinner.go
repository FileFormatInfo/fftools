package main

import (
	"fmt"
	"math"

	"github.com/FileFormatInfo/fftools/internal"
)

var (
	solid  = [...]string{"\u2022", "\u23FA", "\u25CF", "\u2B24"}
	hollow = [...]string{"\u25E6", "\uFFEE", "\u25CB", "\u25EF", "\u2B58"}
)

type Point struct {
	X int
	Y int
}

func main() {
	for _, s := range solid {
		println("Solid spinner character:", s)
	}
	for _, h := range hollow {
		println("Hollow spinner character:", h)
	}
	oldState := internal.Init()
	defer internal.Deinit(oldState)
	width, height := internal.ScreenSize()
	println("Terminal size:", width, "x", height)

	centerX := width / 2
	centerY := height / 2

	numPoints := 16
	points := make([]Point, numPoints)
	radius := 3.0
	for i := 0; i < numPoints; i++ {
		angle := float64(i) * (360.0 / float64(numPoints)) * (3.14159 / 180.0)
		x := centerX + 2*int(math.Round(radius*math.Cos(angle)))
		y := centerY + int(math.Round(radius*math.Sin(angle)))
		points[i] = Point{X: x, Y: y}
	}

	internal.ScreenClear()

	for _, p := range points {
		internal.MoveTo(p.X, p.Y)
		//fmt.Print(solid[i%len(solid)])
		fmt.Print(solid[0])
	}
	internal.MoveTo(1, height)

}
