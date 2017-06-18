package ezcanvas


func (c *Canvas) Triangle(x1, y1, x2, y2, x3, y3 int, r, g, b uint8, mode int) {
    c.triangle(x1, y1, x2, y2, x3, y3, r, g, b, mode, false)
}


func (c *Canvas) Ftriangle(x1, y1, x2, y2, x3, y3 int, r, g, b uint8, mode int) {
    c.triangle(x1, y1, x2, y2, x3, y3, r, g, b, mode, true)
}


func (c *Canvas) triangle(x1, y1, x2, y2, x3, y3 int, r, g, b uint8, mode int, filled bool) {

    p := newPolygon()

    p.line(x1, y1, x2, y2)
    p.line(x1, y1, x3, y3)
    p.line(x2, y2, x3, y3)

    if filled {
        p.drawFilled(c, r, g, b, mode)
    } else {
        p.drawEdges(c, r, g, b, mode)
    }
}
