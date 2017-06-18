package ezcanvas

import (
	"fmt"
)

func min_of_3(a, b, c int) int {
    if a < b && a < c {
        return a
    }
    if b < c {
        return b
    }
    return c
}

func max_of_3(a, b, c int) int {
    if a > b && a > c {
        return a
    }
    if b > c {
        return b
    }
    return c
}

func intercept_x(main_y, x1, y1, x2, y2 int) (int, error) {

	// Given a horizontal line with y = main_y, and 2 points,
	// find the x value at which they cross...

    // No intercept?

    if (y1 < main_y && y2 < main_y) || (y1 > main_y && y2 > main_y) {
        return 0, fmt.Errorf("intercept_x: no intercept")
    }

    // Entire line is identical?

    if y1 == y2 {
        return 0, fmt.Errorf("intercept_x: lines are identical")
    }

    // There is a normal intercept...

    if y2 < y1 {
        x1, y1, x2, y2 = x2, y2, x1, y1
    }

    fraction := float64(main_y - y1) / float64(y2 - y1)

    intercept_f := float64(x2 - x1) * fraction + float64(x1)

    return int(intercept_f), nil
}
