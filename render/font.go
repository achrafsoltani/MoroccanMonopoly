package render

import "github.com/AchrafSoltani/glow"

// DrawText renders text at (x, y) with the given colour and scale.
// Each character is 8x8 pixels at scale 1.
func DrawText(canvas *glow.Canvas, text string, x, y int, color glow.Color, scale int) {
	cx := x
	for _, ch := range text {
		if ch == '\n' {
			y += 8*scale + 2*scale
			cx = x
			continue
		}
		drawChar(canvas, ch, cx, y, color, scale)
		cx += 8 * scale
	}
}

// DrawTextCentered renders text centred horizontally at y.
func DrawTextCentered(canvas *glow.Canvas, text string, centerX, y int, color glow.Color, scale int) {
	w := len(text) * 8 * scale
	x := centerX - w/2
	DrawText(canvas, text, x, y, color, scale)
}

// DrawTextRight renders text right-aligned ending at x.
func DrawTextRight(canvas *glow.Canvas, text string, rightX, y int, color glow.Color, scale int) {
	w := len(text) * 8 * scale
	DrawText(canvas, text, rightX-w, y, color, scale)
}

// DrawTextWrapped renders text within a maximum width, wrapping at word boundaries.
func DrawTextWrapped(canvas *glow.Canvas, text string, x, y, maxWidth int, color glow.Color, scale int) {
	charW := 8 * scale
	lineH := 8*scale + 2*scale
	maxChars := maxWidth / charW
	if maxChars < 1 {
		maxChars = 1
	}

	line := ""
	for _, word := range splitWords(text) {
		test := line
		if len(test) > 0 {
			test += " "
		}
		test += word

		if len(test) > maxChars && len(line) > 0 {
			DrawText(canvas, line, x, y, color, scale)
			y += lineH
			line = word
		} else {
			line = test
		}
	}
	if len(line) > 0 {
		DrawText(canvas, line, x, y, color, scale)
	}
}

func splitWords(s string) []string {
	var words []string
	word := ""
	for _, ch := range s {
		if ch == ' ' {
			if len(word) > 0 {
				words = append(words, word)
				word = ""
			}
		} else {
			word += string(ch)
		}
	}
	if len(word) > 0 {
		words = append(words, word)
	}
	return words
}

func drawChar(canvas *glow.Canvas, ch rune, x, y int, color glow.Color, scale int) {
	idx := int(ch) - 0x20
	if idx < 0 || idx >= len(FontData) {
		return
	}
	glyph := FontData[idx]
	for row := 0; row < 8; row++ {
		bits := glyph[row]
		for col := 0; col < 8; col++ {
			if bits&(1<<(7-col)) != 0 {
				if scale == 1 {
					canvas.SetPixel(x+col, y+row, color)
				} else {
					canvas.DrawRect(x+col*scale, y+row*scale, scale, scale, color)
				}
			}
		}
	}
}

// TextWidth returns the pixel width of a string at the given scale.
func TextWidth(text string, scale int) int {
	return len(text) * 8 * scale
}
