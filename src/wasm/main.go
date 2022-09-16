package main

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"strconv"
	"syscall/js"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/markfarnan/go-canvas/canvas"
)

type Sandpiles struct {
	piles             []int
	size              int
	centerPileHeight  int
	toppleThreshold   int
	toppleDecrementer int
	pixelSize         int
	colorMultiple     int
}

type UI struct {
	colorMultiple   int
	toppleThreshold int
}

var done chan struct{}

var height, width float64
var cvs *canvas.Canvas2d

const SIZE int = 61
const CENTER_PILE_HEIGHT int = 1000000
const TOPPLE_THRESHOLD int = 4
const PIXEL_SIZE int = 4
const COLOR_MULTIPLE int = 8
const TOPPLE_DECREMENTER = 8

var piles []int = make([]int, int(math.Pow(float64(SIZE), 2)))
var sandpiles = Sandpiles{
	piles,
	SIZE,
	CENTER_PILE_HEIGHT,
	TOPPLE_THRESHOLD,
	TOPPLE_DECREMENTER,
	PIXEL_SIZE,
	COLOR_MULTIPLE,
}
var ui = UI{
	COLOR_MULTIPLE,
	TOPPLE_THRESHOLD,
}

func main() {
	var _, err = InitSandpiles(&sandpiles)
	if err != nil {
		fmt.Println("Failed to initialize")
	}

	cvs, _ = canvas.NewCanvas2d(false)
	cvs.Create(SIZE*sandpiles.pixelSize, SIZE*sandpiles.pixelSize)

	height = float64(cvs.Height())
	width = float64(cvs.Width())

	cvs.Start(30, Render)

	<-done
}

func InitSandpiles(s *Sandpiles) (bool, error) {
	s.piles[int(math.Floor(float64(len(s.piles))/2))] = s.centerPileHeight
	return true, nil
}

// Called from the 'requestAnimationFrame' function
func Render(gc *draw2dimg.GraphicContext) bool {
	UpdateUI(&ui)
	Update(&sandpiles, &ui)

	var s *Sandpiles = &sandpiles
	for i, p := range s.piles {
		gc.SetFillColor(color.RGBA{
			uint8(0xdd + p*s.colorMultiple),
			uint8(0xbb + p*s.colorMultiple),
			uint8(0x66 + p*s.colorMultiple),
			0xff,
		})
		gc.SetStrokeColor(color.RGBA{
			uint8(0xdd + p*s.colorMultiple),
			uint8(0xbb + p*s.colorMultiple),
			uint8(0x66 + p*s.colorMultiple),
			0xff,
		})
		if p < 0 {
			gc.SetFillColor(color.RGBA{
				uint8(200),
				uint8(100),
				uint8(200),
				0xff,
			})
		}
		gc.BeginPath()

		// map the 1 dimensional slice of points to the 2d canvas
		xPosition := i % sandpiles.size
		yPosition := int(math.Floor(float64(i / sandpiles.size)))
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

func Update(s *Sandpiles, ui *UI) {
	// update mutable state
	s.colorMultiple = ui.colorMultiple
	s.toppleThreshold = ui.toppleThreshold

	// the sandpiles algorithm
	dupe := make([]int, len(s.piles))
	for i, p := range s.piles {
		if p < s.toppleThreshold {
			dupe[i] = p
		}
	}
	for i, p := range s.piles {
		if p >= s.toppleThreshold {
			dupe[i] += p - s.toppleDecrementer
			// mutate this pile (traditionally is init to 0)
			ToppleCardinal(dupe, s, i, p)
			ToppleDiagonal(dupe, s, i, p)
		}
	}
	s.piles = dupe
	// end the sandpiles algorithm
}

func UpdateUI(ui *UI) error {
	document := js.Global().Get("document")
	if !document.Truthy() {
		return errors.New("Could not retrieve document")
	}
	// Color Multiplier
	colorSlider := document.Call("getElementById", "color-multiple")
	if !colorSlider.Truthy() {
		return errors.New("Could not retrieve color slider")
	}
	colorSliderValue, err := strconv.Atoi(colorSlider.Get("value").String())
	if err != nil {
		return errors.New("Failed to parse slider value")
	}

	// Topple Threshold
	toppleSlider := document.Call("getElementById", "topple-threshold")
	if !toppleSlider.Truthy() {
		return errors.New("Could not retrieve color slider")
	}
	toppleSliderValue, err := strconv.Atoi(toppleSlider.Get("value").String())
	if err != nil {
		return errors.New("Failed to parse slider value")
	}

	ui.colorMultiple   = colorSliderValue
	ui.toppleThreshold = toppleSliderValue

	return nil
}


func ToppleCardinal(dupe []int, s *Sandpiles, i int, p int) {
	//update cardinal neighbors
	// north
	if i > s.size {
		dupe[i-s.size] += 1
	}
	// east
	if i%s.size < s.size-1 {
		dupe[i+1] += 1
	}
	// south
	if i+s.size < len(dupe) {
		dupe[i+s.size] += 1
	}
	// west
	if i%s.size > 0 {
		dupe[i-1] += 1
	}
}

func ToppleDiagonal(dupe []int, s *Sandpiles, i int, p int) {
	//update diagonal neighbors
	// north-east
	if i%s.size < s.size-1 && i > s.size {
		dupe[i-s.size+1] += 1
	}
	// south-east
	if i+s.size < len(dupe) && i%s.size < s.size-1 {
		dupe[i+s.size+1] += 1
	}
	// south-west
	if i+s.size < len(dupe) && i%s.size > 0 {
		dupe[i+s.size-1] += 1
	}
	// north-west
	if i > s.size && i%s.size > 0 {
		dupe[i-s.size-1] += 1
	}
}
