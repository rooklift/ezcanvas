package main

import (
    ezcanvas ".."
)

const MODE = ezcanvas.REPLACE

func main() {

    virtue := ezcanvas.NewCanvas(1024, 768)

    virtue.Clear(255, 255, 255)

    for x := 12 ; x <= 1012 ; x += 50 {
        virtue.Line(0, 0, 0, MODE, 512, 384, x, 12)
        virtue.Line(0, 0, 0, MODE, 512, 384, x, 756)
    }

    virtue.Line(0, 0, 0, MODE, 12, 384, 1012, 384)

    virtue.Circle(0, 0, 0, MODE, 512, 384, 250)
    virtue.Circle(0, 0, 0, MODE, 512, 384, 300)
    virtue.Circle(0, 0, 0, MODE, 512, 384, 350)

    virtue.DumpPNG("lines.png")
}
