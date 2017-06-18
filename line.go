package ezcanvas

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
