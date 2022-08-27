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
	piles            [][]float64
	centerPileHeight float64
	toppleThreshold  float64
	pixelSize        int
}

var done chan struct{}

var height, width float64
var size int = 51
var sandpiles = Sandpiles{make([][]float64, size), 10000, 12, 6}
var cvs *canvas.Canvas2d

// This specifies how long a delay between calls to 'render'.     To get Frame Rate,   1s / renderDelay
var renderDelay time.Duration = 20 * time.Millisecond

func main() {
	var _, err = InitSandpiles(&sandpiles)
	if err != nil {
		fmt.Println("Failed to initialize")
	}

	FrameRate := time.Second / renderDelay
	println("FPS:", FrameRate)
	//cvs, _ = canvas.NewCanvas2d(true)

	cvs, _ = canvas.NewCanvas2d(false)
	cvs.Create(int(js.Global().Get("innerWidth").Float()*0.9), int(js.Global().Get("innerHeight").Float()*0.9)) // Make Canvas 90% of window size.  For testing rendering canvas smaller than full windows

	height = float64(cvs.Height())
	width = float64(cvs.Width())

	cvs.Start(20, Render)

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
	UpdateSandpiles(&sandpiles)

	for row := range sandpiles.piles {
		for col := range sandpiles.piles[row] {
			p := sandpiles.piles
			// fmt.Println(uint8(p[row][col] % 255))
			// draws square if the value is high
			gc.SetFillColor(color.RGBA{
				0x22 + uint8(p[row][col] * 20), 
				0x77 - uint8(p[row][col] * 20), 
				0x22 + uint8(p[row][col] * 20), 
				0xff,
			})
			gc.BeginPath()
			draw2dkit.Rectangle(gc,
				float64(row*sandpiles.pixelSize),
				float64(col*sandpiles.pixelSize),
				float64((row+1)*sandpiles.pixelSize),
				float64((col+1)*sandpiles.pixelSize),
			)
			gc.FillStroke()
			gc.Close()

		}
	}

	return true
}

func InitSandpiles(s *Sandpiles) (bool, error) {
	for i := range s.piles {
		s.piles[i] = make([]float64, len(s.piles))
	}
	s.piles[size/2][size/2] = s.centerPileHeight

	return true, nil
}

func UpdateSandpiles(s *Sandpiles) {
	for row := range s.piles {
		for col := range s.piles[row] {
			if s.piles[row][col] >= s.toppleThreshold {
				// init this pile
				s.piles[row][col] -= s.toppleThreshold

				//update cardinal neighbors
				if row < size-1 {
					s.piles[row+1][col] += math.Floor(float64(s.toppleThreshold) / 6)
				}
				if row > 0 {
					s.piles[row-1][col] += math.Floor(float64(s.toppleThreshold) / 6)
				}
				if col < size-1 {
					s.piles[row][col+1] += math.Floor(float64(s.toppleThreshold) / 6)
				}
				if col > 0 {
					s.piles[row][col-1] += math.Floor(float64(s.toppleThreshold) / 6)
				}

				//update diagonal neighbors
				if row < size-1 && col > 0{
					s.piles[row+1][col-1] += math.Floor(float64(s.toppleThreshold) / 12)
				}
				if row < size-1 && col < size-1 {
					s.piles[row+1][col+1] += math.Floor(float64(s.toppleThreshold) / 12)
				}
				if row > 0 && col > 0 {
					s.piles[row-1][col-1] += math.Floor(float64(s.toppleThreshold) / 12)
				}
				if row > 0 && col < size-1 {
					s.piles[row-1][col+1] += math.Floor(float64(s.toppleThreshold) / 12)
				}
			}
		}
	}
}