package render

import (
	"fmt"
	"strings"

	"github.com/AchrafSoltani/MoroccanMonopoly/board"
	"github.com/AchrafSoltani/MoroccanMonopoly/config"
	"github.com/AchrafSoltani/glow"
)

// SpaceRect holds the pixel bounds of a board space for hit-testing and drawing.
type SpaceRect struct {
	X, Y, W, H int
}

// BoardRenderer draws the board and computes space positions.
type BoardRenderer struct {
	Layout     config.Layout
	SpaceRects [config.SpaceCount]SpaceRect
}

// NewBoardRenderer computes all space rectangles using the default layout.
func NewBoardRenderer() *BoardRenderer {
	br := &BoardRenderer{
		Layout: config.NewLayout(config.WindowWidth, config.WindowHeight),
	}
	br.computeRects()
	return br
}

// Recompute stores a new layout and recalculates space rectangles.
func (br *BoardRenderer) Recompute(l config.Layout) {
	br.Layout = l
	br.computeRects()
}

func (br *BoardRenderer) computeRects() {
	bx := br.Layout.BoardX
	by := br.Layout.BoardY
	bs := br.Layout.BoardSize
	cs := br.Layout.CornerSize
	sw := br.Layout.SpaceWidth

	// Bottom row (right to left: space 0 is bottom-right corner)
	// Space 0: bottom-right corner (GO)
	br.SpaceRects[0] = SpaceRect{bx + bs - cs, by + bs - cs, cs, cs}
	// Spaces 1-9: bottom row, right to left
	for i := 1; i <= 9; i++ {
		br.SpaceRects[i] = SpaceRect{
			X: bx + bs - cs - i*sw,
			Y: by + bs - cs,
			W: sw,
			H: cs,
		}
	}
	// Space 10: bottom-left corner (Jail)
	br.SpaceRects[10] = SpaceRect{bx, by + bs - cs, cs, cs}

	// Left column (bottom to top: spaces 11-19)
	for i := 1; i <= 9; i++ {
		br.SpaceRects[10+i] = SpaceRect{
			X: bx,
			Y: by + bs - cs - i*sw,
			W: cs,
			H: sw,
		}
	}
	// Space 20: top-left corner (Free Parking)
	br.SpaceRects[20] = SpaceRect{bx, by, cs, cs}

	// Top row (left to right: spaces 21-29)
	for i := 1; i <= 9; i++ {
		br.SpaceRects[20+i] = SpaceRect{
			X: bx + cs + (i-1)*sw,
			Y: by,
			W: sw,
			H: cs,
		}
	}
	// Space 30: top-right corner (Go To Jail)
	br.SpaceRects[30] = SpaceRect{bx + bs - cs, by, cs, cs}

	// Right column (top to bottom: spaces 31-39)
	for i := 1; i <= 9; i++ {
		br.SpaceRects[30+i] = SpaceRect{
			X: bx + bs - cs,
			Y: by + cs + (i-1)*sw,
			W: cs,
			H: sw,
		}
	}
}

// Draw renders the entire board.
func (br *BoardRenderer) Draw(canvas *glow.Canvas, b *board.Board) {
	l := br.Layout

	// Board background
	canvas.DrawRect(l.BoardX, l.BoardY, l.BoardSize, l.BoardSize, BoardBg)

	// Draw each space
	for i := 0; i < config.SpaceCount; i++ {
		br.drawSpace(canvas, b, i)
	}

	// Board outer border
	drawRectOutline(canvas, l.BoardX, l.BoardY, l.BoardSize, l.BoardSize, BoardBorder, 2)

	// Inner border (around the centre area)
	cs := l.CornerSize
	innerX := l.BoardX + cs
	innerY := l.BoardY + cs
	innerW := l.BoardSize - 2*cs
	innerH := l.BoardSize - 2*cs
	drawRectOutline(canvas, innerX, innerY, innerW, innerH, BoardBorder, 1)

	// Centre area: game title
	centerX := l.BoardX + l.BoardSize/2
	centerY := l.BoardY + l.BoardSize/2
	DrawTextCentered(canvas, "MONOPOLY", centerX, centerY-30, ZelligeGreen, 3)
	DrawTextCentered(canvas, "MAROC", centerX, centerY+5, ZelligeGold, 2)

	// Decorative diamond in centre
	drawDiamond(canvas, centerX, centerY+40, 20, ZelligeGreen)
}

