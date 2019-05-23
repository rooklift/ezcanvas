package main

import (
    "github.com/fohristiwhirl/ezcanvas"
    "math/rand"
)

const (
    WIDTH = 1024
    HEIGHT = 768
    CIRCLES = 40000
    MIN_RADIUS = 30
    MAX_RADIUS = 60
    BRIGHTNESS = 2
)

func main() {
    virtue := ezcanvas.NewCanvas(WIDTH, HEIGHT)

    for n := 0 ; n < CIRCLES ; n++ {

        x := rand.Int31n(WIDTH)
        y := rand.Int31n(HEIGHT)

        r := rand.Int31n(BRIGHTNESS)
        g := rand.Int31n(BRIGHTNESS)
        b := rand.Int31n(BRIGHTNESS)

        rad := rand.Int31n(MAX_RADIUS - MIN_RADIUS) + MIN_RADIUS

        virtue.Fcircle(uint8(r), uint8(g), uint8(b), ezcanvas.ADD, int(x), int(y), int(rad))
    }

    virtue.DumpPNG("haze.png")
}
