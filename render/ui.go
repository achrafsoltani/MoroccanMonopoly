package render

import "github.com/AchrafSoltani/glow"

// Button represents a clickable UI button.
type Button struct {
	Label   string
	X, Y    int
	W, H    int
	Enabled bool
	Visible bool
}

// NewButton creates a visible, enabled button.
func NewButton(label string, x, y, w, h int) Button {
	return Button{
		Label:   label,
		X:       x,
		Y:       y,
		W:       w,
		H:       h,
		Enabled: true,
		Visible: true,
	}
}

// Contains returns true if (mx, my) is inside the button.
func (b *Button) Contains(mx, my int) bool {
	return mx >= b.X && mx < b.X+b.W && my >= b.Y && my < b.Y+b.H
}

// Draw renders the button.
func (b *Button) Draw(canvas *glow.Canvas, mouseX, mouseY int) {
	if !b.Visible {
		return
	}

	bg := ButtonBg
	if !b.Enabled {
		bg = glow.Color{R: 50, G: 50, B: 50}
	} else if b.Contains(mouseX, mouseY) {
		bg = ButtonHover
	}

	canvas.DrawRect(b.X, b.Y, b.W, b.H, bg)
	canvas.DrawRectOutline(b.X, b.Y, b.W, b.H, PanelBorder)

	textCol := ButtonText
	if !b.Enabled {
		textCol = glow.Color{R: 100, G: 100, B: 100}
	}

	DrawTextCentered(canvas, b.Label, b.X+b.W/2, b.Y+(b.H-8)/2, textCol, 1)
}

// DrawButtonAt draws a standalone button (for dialogs).
func DrawButtonAt(canvas *glow.Canvas, label string, x, y, w, h int, mouseX, mouseY int, enabled bool) bool {
	bg := ButtonBg
	if !enabled {
		bg = glow.Color{R: 50, G: 50, B: 50}
	} else if mouseX >= x && mouseX < x+w && mouseY >= y && mouseY < y+h {
		bg = ButtonHover
	}

	canvas.DrawRect(x, y, w, h, bg)
	canvas.DrawRectOutline(x, y, w, h, PanelBorder)

	textCol := ButtonText
	if !enabled {
		textCol = glow.Color{R: 100, G: 100, B: 100}
	}
	DrawTextCentered(canvas, label, x+w/2, y+(h-8)/2, textCol, 1)

	// Return true if mouse is hovering and button is enabled
	return enabled && mouseX >= x && mouseX < x+w && mouseY >= y && mouseY < y+h
}
