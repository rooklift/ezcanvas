package ezcanvas

// Some helpers to draw polygons.
// Filled polygons work correctly for convex polygons only.

type point struct {
    x int
    y int
}

type left_and_right struct {
    left int
    right int
}

type polygon struct {
    edge_points map[point]bool          // All points that make up the edges of the polygon
    extremes map[int]left_and_right     // Extreme left and right x values of the polygon for any given y value
}

func newPolygon() *polygon {
    p := polygon{}
    p.edge_points = make(map[point]bool)
    p.extremes = make(map[int]left_and_right)
    return &p
}

func (p *polygon) line(x1, y1, x2, y2 int) {

    if x1 == x2 {
        p.lineVertical(x1, y1, y2)
        return
    }
    if y1 == y2 {
        p.lineHorizontal(x1, y1, x2)
        return
    }

    dx := x1 - x2
    dy := y1 - y2
    if dx < 0 { dx *= -1 }
    if dy < 0 { dy *= -1 }

    if dy < dx {
        p.lineGentle(x1, y1, x2, y2)
    } else {
        p.lineSteep(x1, y1, x2, y2)
    }
}

func (p *polygon) lineHorizontal(x1, y, x2 int) {

    if x1 > x2 {
        x1, x2 = x2, x1
    }

    for x := x1 ; x <= x2 ; x++ {
        p.set(x, y)
    }
}

func (p *polygon) lineVertical(x, y1, y2 int) {

    if y1 > y2 {
        y1, y2 = y2, y1
    }

    for y := y1 ; y <= y2 ; y++ {
        p.set(x, y)
    }
}

func (p *polygon) lineGentle(x1, y1, x2, y2 int) {

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

        p.set(n, y1)

        the_error += dy_times_two;
        if the_error > 0 {
            y1 += additive
            the_error -= dx_times_two
        }
    }
}

func (p *polygon) lineSteep(x1, y1, x2, y2 int) {

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

        p.set(x1, n)

        the_error += dx_times_two
        if the_error > 0 {
            x1 += additive
            the_error -= dy_times_two
        }
    }
}

func (p *polygon) set(x, y int) {

    p.edge_points[point{x, y}] = true

    extremes, ok := p.extremes[y]

    if ok == false {

        p.extremes[y] = left_and_right{x, x}

    } else {

        if extremes.left > x { extremes.left = x }
        if extremes.right < x { extremes.right = x }

        p.extremes[y] = extremes
    }
}

func (p *polygon) drawEdges(canvas *Canvas, r, g, b uint8, mode int) {
    for point := range p.edge_points {
        canvas.SetByMode(r, g, b, mode, point.x, point.y)
    }
}

func (p *polygon) drawFilled(canvas *Canvas, r, g, b uint8, mode int) {
    for y, extremes := range p.extremes {
        for x := extremes.left ; x <= extremes.right ; x++ {
            canvas.SetByMode(r, g, b, mode, x, y)
        }
    }
}
