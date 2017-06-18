package ezcanvas

import (
    "math"
)

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