// drawSpace renders a single board space.
func (br *BoardRenderer) drawSpace(canvas *glow.Canvas, b *board.Board, index int) {
	r := br.SpaceRects[index]
	space := b.Spaces[index]

	// Space background
	canvas.DrawRect(r.X, r.Y, r.W, r.H, SpaceBg)

	// Space border
	drawRectOutline(canvas, r.X, r.Y, r.W, r.H, BoardBorder, 1)

	// Colour strip for properties
	br.drawColorStrip(canvas, space, r, index)

	// Space name and price
	br.drawSpaceText(canvas, space, r, index)

	// Houses/hotels
	br.drawHouses(canvas, b, index, r)
}

// drawColorStrip draws the colour band on property spaces.
func (br *BoardRenderer) drawColorStrip(canvas *glow.Canvas, space board.Space, r SpaceRect, index int) {
	if space.Type != board.SpaceProperty {
		return
	}

	col := GroupColor(space.Group)
	stripH := 16

	side := spaceSide(index)
	switch side {
	case 0: // bottom row — strip at top of space
		canvas.DrawRect(r.X+1, r.Y+1, r.W-2, stripH, col)
	case 1: // left column — strip at right of space
		canvas.DrawRect(r.X+r.W-stripH-1, r.Y+1, stripH, r.H-2, col)
	case 2: // top row — strip at bottom of space
		canvas.DrawRect(r.X+1, r.Y+r.H-stripH-1, r.W-2, stripH, col)
	case 3: // right column — strip at left of space
		canvas.DrawRect(r.X+1, r.Y+1, stripH, r.H-2, col)
	}
}

// drawSpaceText draws name and price text on a space.
func (br *BoardRenderer) drawSpaceText(canvas *glow.Canvas, space board.Space, r SpaceRect, index int) {
	side := spaceSide(index)
	isCorner := index%10 == 0

	if isCorner {
		br.drawCornerText(canvas, space, r)
		return
	}

	// Use ShortName if available, otherwise fall back to abbreviation
	name := space.ShortName
	if name == "" {
		name = abbreviate(space.Name, 7)
	}
	price := ""
	if space.Price > 0 {
		price = fmt.Sprintf("%dMAD", space.Price)
	}

	// Special type labels
	switch space.Type {
	case board.SpaceChance:
		name = "CHANCE"
	case board.SpaceCommunityChest:
		name = "CAISSE"
	case board.SpaceTax:
		price = fmt.Sprintf("%dMAD", space.TaxAmount)
	}

	switch side {
	case 0: // bottom row
		stripOffset := 0
		if space.Type == board.SpaceProperty {
			stripOffset = 18
		}
		DrawTextCentered(canvas, name, r.X+r.W/2, r.Y+stripOffset+2, TextDark, 1)
		if price != "" {
			DrawTextCentered(canvas, price, r.X+r.W/2, r.Y+r.H-12, TextDark, 1)
		}
	case 1: // left column
		stripOffset := 0
		if space.Type == board.SpaceProperty {
			stripOffset = 18
		}
		drawTextVertical(canvas, name, r.X+2, r.Y+r.H/2, TextDark, 1)
		if price != "" {
			DrawText(canvas, price, r.X+r.W-stripOffset-TextWidth(price, 1)-2, r.Y+r.H/2-4, TextDark, 1)
		}
	case 2: // top row
		stripOffset := 0
		if space.Type == board.SpaceProperty {
			stripOffset = 18
		}
		DrawTextCentered(canvas, name, r.X+r.W/2, r.Y+2, TextDark, 1)
		if price != "" {
			DrawTextCentered(canvas, price, r.X+r.W/2, r.Y+r.H-stripOffset-12, TextDark, 1)
		}
	case 3: // right column
		stripOffset := 0
		if space.Type == board.SpaceProperty {
			stripOffset = 18
		}
		drawTextVertical(canvas, name, r.X+r.W-10, r.Y+r.H/2, TextDark, 1)
		if price != "" {
			DrawText(canvas, price, r.X+stripOffset+2, r.Y+r.H/2-4, TextDark, 1)
		}
	}
}

