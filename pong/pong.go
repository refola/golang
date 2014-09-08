package pong

import (
	"exp/gui"
	"exp/gui/x11"
	"image/color"
	"math"
	"os"
	"time"
)

// X is from left to right and Y is from top to bottom.
type place struct {
	x, y int
}

type ball struct {
	place
	angle float64 // which way it's going
}

type player struct {
	score int
	place
}

// Lowest Y = highest position
func (p player) up(minY int) {
	if p.y > minY {
		p.y--
	}
}

// Highest Y = lowest position
func (p player) down(maxY int) {
	if p.y < maxY {
		p.y++
	}
}

type game struct {
	gui.Window
	ball
	left, right player
}

func (g *game) tick() {
	switch {
	case g.ball.y < g.right.y:
		g.right.up(g.Screen().Bounds().Min.Y)
	case g.ball.y > g.right.y:
		g.right.down(g.Screen().Bounds().Max.Y)
	}
	g.redoImage()
}

// var color = image.RGBAColor{255, 255, 255, 0} // white

func (g *game) redoImage() {
	t := int(time.Now() / 1e6)
	img := g.Screen()
	minX := img.Bounds().Min.X
	minY := img.Bounds().Min.Y
	maxX := img.Bounds().Max.X
	maxY := img.Bounds().Max.Y
	for x := minX; x < maxX; x++ {
		for y := minY; y < maxY; y++ {
			color := color.RGBA{byte((x + t) % 256), byte((y + t) % 256), byte((x + y + t) % 256), 0} // pseudorandom colors for fun
			img.Set(x, y, color)
		}
	}
	// TODO: implement
	g.FlushImage()
}

func Play() {
	// initial setup
	w, err := x11.NewWindow()
	if err != nil {
		panic("Could not make window!")
	}
	bounds := w.Screen().Bounds()
	midY := (bounds.Max.Y - bounds.Min.Y) / 2
	b := ball{place{(bounds.Max.X - bounds.Min.X) / 2, midY}, math.Pi / 6}
	left := player{0, place{bounds.Min.X, midY}}
	right := player{0, place{bounds.Max.X, midY}}
	g := &game{w, b, left, right}

	// event channel
	events := make(chan func(), 1) // a buffer size of one should be plenty

	// window events
	keys := map[int]func(){'a': func() { left.up(bounds.Min.Y) }, 'z': func() { left.down(bounds.Max.Y) }, 'q': func() { os.Exit(0) }}
	go func() {
		c := w.EventChan()
		for {
			e := <-c
			switch t := e.(type) {
			case gui.ConfigEvent:
				g.redoImage()
			case gui.ErrEvent:
				panic("Cannot handle gui.ErrEvent!")
			case gui.KeyEvent:
				if f, ok := keys[t.Key]; ok {
					events <- f
				}
			case gui.MouseEvent:
				// ignored
			default:
				panic("Unknown event type!")
			}
		}
	}()

	// time events
	go func() {
		t := time.NewTicker(500 * 1e6) // 1e6 = 1 million nanos = 1 milli
		for {
			<-t.C
			events <- func() { g.tick() }
		}
	}()

	// main loop - respond to events
	for {
		(<-events)() // Can event-driven programming get any simpler?
	}
}
