package ezcanvas

import (
    "math"
)

func (c *Canvas) Fcircle(r, g, b uint8, mode int, x, y, radius int) {
    var pyth float64;

    for j := radius ; j >= 0 ; j-- {
        for i := radius ; i >= 0 ; i-- {
            pyth = math.Sqrt(math.Pow(float64(i), 2) + math.Pow(float64(j), 2));
            if (pyth < float64(radius) - 0.5) {
                c.lineHorizontal(r, g, b, mode, x - i - 1, y - j - 1, x + i)
                c.lineHorizontal(r, g, b, mode, x - i - 1, y + j, x + i)
                break
            }
        }
    }
}

func (c *Canvas) Circle(r, g, b uint8, mode int, x, y, radius int) {

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
                    c.lineHorizontal(r, g, b, mode, x - i - 1, y - j - 1, x + i)
                    c.lineHorizontal(r, g, b, mode, x - i - 1, y + j    , x + i)
                    lastiplusone = i + 1
                } else {
                    if lastiplusone == i + 1 {
                        c.SetByMode(r, g, b, mode, x - i - 1, y - j - 1)
                        c.SetByMode(r, g, b, mode, x + i    , y - j - 1)
                        c.SetByMode(r, g, b, mode, x - i - 1, y + j    )
                        c.SetByMode(r, g, b, mode, x + i    , y + j    )
                    } else {
                        c.lineHorizontal(r, g, b, mode, x - i - 1, y - j - 1, x - lastiplusone - 1)
                        c.lineHorizontal(r, g, b, mode, x + lastiplusone , y - j - 1, x + i)
                        c.lineHorizontal(r, g, b, mode, x - i - 1, y + j, x - lastiplusone - 1)
                        c.lineHorizontal(r, g, b, mode, x + lastiplusone , y + j, x + i)
                        lastiplusone = i + 1
                    }
                }
                break
            }
        }
    }
}