// drawCornerText renders text for the four corner spaces.
func (br *BoardRenderer) drawCornerText(canvas *glow.Canvas, space board.Space, r SpaceRect) {
	cx := r.X + r.W/2
	cy := r.Y + r.H/2

	switch space.Type {
	case board.SpaceGo:
		canvas.DrawRect(r.X+1, r.Y+1, r.W-2, r.H-2, ColorGo)
		DrawTextCentered(canvas, "DEPART", cx, cy-12, TextDark, 1)
		DrawTextCentered(canvas, "Collect", cx, cy+2, TextDark, 1)
		DrawTextCentered(canvas, "200 MAD", cx, cy+14, TextDark, 1)
		// Arrow
		canvas.DrawLine(r.X+15, cy, r.X+r.W-15, cy-10, ZelligeGreen)
		canvas.DrawLine(r.X+15, cy, r.X+r.W-15, cy+10, ZelligeGreen)
	case board.SpaceJail:
		canvas.DrawRect(r.X+1, r.Y+1, r.W-2, r.H-2, ColorJail)
		DrawTextCentered(canvas, "EN", cx, cy-12, TextDark, 1)
		DrawTextCentered(canvas, "VISITE", cx, cy, TextDark, 1)
		// Jail bars
		for i := 0; i < 4; i++ {
			bx := r.X + 20 + i*15
			canvas.DrawLine(bx, r.Y+50, bx, r.Y+80, TextDark)
		}
	case board.SpaceFreeParking:
		canvas.DrawRect(r.X+1, r.Y+1, r.W-2, r.H-2, ColorParking)
		DrawTextCentered(canvas, "PARKING", cx, cy-12, TextDark, 1)
		DrawTextCentered(canvas, "GRATUIT", cx, cy+2, TextDark, 1)
	case board.SpaceGoToJail:
		canvas.DrawRect(r.X+1, r.Y+1, r.W-2, r.H-2, ColorGoToJail)
		DrawTextCentered(canvas, "ALLEZ", cx, cy-12, TextDark, 1)
		DrawTextCentered(canvas, "EN", cx, cy, TextDark, 1)
		DrawTextCentered(canvas, "PRISON", cx, cy+12, TextDark, 1)
	}
}

// drawHouses draws house/hotel indicators on a property space.
func (br *BoardRenderer) drawHouses(canvas *glow.Canvas, b *board.Board, index int, r SpaceRect) {
	prop := b.Properties[index]
	if prop.Houses == 0 || prop.Mortgaged {
		return
	}

	side := spaceSide(index)

	if prop.Houses == 5 { // hotel
		br.drawHotelIndicator(canvas, r, side)
	} else {
		for h := 0; h < prop.Houses; h++ {
			br.drawHouseIndicator(canvas, r, side, h, prop.Houses)
		}
	}
}

func (br *BoardRenderer) drawHouseIndicator(canvas *glow.Canvas, r SpaceRect, side, index, total int) {
	size := 6
	gap := 2
	totalW := total*size + (total-1)*gap

	switch side {
	case 0: // bottom row — houses on the colour strip
		startX := r.X + (r.W-totalW)/2 + index*(size+gap)
		canvas.DrawRect(startX, r.Y+3, size, size, HouseColor)
	case 1: // left column
		startY := r.Y + (r.H-totalW)/2 + index*(size+gap)
		canvas.DrawRect(r.X+r.W-20, startY, size, size, HouseColor)
	case 2: // top row
		startX := r.X + (r.W-totalW)/2 + index*(size+gap)
		canvas.DrawRect(startX, r.Y+r.H-20+3, size, size, HouseColor)
	case 3: // right column
		startY := r.Y + (r.H-totalW)/2 + index*(size+gap)
		canvas.DrawRect(r.X+12, startY, size, size, HouseColor)
	}
}

