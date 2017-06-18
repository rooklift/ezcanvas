package ezcanvas

// FIXME? Not sure if I should follow the convention that
// rectangle boundaries are off-by-one.

func (c *Canvas) Rect(r, g, b uint8, mode int, x1, y1, x2, y2 int) {
    c.rect(r, g, b, mode, false, x1, y1, x2, y2)
}

func (c *Canvas) Frect(r, g, b uint8, mode int, x1, y1, x2, y2 int) {
    c.rect(r, g, b, mode, true, x1, y1, x2, y2)
}

func (c *Canvas) rect(r, g, b uint8, mode int, filled bool, x1, y1, x2, y2 int) {

    p := newPolygon()

    p.line(x1, y1, x1, y2)
    p.line(x1, y1, x2, y1)
    p.line(x1, y2, x2, y2)
    p.line(x2, y1, x2, y2)

    if filled {
        p.drawFilled(c, r, g, b, mode)
    } else {
        p.drawEdges(c, r, g, b, mode)
    }
}
