// ezcanvas is a wrapper around some of Golang's "image" package, simplifying the creation of simple PNGs.

package ezcanvas

import (
    "fmt"
    "image"
    "image/color"
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
    field *image.NRGBA
    width int
    height int
}

func (c *Canvas) Get(x, y int) (r, g, b uint8) {
    if x >= 0 && x < c.width && y >= 0 && y < c.height {
        at := c.field.At(x, y)
        r, g, b := at.(color.NRGBA).R, at.(color.NRGBA).G, at.(color.NRGBA).B   // at.(color.NRGBA) is a type assertion
        return r, g, b
    }
    return 0, 0, 0
}

func (c *Canvas) SetByMode(x, y int, r, g, b uint8, mode int) {
    if mode == SET {
        c.field.Set(x, y, color.NRGBA{r, g, b, 255})        // Optimise by using the image library's built-in Set()
    } else if mode == ADD {
        c.Add(x, y, r, g, b)
    } else if mode == SUBTRACT {
        c.Subtract(x, y, r, g, b)
    } else {
        panic("unknown mode")
    }
}

func (c *Canvas) Set(x, y int, r, g, b uint8) {
    c.field.Set(x, y, color.NRGBA{r, g, b, 255})
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

    c.field.Set(x, y, color.NRGBA{new_r, new_g, new_b, 255})
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

    c.field.Set(x, y, color.NRGBA{new_r, new_g, new_b, 255})
}

func (c *Canvas) Clear(r, g, b uint8) {
    col := color.NRGBA{r, g, b, 255}

    for x := 0 ; x < c.width ; x++ {
        for y := 0 ; y < c.height ; y++ {
            c.field.Set(x, y, col)
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

func (c *Canvas) Fcircle(x, y, radius int, r, g, b uint8, mode int) {
    var pyth float64;

    for j := radius; j >= 0 ; j-- {
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
    }
    canvas.Clear(0, 0, 0)
    return &canvas
}
