package main

import (
    ezcanvas ".."
)

const MODE = ezcanvas.REPLACE

func main() {

    canvas1 := ezcanvas.NewCanvas(1024, 768)
    canvas2 := ezcanvas.NewCanvas(1024, 768)

    for n := 384; n >= 0; n-- {
        colour := int(float32(255) - (float32(n) * 255.0 / 384.0))
        canvas1.Fcircle(uint8(colour), uint8(colour), 0, MODE, 512 - ((384 - n) / 2), 384 + ((384 - n) / 2), n)
    }

    for n := 384; n >= 0; n-- {
        colour := int(float32(255) - (float32(n) * 255.0 / 384.0))
        canvas2.Fcircle(0, uint8(colour), uint8(colour), MODE, 512 + ((384 - n) / 2), 384 - ((384 - n) / 2), n)
    }

    canvas1.AddCanvas(canvas2)
    canvas1.DumpJPEG("blendedcircles.jpeg")
}
