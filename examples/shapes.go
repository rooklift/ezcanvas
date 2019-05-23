package main

import (
	"github.com/fohristiwhirl/ezcanvas"
)

const REPLACE = ezcanvas.REPLACE
const SUBTRACT = ezcanvas.SUBTRACT
const ADD = ezcanvas.ADD

func main() {
	virtue := ezcanvas.NewCanvas(640, 240)
	virtue.Clear(0, 0, 0)

	virtue.Rect(255, 255, 255, REPLACE, 20, 20, 60, 60)

	virtue.Frect(255, 0, 0, ADD, 30, 30, 70, 70)

	virtue.Fpolygon(0, 128, 0, REPLACE, 20, 120, 40, 220, 620, 120, 600, 20)
	virtue.Polygon(0, 255, 0, REPLACE, 20, 120, 40, 220, 620, 120, 600, 20)

	virtue.Fpolygon(0, 64, 0, SUBTRACT, 300, 80, 400, 80, 350, 30)
	virtue.Polygon(0, 255, 0, REPLACE, 300, 80, 400, 80, 350, 30)

	virtue.DumpPNG("shapes.png")
}
