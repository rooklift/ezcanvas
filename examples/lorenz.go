package main

import (
	"github.com/fohristiwhirl/ezcanvas"
)

const (
	ITERATIONS = 1500000

	WIDTH = 1800
	HEIGHT = 800

	A = 10
	B = 28
	C = 2

	MODIFIER = 0.01
	ZOOM = 18
)

func main() {

	virtue := ezcanvas.NewCanvas(WIDTH, HEIGHT)
	virtue.Clear(0,0,0)

	x := 0.2
	y := 0.1
	z := 0.1

	for n := 0 ; n < ITERATIONS ; n++ {

		nextx := x + A * (y - x) * MODIFIER
		nexty := y + ((x * (B - z)) - y) * MODIFIER
		nextz := z + ((x * y) - (C * z)) * MODIFIER

		x1 := int(z * ZOOM + WIDTH / 2)
		y1 := int(y * ZOOM + HEIGHT / 2)
		x2 := int(nextz * ZOOM + WIDTH / 2)
		y2 := int(nexty * ZOOM + HEIGHT / 2)

		virtue.Line(1, 2, 0, ezcanvas.ADD, x1, y1, x2, y2)

		x = nextx
		y = nexty
		z = nextz
	}

	virtue.DumpPNG("lorenz.png")
}
