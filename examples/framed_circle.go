package main

import (
    "github.com/fohristiwhirl/ezcanvas"
)

const (
    MODE = ezcanvas.REPLACE
    WIDTH = 1024
    HEIGHT = 768
)

func main() {
    virtue := ezcanvas.NewCanvas(WIDTH, HEIGHT)
    virtue.Clear(0, 0, 0)

    for n := 100 ; n < WIDTH / 3 ; n++ {
        virtue.Rect(1, 1, 1, ezcanvas.ADD, 40, 40, n, HEIGHT - 40)
    }

    for n := 100 ; n < WIDTH / 3 ; n++ {
        virtue.Rect(1, 1, 1, ezcanvas.ADD, WIDTH - 40, 40, WIDTH - n, HEIGHT - 40)
    }

    for n := 0; n < 256; n += 2 {
        virtue.Fcircle(1, 0, 0, ezcanvas.ADD, WIDTH / 2, HEIGHT / 2, n)
    }

    virtue.DumpPNG("framed_circle.png")
}
