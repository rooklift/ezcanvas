package main

// Buddhabrot. We iterate through the Mandelbrot function; particles which escape
// have their escape route (i.e. all points they touched) brightened slightly.

import (
    "fmt"
    "github.com/fohristiwhirl/ezcanvas"
    "math"
    "math/cmplx"
    "math/rand"
    "sync"
)

const (
    OUTFILE = "buddhabrot.png"
    WIDTH = 1920
    HEIGHT = 1080
    ZOOM = 400
    X_OFFSET = 1000
    MAX_ITERATIONS = 1000
    SOURCES = 10000000
    THREADS = 4
)

var done_chan chan bool = make(chan bool)
var progress_chan chan int = make(chan int)
var threads_running int
var mutex sync.Mutex

var virtue *ezcanvas.Canvas = ezcanvas.NewCanvas(WIDTH, HEIGHT)


func main() {

    go progress_bar(SOURCES)

    for n := 1 ; n < SOURCES ; n++ {

        if n % (SOURCES / 100) == 0 {
            progress_chan <- n              // Send our progress to the progress bar drawer
        }

        var c complex128 = complex(rand.Float64() * 3.32 - 2.11, rand.Float64() * 2.48 - 1.24)

        if threads_running >= THREADS {
            <- done_chan                    // Wait till a thread ends before starting a new one
            go iterator(c)
        } else {
            go iterator(c)
            threads_running += 1
        }
    }

    for threads_running > 0 {               // Wait till all threads have finished
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


func iterator(c complex128) {

    var z complex128 = c
    var list [MAX_ITERATIONS]complex128

    for n := 0 ; n < MAX_ITERATIONS ; n++ {

        z = z * z + c
        list[n] = z

        if cmplx.Abs(z) > 2 {               // The particle does escape, so draw the list of points
            for i := 0 ; i <= n ; i++ {
                pixel_x := int(math.Floor(real(list[i]) * ZOOM)) + X_OFFSET         // It's OK if x,y is out of bounds
                pixel_y := int(math.Floor(imag(list[i]) * ZOOM)) + HEIGHT / 2
                mutex.Lock()
                virtue.Add(1, 1, 1, pixel_x, pixel_y)
                virtue.Add(1, 1, 1, pixel_x, HEIGHT - pixel_y)
                mutex.Unlock()
            }
            break
        }
    }
    done_chan <- true
}
