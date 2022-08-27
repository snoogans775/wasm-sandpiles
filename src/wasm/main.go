package main

import (
	"fmt"
	"image/color"
	"math"
	"syscall/js"
	"time"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/markfarnan/go-canvas/canvas"
)

type Sandpiles struct {
	piles            []int
	size             int
	centerPileHeight int
	toppleThreshold  int
	pixelSize        int
}

var done chan struct{}

var height, width float64

const SIZE int = 39
const COLOR_MULTIPLE int = 9
const RED_VALUE int = 45
const GREEN_VALUE int = 29
const BLUE_VALUE int = 140

var piles []int = make([]int, int(math.Pow(float64(SIZE), 2)))
var sandpiles = Sandpiles{piles, SIZE, 1000000, 4, 8}
var cvs *canvas.Canvas2d

// This specifies how long a delay between calls to 'render'.     To get Frame Rate,   1s / renderDelay
var renderDelay time.Duration = 20 * time.Millisecond

func main() {
	var _, err = InitSandpiles(&sandpiles)
	if err != nil {
		fmt.Println("Failed to initialize")
	}

	//cvs, _ = canvas.NewCanvas2d(true)

	cvs, _ = canvas.NewCanvas2d(false)
	cvs.Create(int(js.Global().Get("innerWidth").Float()*0.9), int(js.Global().Get("innerHeight").Float()*0.9)) // Make Canvas 90% of window size.  For testing rendering canvas smaller than full windows

	height = float64(cvs.Height())
	width = float64(cvs.Width())

	cvs.Start(30, Render)

	//go doEvery(renderDelay, Render) // Kick off the Render function as go routine as it never returns
	<-done
}

// Helper function which calls the required func (in this case 'render') every time.Duration,  Call as a go-routine to prevent blocking, as this never returns
func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

// Called from the 'requestAnimationFrame' function.   It may also be called seperatly from a 'doEvery' function, if the user prefers drawing to be seperate from the animationFrame callback
func Render(gc *draw2dimg.GraphicContext) bool {
	Update(&sandpiles)

	for i, p := range sandpiles.piles {
		// fmt.Println(uint8(p[row][col] % 255))
		// draws square if the value is high
		gc.SetFillColor(color.RGBA{
			uint8(math.Min(float64(RED_VALUE+p*sandpiles.toppleThreshold*COLOR_MULTIPLE), 0xff)),
			uint8(math.Min(float64(0x22+p*sandpiles.toppleThreshold*COLOR_MULTIPLE), 0xff)),
			uint8(math.Min(float64(0x88+p*sandpiles.toppleThreshold*COLOR_MULTIPLE), 0xff)),
			0xff,
		})
		gc.BeginPath()

		// map the 1 dimensional slice of points to the 2d canvas
		xPosition := i % sandpiles.size
		yPosition := i / sandpiles.size
		draw2dkit.Rectangle(gc,
			float64(xPosition*sandpiles.pixelSize),
			float64(yPosition*sandpiles.pixelSize),
			float64(xPosition*sandpiles.pixelSize+sandpiles.pixelSize),
			float64(yPosition*sandpiles.pixelSize+sandpiles.pixelSize),
		)
		gc.FillStroke()
		gc.Close()

	}

	return true
}

func InitSandpiles(s *Sandpiles) (bool, error) {
	s.piles[int(math.Floor(float64(len(s.piles))/2))] = s.centerPileHeight
	return true, nil
}

func Update(s *Sandpiles) {
	for i, p := range s.piles {
		if p >= s.toppleThreshold {
			// init this pile
			s.piles[i] -= s.toppleThreshold

			//update cardinal neighbors
			if i%s.size < s.size-1 {
				s.piles[i+1] += 1
			}
			if i+s.size < len(s.piles) {
				s.piles[i+s.size] += 1
			}
			if i%s.size > 0 {
				s.piles[i-1] += 1
			}
			if i > s.size {
				s.piles[i-s.size] += 1
			}

			//update diagonal neighbors
			// FIXME: This is still a copy of the cardinal function
			// Implement with multiples of 12 to give 2 to cardinal and 1 to diagonals
			// if i%s.size < s.size-1 {
			// 	s.piles[i+1] += int(math.Floor(float64(s.toppleThreshold) / 6))
			// }
			// if i%s.size > 0 {
			// 	s.piles[i-1] += int(math.Floor(float64(s.toppleThreshold) / 6))
			// }
			// if i > s.size {
			// 	s.piles[i-s.size] += int(math.Floor(float64(s.toppleThreshold) / 6))
			// }
			// if i+s.size < len(s.piles) {
			// 	s.piles[i+s.size] += int(math.Floor(float64(s.toppleThreshold) / 6))
			// }
		}
	}
}
