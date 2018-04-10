// ezcanvas is a wrapper around some of Golang's "image" package, simplifying the creation of simple images.
// i.e. the sort of image that can be created from simple shapes like rectangles, circles, lines.
// At the moment transparency has been removed for simplicity's sake.

package ezcanvas

import (
    "fmt"
    "image"
    "image/png"
    "os"
)

const (
    SET = 0
    REPLACE = 0

    ADD = 1

    SUBTRACT = 2
    SUB = 2
)

type Canvas struct {
    field *image.NRGBA      // https://golang.org/pkg/image/#NRGBA
    width int
    height int
    arraylen int            // is equivalent to len(c.field.Pix), the total number of r/g/b/a values, i.e. pixels * 4
}

// The image.NRGBA has a .Pix field that holds the actual pixels. From the docs:
//
//      Pix holds the image's pixels, in R, G, B, A order. The pixel at
//      (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
//
// Since we're always having the Min.X and Min.Y == 0, this simplifies to:
//
//      Pix[y * Stride + x * 4]


func (c *Canvas) Field() *image.NRGBA {
    return c.field
}

func (c *Canvas) Get(x, y int) (r, g, b uint8) {

    if x >= 0 && x < c.width && y >= 0 && y < c.height {

        index := y * c.field.Stride + x * 4
        r := c.field.Pix[index]

        index++
        g := c.field.Pix[index]

        index++
        b := c.field.Pix[index]

        return r, g, b
    }
    return 0, 0, 0
}

func (c *Canvas) SetByMode(r, g, b uint8, mode int, x, y int) {
    if mode == SET {
        c.Set(r, g, b, x, y)
    } else if mode == ADD {
        c.Add(r, g, b, x, y)
    } else if mode == SUBTRACT {
        c.Subtract(r, g, b, x, y)
    } else {
        panic("unknown mode")
    }
}

func (c *Canvas) Set(r, g, b uint8, x, y int) {

    if x >= 0 && x < c.width && y >= 0 && y < c.height {

        index := y * c.field.Stride + x * 4
        c.field.Pix[index] = r

        index++
        c.field.Pix[index] = g

        index++
        c.field.Pix[index] = b

        index++
        c.field.Pix[index] = 255        // alpha value : fully opaque
    }
}

func (c *Canvas) Add(r, g, b uint8, x, y int) {

    old_r, old_g, old_b := c.Get(x, y)

    new_r := old_r + r
    new_g := old_g + g
    new_b := old_b + b

    // In case of overflows, set to max value instead...

    if new_r < old_r { new_r = 255 }
    if new_g < old_g { new_g = 255 }
    if new_b < old_b { new_b = 255 }

    c.Set(new_r, new_g, new_b, x, y)
}

func (c *Canvas) Subtract(r, g, b uint8, x, y int) {

    old_r, old_g, old_b := c.Get(x, y)

    new_r := old_r - r
    new_g := old_g - g
    new_b := old_b - b

    // In case of underflows, set to min value instead...

    if new_r > old_r { new_r = 0 }
    if new_g > old_g { new_g = 0 }
    if new_b > old_b { new_b = 0 }

    c.Set(new_r, new_g, new_b, x, y)
}

func (c *Canvas) Clear(r, g, b uint8) {
    for x := 0 ; x < c.width ; x++ {
        for y := 0 ; y < c.height ; y++ {
            c.Set(r, g, b, x, y)
        }
    }
}

func (c *Canvas) AddCanvas(other *Canvas) {
    if c.width != other.width || c.height != other.height {
        panic("dimensions did not match for AddCanvas")
    }

    for x := 0 ; x < c.width ; x++ {
        for y := 0 ; y < c.height ; y++ {
            r, g, b := other.Get(x, y)
            c.Add(r, g, b, x, y)
        }
    }
}

func (c *Canvas) DumpPNG(filename string) error {
    outfile, err := os.Create(filename)
    if err != nil {
        if outfile != nil {
            outfile.Close()
        }
        return fmt.Errorf("Couldn't create output file '%s'", filename)
    }
    png.Encode(outfile, c.field)
    outfile.Close()
    return nil
}

func NewCanvas(width int, height int) *Canvas {
    canvas := Canvas{
                field : image.NewNRGBA(image.Rect(0, 0, width, height)),
                width : width,
                height : height,
                arraylen : width * height * 4,
    }
    canvas.Clear(0, 0, 0)
    return &canvas
}
