package render

import (
	"fmt"

	"github.com/AchrafSoltani/glow"
)

// DialogData holds everything needed to render a modal dialog.
type DialogData struct {
	Title   string
	Lines   []string
	Buttons []DialogButton
}

// DialogButton is a button within a dialog.
type DialogButton struct {
	Label   string
	ID      int
	Enabled bool
}

// DrawDialog renders a centred modal dialog.
func DrawDialog(canvas *glow.Canvas, data DialogData, mouseX, mouseY int) int {
	// Semi-transparent overlay
	for y := 0; y < canvas.Height(); y += 2 {
		for x := 0; x < canvas.Width(); x += 2 {
			canvas.SetPixel(x, y, glow.Color{R: 0, G: 0, B: 0})
		}
	}

	w := 380
	h := 80 + len(data.Lines)*14 + len(data.Buttons)*34
	x := (canvas.Width() - w) / 2
	y := (canvas.Height() - h) / 2

	// Dialog background
	canvas.DrawRect(x, y, w, h, DialogBg)
	drawThickRectOutline(canvas, x, y, w, h, DialogBorder, 2)

	// Title
	DrawTextCentered(canvas, data.Title, x+w/2, y+10, TextGold, 2)

	// Separator
	canvas.DrawLine(x+10, y+32, x+w-10, y+32, DialogBorder)

	// Text lines
	ty := y + 42
	for _, line := range data.Lines {
		DrawTextWrapped(canvas, line, x+15, ty, w-30, TextLight, 1)
		ty += 14
	}

	// Buttons
	ty += 10
	hoveredID := -1
	bw := w - 40
	bh := 26
	for _, btn := range data.Buttons {
		bx := x + 20
		hovered := DrawButtonAt(canvas, btn.Label, bx, ty, bw, bh, mouseX, mouseY, btn.Enabled)
		if hovered {
			hoveredID = btn.ID
		}
		ty += bh + 6
	}

	return hoveredID
}

// DrawPropertyCard renders a property info card on the HUD.
func DrawPropertyCard(canvas *glow.Canvas, x, y, w int, name string, price int, rent [6]int, houseCost int, groupColor glow.Color, ownerName string, houses int, mortgaged bool) {
	h := 160

	// Card background
	canvas.DrawRect(x, y, w, h, glow.Color{R: 250, G: 245, B: 235})
	canvas.DrawRectOutline(x, y, w, h, TextDark)

	// Colour strip
	canvas.DrawRect(x+1, y+1, w-2, 22, groupColor)

	// Name
	DrawTextCentered(canvas, name, x+w/2, y+5, TextDark, 1)

	// Price
	DrawText(canvas, fmt.Sprintf("Price: %d MAD", price), x+8, y+28, TextDark, 1)

	// Rent table
	labels := []string{"Base", "1 House", "2 Houses", "3 Houses", "4 Houses", "Hotel"}
	for i, label := range labels {
		ry := y + 42 + i*14
		DrawText(canvas, label, x+8, ry, TextDark, 1)
		DrawTextRight(canvas, fmt.Sprintf("%d", rent[i]), x+w-8, ry, TextDark, 1)
	}

	// House cost
	DrawText(canvas, fmt.Sprintf("House: %d MAD", houseCost), x+8, y+128, TextDark, 1)

	// Owner
	if ownerName != "" {
		DrawText(canvas, "Owner: "+ownerName, x+8, y+142, TextDark, 1)
	}
	if mortgaged {
		DrawTextCentered(canvas, "MORTGAGED", x+w/2, y+h-14, ColorRed, 1)
	}
}

func drawThickRectOutline(canvas *glow.Canvas, x, y, w, h int, color glow.Color, thickness int) {
	for t := 0; t < thickness; t++ {
		canvas.DrawRectOutline(x+t, y+t, w-2*t, h-2*t, color)
	}
}
