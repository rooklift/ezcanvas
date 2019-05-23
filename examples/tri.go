package main

import (
	"math/rand"
	ezcanvas ".."
)

func main() {

	virtue := ezcanvas.NewCanvas(1024, 768)

	virtue.Clear(0, 0, 0)

	for n := 0 ; n < 750 ; n++ {

		x1 := int(rand.Int31n(1024))
		x2 := int(rand.Int31n(1024))
		x3 := int(rand.Int31n(1024))

		y1 := int(rand.Int31n(768))
		y2 := int(rand.Int31n(768))
		y3 := int(rand.Int31n(768))

		r := uint8(rand.Int31n(10))
		g := uint8(rand.Int31n(10))
		b := uint8(rand.Int31n(10))

		virtue.Fpolygon(r, g, b, ezcanvas.ADD, x1, y1, x2, y2, x3, y3)
	}

	virtue.Fpolygon(255, 255, 0, ezcanvas.REPLACE, 512, 400, 462, 350, 562, 350)
	virtue.Polygon(0, 0, 0, ezcanvas.REPLACE, 512, 400, 462, 350, 562, 350)

	virtue.DumpPNG("tri.png")
}
