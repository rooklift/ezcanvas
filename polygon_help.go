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
	} else if y1 == y2 {
		p.lineHorizontal(x1, y1, x2)
	} else {
		p.lineBresenham(x1, y1, x2, y2)
	}
}

func (p *polygon) lineHorizontal(x1, y, x2 int) {

	var sx int

	if x1 < x2 {
		sx = 1
	} else {
		sx = -1
	}

	for {
		p.set(x1, y)
		if x1 == x2 {
			break
		}
		x1 += sx
	}
}

func (p *polygon) lineVertical(x, y1, y2 int) {

	var sy int

	if y1 < y2 {
		sy = 1
	} else {
		sy = -1
	}

	for {
		p.set(x, y1)
		if y1 == y2 {
			break
		}
		y1 += sy
	}
}

func (p *polygon) lineBresenham(x1, y1, x2, y2 int) {

	// http://members.chello.at/~easyfilter/bresenham.html

	dx := abs(x2 - x1)
	dy := abs(y2 - y1) * -1

	var sx int
	var sy int

	if x1 < x2 { sx = 1 } else { sx = -1 }
	if y1 < y2 { sy = 1 } else { sy = -1 }

	err := dx + dy

	for {
		p.set(x1, y1)
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

func (p *polygon) set(x, y int) {

	// Set a single point as being on the edge of the polygon,
	// and update the data structure.

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
