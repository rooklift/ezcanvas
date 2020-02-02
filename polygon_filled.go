package ezcanvas

// NOTE: this only works properly for convex polygons.

type polygon_filled struct {
	extremes		map[int][]int			// Extreme left and right x values of the polygon for any given y value
}

func (p *polygon_filled) note_edge_point(x, y int) {

	if p.extremes[y] == nil {
		p.extremes[y] = append(p.extremes[y], x)
		p.extremes[y] = append(p.extremes[y], x)
		return
	}

	if x < p.extremes[y][0] {
		p.extremes[y][0] = x
	}

	if x > p.extremes[y][1] {
		p.extremes[y][1] = x
	}
}

func (p *polygon_filled) line(x1, y1, x2, y2 int) {
	if x1 == x2 {
		p.lineVertical(x1, y1, y2)
	} else if y1 == y2 {
		p.lineHorizontal(x1, y1, x2)
	} else {
		p.lineBresenham(x1, y1, x2, y2)
	}
}

func (p *polygon_filled) lineHorizontal(x1, y, x2 int) {
	p.note_edge_point(x1, y)
	p.note_edge_point(x2, y)
}

func (p *polygon_filled) lineVertical(x, y1, y2 int) {

	var sy int

	if y1 < y2 {
		sy = 1
	} else {
		sy = -1
	}

	for {
		p.note_edge_point(x, y1)
		if y1 == y2 {
			break
		}
		y1 += sy
	}
}

func (p *polygon_filled) lineBresenham(x1, y1, x2, y2 int) {

	// http://members.chello.at/~easyfilter/bresenham.html

	dx := abs(x2 - x1)
	dy := abs(y2 - y1) * -1

	var sx int
	var sy int

	if x1 < x2 { sx = 1 } else { sx = -1 }
	if y1 < y2 { sy = 1 } else { sy = -1 }

	err := dx + dy

	for {
		p.note_edge_point(x1, y1)
		if (x1 == x2 && y1 == y2) {
			break
		}
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

func (c *Canvas) Fpolygon(r, g, b uint8, mode int, args... int) {

	if len(args) < 6 {
		panic("Fpolygon needs 3 or more points")
	}

	if len(args) % 2 == 1 {
		panic("Fpolygon needs whole points (2 args each)")
	}

	var pol polygon_filled
	pol.extremes = make(map[int][]int)

/*
	// This could be handy info at some point...

	pol.min_x = args[0]
	pol.max_x = args[0]
	pol.min_y = args[1]
	pol.max_y = args[1]

	for x := args[2]; x < len(args); x += 2 {
		if x < pol.min_x {
			pol.min_x = x
		}
		if x > pol.max_x {
			pol.max_x = x
		}
	}

	for y := args[3]; y < len(args); y += 2 {
		if y < pol.min_y {
			pol.min_y = y
		}
		if y > pol.max_y {
			pol.max_y = y
		}
	}
*/

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

	for y, extremes := range pol.extremes {
		for x := extremes[0]; x <= extremes[1]; x++ {
			c.SetByMode(r, g, b, mode, x, y)
		}
	}
}
