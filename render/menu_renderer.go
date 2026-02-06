package render

import (
	"math"

	"github.com/AchrafSoltani/glow"
)

// DrawMenuBackground draws Moroccan zellige-inspired geometric patterns.
func DrawMenuBackground(canvas *glow.Canvas, timer float64) {
	w := canvas.Width()
	h := canvas.Height()

	// Dark background
	canvas.Clear(glow.Color{R: 15, G: 30, B: 15})

	// Animated zellige pattern
	spacing := 60
	patternOffset := int(timer * 10)

	for y := -spacing; y < h+spacing; y += spacing {
		for x := -spacing; x < w+spacing; x += spacing {
			px := x + (patternOffset % spacing)
			py := y

			alpha := 0.15 + 0.05*math.Sin(timer+float64(x)*0.01+float64(y)*0.01)
			if alpha > 0.2 {
				alpha = 0.2
			}

			intensity := uint8(alpha * 255)

			// 8-pointed star pattern (Moroccan zellige motif)
			drawZelligeStar(canvas, px, py, 12, glow.Color{R: 0, G: intensity, B: intensity / 2})

			// Diamond between stars
			if (x/spacing+y/spacing)%2 == 0 {
				drawSmallDiamond(canvas, px+spacing/2, py+spacing/2, 5,
					glow.Color{R: intensity, G: intensity / 2, B: 0})
			}
		}
	}

	// Border frame
	borderW := 3
	frameColor := ZelligeGold
	for t := 0; t < borderW; t++ {
		canvas.DrawRectOutline(20+t, 20+t, w-40-2*t, h-40-2*t, frameColor)
	}

	// Corner ornaments
	drawCornerOrnament(canvas, 20, 20, 1)      // top-left
	drawCornerOrnament(canvas, w-20, 20, 2)     // top-right
	drawCornerOrnament(canvas, 20, h-20, 3)     // bottom-left
	drawCornerOrnament(canvas, w-20, h-20, 4)   // bottom-right
}

func drawZelligeStar(canvas *glow.Canvas, cx, cy, size int, color glow.Color) {
	// Draw an 8-pointed star using two overlapping squares (rotated 45 degrees)
	// First square (axis-aligned)
	halfS := size / 2
	for dy := -halfS; dy <= halfS; dy++ {
		for dx := -halfS; dx <= halfS; dx++ {
			if abs(dx)+abs(dy) <= halfS {
				canvas.SetPixel(cx+dx, cy+dy, color)
			}
		}
	}
}

func drawSmallDiamond(canvas *glow.Canvas, cx, cy, size int, color glow.Color) {
	for dy := -size; dy <= size; dy++ {
		w := size - abs(dy)
		for dx := -w; dx <= w; dx++ {
			canvas.SetPixel(cx+dx, cy+dy, color)
		}
	}
}

func drawCornerOrnament(canvas *glow.Canvas, x, y, corner int) {
	size := 15
	col := ZelligeGold

	// Draw L-shaped ornament in each corner
	switch corner {
	case 1: // top-left
		canvas.DrawLine(x, y, x+size, y, col)
		canvas.DrawLine(x, y, x, y+size, col)
		canvas.DrawLine(x+2, y+2, x+size-3, y+2, col)
		canvas.DrawLine(x+2, y+2, x+2, y+size-3, col)
	case 2: // top-right
		canvas.DrawLine(x, y, x-size, y, col)
		canvas.DrawLine(x, y, x, y+size, col)
		canvas.DrawLine(x-2, y+2, x-size+3, y+2, col)
		canvas.DrawLine(x-2, y+2, x-2, y+size-3, col)
	case 3: // bottom-left
		canvas.DrawLine(x, y, x+size, y, col)
		canvas.DrawLine(x, y, x, y-size, col)
		canvas.DrawLine(x+2, y-2, x+size-3, y-2, col)
		canvas.DrawLine(x+2, y-2, x+2, y-size+3, col)
	case 4: // bottom-right
		canvas.DrawLine(x, y, x-size, y, col)
		canvas.DrawLine(x, y, x, y-size, col)
		canvas.DrawLine(x-2, y-2, x-size+3, y-2, col)
		canvas.DrawLine(x-2, y-2, x-2, y-size+3, col)
	}
}