func (br *BoardRenderer) drawHotelIndicator(canvas *glow.Canvas, r SpaceRect, side int) {
	size := 10
	switch side {
	case 0:
		canvas.DrawRect(r.X+(r.W-size)/2, r.Y+3, size, size, HotelColor)
	case 1:
		canvas.DrawRect(r.X+r.W-20, r.Y+(r.H-size)/2, size, size, HotelColor)
	case 2:
		canvas.DrawRect(r.X+(r.W-size)/2, r.Y+r.H-20+3, size, size, HotelColor)
	case 3:
		canvas.DrawRect(r.X+12, r.Y+(r.H-size)/2, size, size, HotelColor)
	}
}

// DrawOwnershipDots draws small dots on properties to show ownership.
func (br *BoardRenderer) DrawOwnershipDots(canvas *glow.Canvas, b *board.Board) {
	for i := 0; i < config.SpaceCount; i++ {
		prop := b.Properties[i]
		if prop.OwnerID < 0 {
			continue
		}
		r := br.SpaceRects[i]
		col := PlayerColors[prop.OwnerID%4]
		if prop.Mortgaged {
			col = MortgageColor
		}
		// Small dot in centre-bottom of space
		cx := r.X + r.W/2
		cy := r.Y + r.H - 8
		canvas.FillCircle(cx, cy, 3, col)
	}
}

// Helper functions

// spaceSide returns which side of the board a space is on:
// 0=bottom, 1=left, 2=top, 3=right
func spaceSide(index int) int {
	switch {
	case index <= 10:
		return 0
	case index <= 20:
		return 1
	case index <= 30:
		return 2
	default:
		return 3
	}
}

// GroupColor returns the display colour for a colour group.
func GroupColor(g board.ColorGroup) glow.Color {
	switch g {
	case board.GroupBrown:
		return ColorBrown
	case board.GroupLightBlue:
		return ColorLightBlue
	case board.GroupPink:
		return ColorPink
	case board.GroupOrange:
		return ColorOrange
	case board.GroupRed:
		return ColorRed
	case board.GroupYellow:
		return ColorYellow
	case board.GroupGreen:
		return ColorGreen
	case board.GroupDarkBlue:
		return ColorDarkBlue
	default:
		return SpaceBg
	}
}

// abbreviate shortens a name to fit in a small space.
func abbreviate(name string, maxLen int) string {
	if len(name) <= maxLen {
		return name
	}
	// Try to find a meaningful abbreviation
	parts := strings.Fields(name)
	if len(parts) > 1 {
		// Use first word if short enough
		if len(parts[0]) <= maxLen {
			return parts[0]
		}
	}
	return name[:maxLen]
}

// drawTextVertical draws text vertically (one char per line).
func drawTextVertical(canvas *glow.Canvas, text string, x, centerY int, color glow.Color, scale int) {
	charH := 8 * scale
	totalH := len(text) * (charH + scale)
	startY := centerY - totalH/2
	for i, ch := range text {
		drawChar(canvas, ch, x, startY+i*(charH+scale), color, scale)
	}
}

// drawRectOutline draws a rectangle outline with thickness.
func drawRectOutline(canvas *glow.Canvas, x, y, w, h int, color glow.Color, thickness int) {
	for t := 0; t < thickness; t++ {
		canvas.DrawRectOutline(x+t, y+t, w-2*t, h-2*t, color)
	}
}

// drawDiamond draws a filled diamond shape.
func drawDiamond(canvas *glow.Canvas, cx, cy, size int, color glow.Color) {
	for dy := -size; dy <= size; dy++ {
		w := size - abs(dy)
		for dx := -w; dx <= w; dx++ {
			canvas.SetPixel(cx+dx, cy+dy, color)
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
