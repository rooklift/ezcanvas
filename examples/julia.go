package main

// A small modification of mandelbrot.go to make it produce Julia sets instead

import (
	"fmt"
	"math/cmplx"
	"sync"

	"github.com/fohristiwhirl/ezcanvas"
)

const (
	OUTFILE = "julia.png"
	WIDTH = 1920
	HEIGHT = 1080
	CENTRE_X = 0
	CENTRE_Y = 0
	OVERSAMPLE = 3
	ZOOM = 500
	MAX_ITERATIONS = 10000
	THREADS = 4
	JULIA_C = complex(-0.835, -0.2321)		// complex(-0.8, 0.156)
)

type colour struct {
	r int
	g int
	b int
}

var done_chan chan bool = make(chan bool)
var progress_chan chan int = make(chan int)
var threads_running int
var mutex sync.Mutex

var virtue *ezcanvas.Canvas = ezcanvas.NewCanvas(WIDTH, HEIGHT)


func main() {

	var i int

	var startx float64 = CENTRE_X - (WIDTH / 2.0) / ZOOM		// The decimal in 2.0 is annoyingly required
	var endx float64 = CENTRE_X + (WIDTH / 2.0) / ZOOM
	var starty float64 = CENTRE_Y - (HEIGHT / 2.0) / ZOOM

	dx := endx - startx
	pixel_size := dx / WIDTH
	subpixel_size := pixel_size / OVERSAMPLE

	go progress_bar(WIDTH * HEIGHT)

	x := startx

	for x_pixel := 0 ; x_pixel < WIDTH ; x_pixel++ {

		y := starty

		for y_pixel := 0 ; y_pixel < HEIGHT ; y_pixel++ {

			i++

			if i % 1000 == 0 {
				progress_chan <- i
			}

			if threads_running >= THREADS {
				<- done_chan                    // Wait till a thread ends before starting a new one
				threads_running -= 1
			}

			go paint_pixel(x_pixel, y_pixel, x, y, subpixel_size)
			threads_running += 1

			y += pixel_size

		}

		x += pixel_size
	}

	for threads_running > 0 {               	// Wait till all threads have finished
		<- done_chan
		threads_running -= 1
	}

	virtue.DumpPNG(OUTFILE)
	fmt.Printf("\n")
}

func progress_bar(iterations int) {

	// Get the current iteration from channel; and draw some sort of progress report

	var chars_to_erase int

	fmt.Printf("Working: ")

	for {
		current_iteration := <- progress_chan
		for c := 0 ; c < chars_to_erase ; c++ {
			fmt.Printf("\b")
		}
		chars_to_erase, _ = fmt.Printf("%d%%", current_iteration * 100 / iterations)
	}
}

func paint_pixel(x_pixel, y_pixel int, x, y float64, subpixel_size float64) {

	var results []colour

	for sub_x := 0 ; sub_x < OVERSAMPLE ; sub_x++ {

		for sub_y := 0 ; sub_y < OVERSAMPLE ; sub_y++ {

			actual_x := x + float64(sub_x) * subpixel_size
			actual_y := y + float64(sub_y) * subpixel_size

			var c complex128 = complex(actual_x, actual_y)

			n := iterations_till_escape(c)
			results = append(results, choose_colour(n))
		}
	}

	av := colour_average(results)

	mutex.Lock()
	virtue.Set(uint8(av.r), uint8(av.g), uint8(av.b), x_pixel, y_pixel)
	mutex.Unlock()

	done_chan <- true
}

func iterations_till_escape(c complex128) int {				// How long it takes for the particle to be obviously escaping...

	var z complex128 = c

	for n := 0 ; n < MAX_ITERATIONS ; n++ {
		z = z * z + JULIA_C
		if cmplx.Abs(z) > 2 {               				// ...for this definition of "obviously".
			return n
		}
	}

	return -1
}

func min(a, b int) int {
	if a < b { return a } else { return b }
}

func max(a, b int) int {
	if a > b { return a } else { return b }
}

func colour_average(cols []colour) colour {
	var total_r, total_g, total_b int
	for _, c := range cols {
		total_r += c.r
		total_g += c.g
		total_b += c.b
	}
	return colour{total_r / len(cols), total_g / len(cols), total_b / len(cols)}
}

func choose_colour(n int) colour {

	if n == -1 {											// No-escape
		return colour{0, 0, 0}
	}

	r := min(255, n)
	g := min(255, n * 2)
	b := min(255, n * 4)

	r = max(r, 0)
	g = max(g, 0)
	b = max(b, 0)

	return colour{r, g, b}
}
