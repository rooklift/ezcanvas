// ezcanvas is a wrapper around some of Golang's "image" package, simplifying the creation of simple images.
// i.e. the sort of image that can be created from simple shapes like rectangles, circles, lines.
// At the moment transparency has been removed for simplicity's sake.

package ezcanvas

import (
    "fmt"
    "image"
    "image/png"
    "math"
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


func (c *Canvas) Get(x, y int) (r, g, b uint8) {

    index := y * c.field.Stride + x * 4
    if index >= 0 && index < c.arraylen - 3 {

        r := c.field.Pix[index]

        index++
        g := c.field.Pix[index]

        index++
        b := c.field.Pix[index]

        return r, g, b
    }
    return 0, 0, 0
}

func (c *Canvas) SetByMode(x, y int, r, g, b uint8, mode int) {
    if mode == SET {
        c.Set(x, y, r, g, b)
    } else if mode == ADD {
        c.Add(x, y, r, g, b)
    } else if mode == SUBTRACT {
        c.Subtract(x, y, r, g, b)
    } else {
        panic("unknown mode")
    }
}

func (c *Canvas) Set(x, y int, r, g, b uint8) {

    index := y * c.field.Stride + x * 4
    if index >= 0 && index < c.arraylen - 3 {

        c.field.Pix[index] = r

        index++
        c.field.Pix[index] = g

        index++
        c.field.Pix[index] = b

        index++
        c.field.Pix[index] = 255        // alpha value : fully opaque
    }
}

func (c *Canvas) Add(x, y int, r, g, b uint8) {

    new_r, new_g, new_b := c.Get(x, y)

    new_r += r
    new_g += g
    new_b += b

    // In case of overflows, set to max value instead...

    if new_r < r { new_r = 255 }
    if new_g < g { new_g = 255 }
    if new_b < b { new_b = 255 }

    c.Set(x, y, new_r, new_g, new_b)
}

func (c *Canvas) Subtract(x, y int, r, g, b uint8) {

    new_r, new_g, new_b := c.Get(x, y)

    new_r -= r
    new_g -= g
    new_b -= b

    // In case of underflows, set to min value instead...

    if new_r > r { new_r = 0 }
    if new_g > g { new_g = 0 }
    if new_b > b { new_b = 0 }

    c.Set(x, y, new_r, new_g, new_b)
}

func (c *Canvas) Clear(r, g, b uint8) {
    for x := 0 ; x < c.width ; x++ {
        for y := 0 ; y < c.height ; y++ {
            c.Set(x, y, r, g, b)
        }
    }
}

func (c *Canvas) Frect(x1, y1, x2, y2 int, r, g, b uint8, mode int) {

    if x1 > x2 {
        x1, x2 = x2, x1
    }

    if y1 > y2 {
        y1, y2 = y2, y1
    }

    for x := x1 ; x < x2 ; x++ {
        for y := y1 ; y < y2 ; y++ {
            c.SetByMode(x, y, r, g, b, mode)
        }
    }
}

func (c *Canvas) Rect(x1, y1, x2, y2 int, r, g, b uint8, mode int) {

    if x1 > x2 {
        x1, x2 = x2, x1
    }

    if y1 > y2 {
        y1, y2 = y2, y1
    }

    c.lineHorizontal(x1, y1, x2, r, g, b, mode)
    c.lineHorizontal(x1, y2, x2, r, g, b, mode)
    c.lineVertical(x1, y1, y2, r, g, b, mode)
    c.lineVertical(x2, y1, y2, r, g, b, mode)

}

func (c *Canvas) Fcircle(x, y, radius int, r, g, b uint8, mode int) {
    var pyth float64;

    for j := radius ; j >= 0 ; j-- {
        for i := radius ; i >= 0 ; i-- {
            pyth = math.Sqrt(math.Pow(float64(i), 2) + math.Pow(float64(j), 2));
            if (pyth < float64(radius) - 0.5) {
                c.lineHorizontal(x - i - 1, y - j - 1, x + i, r, g, b, mode)
                c.lineHorizontal(x - i - 1, y + j, x + i, r, g, b, mode)
                break
            }
        }
    }
}

func (c *Canvas) Circle(x, y, radius int, r, g, b uint8, mode int) {

    // I wrote this algorithm 15 years ago for C and can't remember how it works. But it does.

    var pyth float64
    var topline bool = true
    var lastiplusone int

    for j := radius - 1 ; j >= 0 ; j-- {
        for i := radius - 1 ; i >= 0 ; i-- {
            pyth = math.Sqrt(math.Pow(float64(i), 2) + math.Pow(float64(j), 2))
            if (pyth < float64(radius) - 0.5) {
                if topline {                    // i.e. if we're on the top (and, with mirroring, bottom) lines
                    topline = false
                    c.lineHorizontal(x - i - 1, y - j - 1, x + i, r, g, b, mode)
                    c.lineHorizontal(x - i - 1, y + j    , x + i, r, g, b, mode)
                    lastiplusone = i + 1
                } else {
                    if lastiplusone == i + 1 {
                        c.SetByMode(x - i - 1, y - j - 1, r, g, b, mode)
                        c.SetByMode(x + i    , y - j - 1, r, g, b, mode)
                        c.SetByMode(x - i - 1, y + j    , r, g, b, mode)
                        c.SetByMode(x + i    , y + j    , r, g, b, mode)
                    } else {
                        c.lineHorizontal(x - i - 1, y - j - 1, x - lastiplusone - 1, r, g, b, mode)
                        c.lineHorizontal(x + lastiplusone , y - j - 1, x + i, r, g, b, mode)
                        c.lineHorizontal(x - i - 1, y + j, x - lastiplusone - 1, r, g, b, mode)
                        c.lineHorizontal(x + lastiplusone , y + j, x + i, r, g, b, mode)
                        lastiplusone = i + 1
                    }
                }
                break
            }
        }
    }
}

func (c *Canvas) Line(x1, y1, x2, y2 int, r, g, b uint8, mode int) {

    if x1 == x2 {
        c.lineVertical(x1, y1, y2, r, g, b, mode)
        return
    }
    if y1 == y2 {
        c.lineHorizontal(x1, y1, x2, r, g, b, mode)
        return
    }

    dx := x1 - x2
    dy := y1 - y2
    if dx < 0 { dx *= -1 }
    if dy < 0 { dy *= -1 }

    if dy < dx {
        c.lineGentle(x1, y1, x2, y2, r, g, b, mode)
    } else {
        c.lineSteep(x1, y1, x2, y2, r, g, b, mode)
    }
}

func (c *Canvas) lineHorizontal(x1, y, x2 int, r, g, b uint8, mode int) {

    if x1 > x2 {
        x1, x2 = x2, x1
    }

    for x := x1 ; x <= x2 ; x++ {
        c.SetByMode(x, y, r, g, b, mode)
    }
}

func (c *Canvas) lineVertical(x, y1, y2 int, r, g, b uint8, mode int) {

    if y1 > y2 {
        y1, y2 = y2, y1
    }

    for y := y1 ; y <= y2 ; y++ {
        c.SetByMode(x, y, r, g, b, mode)
    }
}

func (c *Canvas) lineGentle(x1, y1, x2, y2 int, r, g, b uint8, mode int) {

    // Based on an algorithm I read on the web 15 years ago;
    // The webpage has long since vanished.

    var additive int

    if x1 > x2 {
        x1, x2 = x2, x1
        y1, y2 = y2, y1
    }

    if (y1 < y2) {
        additive = 1;
    } else {
        additive = -1;
    }

    dy_times_two := (y2 - y1) * 2
    if dy_times_two < 0 { dy_times_two *= -1 }

    dx_times_two := (x2 - x1) * 2       // We know we're going right, no need to check for < 0

    the_error := x1 - x2

    for n := x1 ; n <= x2 ; n++ {

        c.SetByMode(n, y1, r, g, b, mode)

        the_error += dy_times_two;
        if the_error > 0 {
            y1 += additive
            the_error -= dx_times_two
        }
    }
}

func (c *Canvas) lineSteep(x1, y1, x2, y2 int, r, g, b uint8, mode int) {

    var additive int

    if y1 > y2 {
        x1, x2 = x2, x1
        y1, y2 = y2, y1
    }

    if (x1 < x2) {
        additive = 1;
    } else {
        additive = -1;
    }

    dy_times_two := (y2 - y1) * 2       // We know we're going down, no need to check for < 0

    dx_times_two := (x2 - x1) * 2
    if dx_times_two < 0 { dx_times_two *= -1 }

    the_error := y1 - y2;

    for n := y1 ; n <= y2 ; n++ {

        c.SetByMode(x1, n, r, g, b, mode)

        the_error += dx_times_two
        if the_error > 0 {
            x1 += additive
            the_error -= dy_times_two
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
            c.Add(x, y, r, g, b)
        }
    }
}

func (c *Canvas) DumpPNG(filename string) error {
    outfile, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("Couldn't create output file '%s'", filename)
    }
    png.Encode(outfile, c.field)
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
