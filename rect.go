package ezcanvas


// FIXME? Not sure if I should follow the convention that
// rectangle boundaries are off-by-one.


func (c *Canvas) Rect(x1, y1, x2, y2 int, r, g, b uint8, mode int) {
    c.rect(x1, y1, x2, y2, r, g, b, mode, false)
}


func (c *Canvas) Frect(x1, y1, x2, y2 int, r, g, b uint8, mode int) {
    c.rect(x1, y1, x2, y2, r, g, b, mode, true)
}


func (c *Canvas) rect(x1, y1, x2, y2 int, r, g, b uint8, mode int, filled bool) {

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
