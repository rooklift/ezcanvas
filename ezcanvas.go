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

type Canvas struct {
    field *image.NRGBA
    width int
    height int
}

func (c *Canvas) Set(x, y int, r, g, b uint8) {
    c.field.Set(x, y, color.NRGBA{r, g, b, 255})
}

func (c *Canvas) Get(x, y int) (r, g, b uint8) {
    if x >= 0 && x < c.width && y >= 0 && y < c.height {
        at := c.field.At(x, y)
        r, g, b := at.(color.NRGBA).R, at.(color.NRGBA).G, at.(color.NRGBA).B   // at.(color.NRGBA) is a type assertion
        return r, g, b
    }
    return 0, 0, 0
}

func (c *Canvas) Clear(r, g, b uint8) {
    col := color.NRGBA{r, g, b, 255}

    for x := 0 ; x < c.width ; x++ {
        for y := 0 ; y < c.height ; y++ {
            c.field.Set(x, y, col)
        }
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

func (c *Canvas) Frect(left, top, right, bottom int, r, g, b uint8, mode string) {

    if left > right {
        left, right = right, left
    }

    if top > bottom {
        top, bottom = bottom, top
    }

    if mode == "add" {
        for x := left ; x < right ; x++ {
            for y := top ; y < bottom ; y++ {
                c.Add(x, y, r, g, b)
            }
        }
    } else if mode == "subtract" {
        for x := left ; x < right ; x++ {
            for y := top ; y < bottom ; y++ {
                c.Subtract(x, y, r, g, b)
            }
        }
    } else if mode == "set" {
        col := color.NRGBA{r, g, b, 255}
        for x := left ; x < right ; x++ {
            for y := top ; y < bottom ; y++ {
                c.field.Set(x, y, col)          // Optimise by using the image library's built-in Set()
            }
        }
    } else {
        panic("unknown mode")
    }
}

func (c *Canvas) Fcircle (x, y, radius int, r, g, b uint8, mode string) {
    var pyth float64;

    for j := radius; j >= 0 ; j-- {
        for i := radius ; i >= 0 ; i-- {
            pyth = math.Sqrt(math.Pow(float64(i), 2) + math.Pow(float64(j), 2));
            if (pyth < float64(radius) - 0.5) {
                c.LineHorizontal(x - i - 1, y - j - 1, x + i, r, g, b, mode)
                c.LineHorizontal(x - i - 1, y + j, x + i, r, g, b, mode)
                break
            }
        }
    }
}

func (c *Canvas) LineHorizontal(x1, y, x2 int, r, g, b uint8, mode string) {

    if x1 > x2 {
        x1, x2 = x2, x1
    }

    if mode == "set" {
        col := color.NRGBA{r, g, b, 255}
        for x := x1 ; x <= x2 ; x++ {
            c.field.Set(x, y, col)          // Optimise by using the image library's built-in Set()
        }
    } else if mode == "add" {
        for x := x1 ; x <= x2 ; x++ {
            c.Add(x, y, r, g, b)
        }
    } else if mode == "subtract" {
        for x := x1 ; x <= x2 ; x++ {
            c.Subtract(x, y, r, g, b)
        }
    } else {
        panic("unknown mode")
    }
}

func (c *Canvas) LineVertical(x, y1, y2 int, r, g, b uint8, mode string) {

    if y1 > y2 {
        y1, y2 = y2, y1
    }

    if mode == "set" {
        col := color.NRGBA{r, g, b, 255}
        for y := y1 ; y <= y2 ; y++ {
            c.field.Set(x, y, col)          // Optimise by using the image library's built-in Set()
        }
    } else if mode == "add" {
        for y := y1 ; y <= y2 ; y++ {
            c.Add(x, y, r, g, b)
        }
    } else if mode == "subtract" {
        for y := y1 ; y <= y2 ; y++ {
            c.Subtract(x, y, r, g, b)
        }
    } else {
        panic("unknown mode")
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
