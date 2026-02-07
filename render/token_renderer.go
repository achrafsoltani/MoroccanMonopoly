package render

import (
	"math"

	"github.com/AchrafSoltani/MoroccanMonopoly/player"
	"github.com/AchrafSoltani/glow"
)

// TokenShape identifies the geometric shape for each player.
type TokenShape int

const (
	TokenCircle  TokenShape = iota
	TokenDiamond
	TokenSquare
	TokenTriangle
)

// DrawToken draws a player's token on the board at the given space rect.
func DrawToken(canvas *glow.Canvas, p *player.Player, r SpaceRect) {
	col := PlayerColors[p.ID%4]
	shape := TokenShape(p.ID % 4)

	var cx, cy int
	if p.InJail {
		// Position jailed tokens in the lower-left "In Jail" subsection of the jail corner
		cx = r.X + 12 + (p.ID%2)*18
		cy = r.Y + r.H - 18 + (p.ID/2)*10
	} else {
		// Offset tokens so multiple players on the same space don't overlap
		ox := 15 + (p.ID%2)*25
		oy := 30 + (p.ID/2)*25
		cx = r.X + ox
		cy = r.Y + oy
	}
	size := 8

	switch shape {
	case TokenCircle:
		canvas.FillCircle(cx, cy, size, col)
		canvas.DrawCircle(cx, cy, size, TextDark)
	case TokenDiamond:
		drawFilledDiamond(canvas, cx, cy, size, col)
		drawDiamondOutline(canvas, cx, cy, size, TextDark)
	case TokenSquare:
		canvas.DrawRect(cx-size, cy-size, size*2, size*2, col)
		canvas.DrawRectOutline(cx-size, cy-size, size*2, size*2, TextDark)
	case TokenTriangle:
		drawFilledTriangle(canvas, cx, cy-size, cx-size, cy+size, cx+size, cy+size, col)
		canvas.DrawTriangle(cx, cy-size, cx-size, cy+size, cx+size, cy+size, TextDark)
	}

	// Draw jail bars over jailed tokens
	if p.InJail {
		barCol := glow.Color{R: 100, G: 100, B: 100}
		for i := 0; i < 3; i++ {
			bx := cx - 8 + i*8
			canvas.DrawLine(bx, cy-10, bx, cy+10, barCol)
		}
	}
}

// DrawTokenHighlight draws a pulsing ring around the current player's token.
func DrawTokenHighlight(canvas *glow.Canvas, p *player.Player, r SpaceRect, timer float64) {
	col := PlayerColors[p.ID%4]

	var cx, cy int
	if p.InJail {
		cx = r.X + 12 + (p.ID%2)*18
		cy = r.Y + r.H - 18 + (p.ID/2)*10
	} else {
		ox := 15 + (p.ID%2)*25
		oy := 30 + (p.ID/2)*25
		cx = r.X + ox
		cy = r.Y + oy
	}

	// Pulsing radius
	pulse := math.Sin(timer * 4.0)
	radius := 12 + int(pulse*3)

	// Draw concentric ring using the player's colour
	canvas.DrawCircle(cx, cy, radius, col)
	canvas.DrawCircle(cx, cy, radius+1, col)
}

// DrawTokenAt draws a token at an arbitrary position (for HUD/dialogs).
func DrawTokenAt(canvas *glow.Canvas, playerID int, cx, cy, size int) {
	col := PlayerColors[playerID%4]
	shape := TokenShape(playerID % 4)

	switch shape {
	case TokenCircle:
		canvas.FillCircle(cx, cy, size, col)
		canvas.DrawCircle(cx, cy, size, TextDark)
	case TokenDiamond:
		drawFilledDiamond(canvas, cx, cy, size, col)
	case TokenSquare:
		canvas.DrawRect(cx-size, cy-size, size*2, size*2, col)
	case TokenTriangle:
		drawFilledTriangle(canvas, cx, cy-size, cx-size, cy+size, cx+size, cy+size, col)
	}
}

func drawFilledDiamond(canvas *glow.Canvas, cx, cy, size int, color glow.Color) {
	for dy := -size; dy <= size; dy++ {
		w := size - abs(dy)
		for dx := -w; dx <= w; dx++ {
			canvas.SetPixel(cx+dx, cy+dy, color)
		}
	}
}

func drawDiamondOutline(canvas *glow.Canvas, cx, cy, size int, color glow.Color) {
	canvas.DrawLine(cx, cy-size, cx+size, cy, color)
	canvas.DrawLine(cx+size, cy, cx, cy+size, color)
	canvas.DrawLine(cx, cy+size, cx-size, cy, color)
	canvas.DrawLine(cx-size, cy, cx, cy-size, color)
}

func drawFilledTriangle(canvas *glow.Canvas, x0, y0, x1, y1, x2, y2 int, color glow.Color) {
	// Simple scanline fill for triangle
	minY := min3(y0, y1, y2)
	maxY := max3(y0, y1, y2)

	for y := minY; y <= maxY; y++ {
		minX := 10000
		maxX := -10000

		// Check each edge
		edges := [][4]int{{x0, y0, x1, y1}, {x1, y1, x2, y2}, {x2, y2, x0, y0}}
		for _, e := range edges {
			ex0, ey0, ex1, ey1 := e[0], e[1], e[2], e[3]
			if (ey0 <= y && ey1 > y) || (ey1 <= y && ey0 > y) {
				x := ex0 + (y-ey0)*(ex1-ex0)/(ey1-ey0)
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
			}
		}

		if minX <= maxX {
			for x := minX; x <= maxX; x++ {
				canvas.SetPixel(x, y, color)
			}
		}
	}
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func max3(a, b, c int) int {
	if a > b {
		if a > c {
			return a
		}
		return c
	}
	if b > c {
		return b
	}
	return c
}
