package render

import (
	"math"

	"github.com/AchrafSoltani/glow"
)

// DrawDice renders two dice at the given position.
func DrawDice(canvas *glow.Canvas, die1, die2 int, x, y int, rolling bool, animTimer float64) {
	dieSize := 36
	gap := 12

	if rolling {
		// Randomise displayed values during animation
		d1 := int(math.Abs(math.Sin(animTimer*17.3))*5) + 1
		d2 := int(math.Abs(math.Cos(animTimer*23.7))*5) + 1
		drawSingleDie(canvas, d1, x, y, dieSize)
		drawSingleDie(canvas, d2, x+dieSize+gap, y, dieSize)
	} else {
		drawSingleDie(canvas, die1, x, y, dieSize)
		drawSingleDie(canvas, die2, x+dieSize+gap, y, dieSize)
	}
}

func drawSingleDie(canvas *glow.Canvas, value int, x, y, size int) {
	// White die with rounded corners (simulated with filled rect)
	bg := glow.Color{R: 250, G: 250, B: 245}
	border := glow.Color{R: 60, G: 60, B: 60}
	dotCol := glow.Color{R: 30, G: 30, B: 30}

	canvas.DrawRect(x, y, size, size, bg)
	canvas.DrawRectOutline(x, y, size, size, border)

	// Dot positions (relative to die)
	dotR := 3
	cx := x + size/2
	cy := y + size/2
	q := size / 4 // quarter offset

	// Standard die face patterns
	switch value {
	case 1:
		canvas.FillCircle(cx, cy, dotR, dotCol)
	case 2:
		canvas.FillCircle(cx-q, cy-q, dotR, dotCol)
		canvas.FillCircle(cx+q, cy+q, dotR, dotCol)
	case 3:
		canvas.FillCircle(cx-q, cy-q, dotR, dotCol)
		canvas.FillCircle(cx, cy, dotR, dotCol)
		canvas.FillCircle(cx+q, cy+q, dotR, dotCol)
	case 4:
		canvas.FillCircle(cx-q, cy-q, dotR, dotCol)
		canvas.FillCircle(cx+q, cy-q, dotR, dotCol)
		canvas.FillCircle(cx-q, cy+q, dotR, dotCol)
		canvas.FillCircle(cx+q, cy+q, dotR, dotCol)
	case 5:
		canvas.FillCircle(cx-q, cy-q, dotR, dotCol)
		canvas.FillCircle(cx+q, cy-q, dotR, dotCol)
		canvas.FillCircle(cx, cy, dotR, dotCol)
		canvas.FillCircle(cx-q, cy+q, dotR, dotCol)
		canvas.FillCircle(cx+q, cy+q, dotR, dotCol)
	case 6:
		canvas.FillCircle(cx-q, cy-q, dotR, dotCol)
		canvas.FillCircle(cx+q, cy-q, dotR, dotCol)
		canvas.FillCircle(cx-q, cy, dotR, dotCol)
		canvas.FillCircle(cx+q, cy, dotR, dotCol)
		canvas.FillCircle(cx-q, cy+q, dotR, dotCol)
		canvas.FillCircle(cx+q, cy+q, dotR, dotCol)
	}
}
