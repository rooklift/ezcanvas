package ezcanvas

type point struct {
	x int
	y int
}

type polygon_unfilled struct {
	points			map[point]bool
}

func (p *polygon_unfilled) line(x1, y1, x2, y2 int) {
	if x1 == x2 {
		p.lineVertical(x1, y1, y2)
	} else if y1 == y2 {
		p.lineHorizontal(x1, y1, x2)
	} else {
		p.lineBresenham(x1, y1, x2, y2)
	}
}

func (p *polygon_unfilled) lineHorizontal(x1, y, x2 int) {

	var sx int

	if x1 < x2 {
		sx = 1
	} else {
		sx = -1
	}

	for {
		p.points[point{x1, y}] = true
		if x1 == x2 {
			break
		}
		x1 += sx
	}
}

func (p *polygon_unfilled) lineVertical(x, y1, y2 int) {

	var sy int

	if y1 < y2 {
		sy = 1
	} else {
		sy = -1
	}

	for {
		p.points[point{x, y1}] = true
		if y1 == y2 {
			break
		}
		y1 += sy
	}
}

func (p *polygon_unfilled) lineBresenham(x1, y1, x2, y2 int) {

	// http://members.chello.at/~easyfilter/bresenham.html

	dx := abs(x2 - x1)
	dy := abs(y2 - y1) * -1

	var sx int
	var sy int

	if x1 < x2 { sx = 1 } else { sx = -1 }
	if y1 < y2 { sy = 1 } else { sy = -1 }

	err := dx + dy

	for {
		p.points[point{x1, y1}] = true
		if (x1 == x2 && y1 == y2) { break }
		e2 := 2 * err
		if (e2 >= dy) {
			err += dy
			x1 += sx
		}
		if (e2 <= dx) {
			err += dx
			y1 += sy
		}
	}
}

func (c *Canvas) Polygon(r, g, b uint8, mode int, args... int) {

	if len(args) < 6 {
		panic("Polygon needs 3 or more points")
	}

	if len(args) % 2 == 1 {
		panic("Polygon needs whole points (2 args each)")
	}

	var pol polygon_unfilled
	pol.points = make(map[point]bool)

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

	for point, _ := range pol.points {
		c.SetByMode(r, g, b, mode, point.x, point.y)
	}
}
