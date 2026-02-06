package render

import (
	"fmt"

	"github.com/AchrafSoltani/glow"
)

// PlayerInfo holds the data the HUD needs about a player.
type PlayerInfo struct {
	ID       int
	Name     string
	Money    int
	Bankrupt bool
	InJail   bool
	IsAI     bool
}

// HUDData holds all the data the HUD needs to render.
type HUDData struct {
	CurrentPlayerID int
	Players         []PlayerInfo
	Messages        []string
	Die1, Die2      int
	Phase           string
}

// DrawHUD renders the right-side info panel.
func DrawHUD(canvas *glow.Canvas, data HUDData, panelX, panelWidth, winHeight int) {
	px := panelX
	pw := panelWidth

	// Panel background
	canvas.DrawRect(px, 0, pw, winHeight, PanelBg)
	canvas.DrawLine(px, 0, px, winHeight, PanelBorder)

	// Title bar
	DrawTextCentered(canvas, "MONOPOLY MAROC", px+pw/2, 10, TextGold, 2)
	canvas.DrawLine(px+10, 35, px+pw-10, 35, PanelBorder)

	// Current player info
	y := 45

	DrawText(canvas, "Current Turn:", px+15, y, TextLight, 1)
	y += 14

	for _, p := range data.Players {
		if p.ID == data.CurrentPlayerID {
			col := PlayerColors[p.ID%4]
			DrawTokenAt(canvas, p.ID, px+25, y+6, 5)
			tag := ""
			if p.IsAI {
				tag = " (AI)"
			}
			DrawText(canvas, fmt.Sprintf("%s%s", p.Name, tag), px+40, y, col, 1)
			y += 12
			DrawText(canvas, fmt.Sprintf("Money: %d MAD", p.Money), px+40, y, TextLight, 1)
			if p.InJail {
				y += 12
				DrawText(canvas, "** IN JAIL **", px+40, y, ColorRed, 1)
			}
			break
		}
	}

	y += 20
	canvas.DrawLine(px+10, y, px+pw-10, y, PanelBorder)
	y += 8

	// Dice display
	if data.Die1 > 0 {
		DrawText(canvas, fmt.Sprintf("Dice: %d + %d = %d", data.Die1, data.Die2, data.Die1+data.Die2), px+15, y, TextLight, 1)
		if data.Die1 == data.Die2 {
			w := TextWidth(fmt.Sprintf("Dice: %d + %d = %d", data.Die1, data.Die2, data.Die1+data.Die2), 1)
			DrawText(canvas, " DOUBLES!", px+15+w, y, TextGold, 1)
		}
	}
	y += 14

	// Phase indicator
	if data.Phase != "" {
		DrawText(canvas, data.Phase, px+15, y, glow.Color{R: 150, G: 180, B: 150}, 1)
	}
	y += 16

	canvas.DrawLine(px+10, y, px+pw-10, y, PanelBorder)
	y += 8

	// All players summary
	DrawText(canvas, "Players:", px+15, y, TextLight, 1)
	y += 14

	for _, p := range data.Players {
		col := PlayerColors[p.ID%4]
		status := fmt.Sprintf("%d MAD", p.Money)
		if p.Bankrupt {
			status = "BANKRUPT"
			col = MortgageColor
		} else if p.InJail {
			status += " [JAIL]"
		}
		tag := ""
		if p.IsAI {
			tag = "(AI) "
		}

		indicator := "  "
		if p.ID == data.CurrentPlayerID {
			indicator = "> "
		}

		DrawTokenAt(canvas, p.ID, px+25, y+4, 4)
		text := fmt.Sprintf("%s%s%s: %s", indicator, tag, p.Name, status)
		DrawText(canvas, text, px+38, y, col, 1)
		y += 14
	}

	y += 10
	canvas.DrawLine(px+10, y, px+pw-10, y, PanelBorder)
	y += 8

	// Message log
	DrawText(canvas, "Log:", px+15, y, TextLight, 1)
	y += 14

	maxShow := 16
	start := 0
	if len(data.Messages) > maxShow {
		start = len(data.Messages) - maxShow
	}
	for i := start; i < len(data.Messages); i++ {
		DrawText(canvas, data.Messages[i], px+15, y, glow.Color{R: 180, G: 200, B: 180}, 1)
		y += 10
	}
}
