package ezcanvas

// FIXME? Not sure if I should follow the convention that
// rectangle boundaries are off-by-one.

func (c *Canvas) Rect(r, g, b uint8, mode int, x1, y1, x2, y2 int) {

    if x1 > x2 { x1, x2 = x2, x1 }
    if y1 > y2 { y1, y2 = y2, y1 }

    c.Polygon(r, g, b, mode, x1, y1, x2 - 1, y1, x2 - 1, y2 - 1, x1, y2 - 1)
}

func (c *Canvas) Frect(r, g, b uint8, mode int, x1, y1, x2, y2 int) {

    if x1 > x2 { x1, x2 = x2, x1 }
    if y1 > y2 { y1, y2 = y2, y1 }

    for x := x1; x < x2; x++ {
        for y := y1; y < y2; y++ {
            c.SetByMode(r, g, b, mode, x, y)
        }
    }
}
