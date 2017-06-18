package ezcanvas

func (c *Canvas) Polygon(r, g, b uint8, mode int, args... int) {
    c.draw_polygon(r, g, b, mode, false, args...)
}

func (c *Canvas) Fpolygon(r, g, b uint8, mode int, args... int) {
    c.draw_polygon(r, g, b, mode, true, args...)
}

// Filled polygons work correctly for convex polygons only.

func (c *Canvas) draw_polygon(r, g, b uint8, mode int, filled bool, args... int) {

    if len(args) < 4 {
        return
    }

    pol := newPolygon()

    last_x := args[0]
    last_y := args[1]

    for n := 2 ; n < len(args) - 1 ; n += 2 {

        next_x := args[n]
        next_y := args[n + 1]

        pol.line(last_x, last_y, next_x, next_y)

        last_x = next_x
        last_y = next_y
    }

    pol.line(last_x, last_y, args[0], args[1])

    if filled {
        pol.drawFilled(c, r, g, b, mode)
    } else {
        pol.drawEdges(c, r, g, b, mode)
    }
}
