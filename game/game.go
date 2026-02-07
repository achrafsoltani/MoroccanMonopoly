package game

import (
	"fmt"
	"sort"

	"github.com/AchrafSoltani/MoroccanMonopoly/audio"
	"github.com/AchrafSoltani/MoroccanMonopoly/board"
	"github.com/AchrafSoltani/MoroccanMonopoly/config"
	"github.com/AchrafSoltani/MoroccanMonopoly/player"
	"github.com/AchrafSoltani/MoroccanMonopoly/render"
	"github.com/AchrafSoltani/MoroccanMonopoly/save"
	"github.com/AchrafSoltani/glow"
)

// Game holds all game state and coordinates updates and rendering.
type Game struct {
	State     GameState
	Phase     TurnPhase
	Dialog    DialogType
	Board     *board.Board
	Players   []*player.Player
	Current   int // index of current player
	GameTimer float64
	Layout    config.Layout

	// Dice
	Die1, Die2    int
	Doubles       bool
	DoublesCount  int
	DiceAnimTimer float64
	DiceRolling   bool

	// Token movement animation
	MoveFrom    int
	MoveTo      int
	MoveTimer   float64
	MoveSteps   int
	MoveCurrent int

	// AI pacing
	AITimer float64

	// Message log
	Messages    []string
	MaxMessages int

	// Renderers
	BoardRenderer *render.BoardRenderer
	Audio         *audio.Engine

	// Mouse state
	MouseX, MouseY int
	MouseClicked   bool
	DialogHovered  int // button ID hovered in dialog (-1 = none)

	// Build/mortgage selection
	SelectableSpaces []int // space indices player can choose from
	SelectedSpace    int   // currently selected space index (-1 = none)

	// Auction state
	AuctionSpaceIdx   int
	AuctionBids       [config.MaxPlayers]int
	AuctionActive     [config.MaxPlayers]bool
	AuctionCurrent    int
	AuctionHighBid    int
	AuctionHighBidder int

	// Trade state
	TradePartner      int // target player index
	TradeOfferedProps  []int
	TradeWantedProps   []int
	TradeOfferedMoney  int
	TradeWantedMoney   int
	TradeStage         TradeState
	TradeOfferJailCard bool // offering a jail card
	TradeWantJailCard  bool // requesting a jail card
	PendingOffer       *TradeOffer // for human-to-human confirmation

	// Buttons
	Buttons []render.Button
}

// NewGame creates a new game with default state.
func NewGame() *Game {
	g := &Game{
		State:         StateMenu,
		Board:         board.NewBoard(),
		BoardRenderer: render.NewBoardRenderer(),
		Audio:         audio.NewEngine(),
		MaxMessages:   12,
		Layout:        config.NewLayout(config.WindowWidth, config.WindowHeight),
	}
	return g
}

// OnResize recomputes layout and repositions all UI elements.
func (g *Game) OnResize(width, height int) {
	g.Layout = config.NewLayout(width, height)
	g.BoardRenderer.Recompute(g.Layout)
	if len(g.Buttons) > 0 {
		g.repositionButtons()
	}
}

// repositionButtons updates button positions from the current layout.
func (g *Game) repositionButtons() {
	margin := 20
	bx := g.Layout.PanelX + margin
	availW := g.Layout.PanelWidth - 2*margin
	gap := 6
	bw := (availW - gap) / 2
	bh := 28

	// Place buttons starting from roughly the vertical middle of the window,
	// clamped to a sensible range.
	by := g.Layout.WinH/2 - 30
	if by < 260 {
		by = 260
	}
	if by > 400 {
		by = 400
	}

	positions := [][4]int{
		{bx, by, bw*2 + gap, bh},                 // Roll Dice (full width)
		{bx, by + bh + gap, bw, bh},               // Buy
		{bx + bw + gap, by + bh + gap, bw, bh},    // Auction
		{bx, by + 2*(bh+gap), bw, bh},             // Build
		{bx + bw + gap, by + 2*(bh+gap), bw, bh},  // Mortgage
		{bx, by + 3*(bh+gap), bw, bh},             // Trade
		{bx + bw + gap, by + 3*(bh+gap), bw, bh},  // End Turn
	}

	for i := range g.Buttons {
		if i < len(positions) {
			g.Buttons[i].X = positions[i][0]
			g.Buttons[i].Y = positions[i][1]
			g.Buttons[i].W = positions[i][2]
			g.Buttons[i].H = positions[i][3]
		}
	}
}

// Update advances the game state by dt seconds.
func (g *Game) Update(dt float64) {
	g.GameTimer += dt

	switch g.State {
	case StateMenu:
		g.updateMenu(dt)
	case StateSetup:
		g.updateSetup(dt)
	case StatePlaying:
		g.updatePlaying(dt)
	case StateGameOver:
		g.updateGameOver(dt)
	}

	g.MouseClicked = false
}

// Draw renders the entire game.
func (g *Game) Draw(canvas *glow.Canvas) {
	switch g.State {
	case StateMenu:
		g.drawMenu(canvas)
	case StateSetup:
		g.drawSetup(canvas)
	case StatePlaying:
		g.drawPlaying(canvas)
	case StateGameOver:
		g.drawGameOver(canvas)
	}
}

// MouseDown handles mouse button press.
func (g *Game) MouseDown(x, y int, button glow.MouseButton) {
	if button == glow.MouseLeft {
		g.MouseX = x
		g.MouseY = y
		g.MouseClicked = true
	}
}

// MouseMove handles mouse movement.
func (g *Game) MouseMove(x, y int) {
	g.MouseX = x
	g.MouseY = y
}

// KeyDown handles key press.
func (g *Game) KeyDown(key glow.Key) {
	switch g.State {
	case StateMenu:
		g.keyMenu(key)
	case StateSetup:
		g.keySetup(key)
	case StatePlaying:
		g.keyPlaying(key)
	case StateGameOver:
		if key == glow.KeyEnter {
			save.DeleteSave()
			g.State = StateMenu
		}
	}
}

// AddMessage appends a message to the log.
func (g *Game) AddMessage(msg string) {
	g.Messages = append(g.Messages, msg)
	if len(g.Messages) > g.MaxMessages {
		g.Messages = g.Messages[len(g.Messages)-g.MaxMessages:]
	}
}

// StartGame initialises a new game with the given players.
func (g *Game) StartGame(players []*player.Player) {
	g.Players = players
	g.Current = 0
	g.State = StatePlaying
	g.Phase = PhasePreRoll
	g.Board = board.NewBoard()
	g.Messages = nil
	g.AddMessage("Game started! Roll the dice.")

	// Set up buttons
	g.setupButtons()
}

func (g *Game) setupButtons() {
	// Create buttons with placeholder positions; repositionButtons() will set the real coords.
	g.Buttons = []render.Button{
		render.NewButton("Roll Dice", 0, 0, 0, 0),
		render.NewButton("Buy", 0, 0, 0, 0),
		render.NewButton("Auction", 0, 0, 0, 0),
		render.NewButton("Build", 0, 0, 0, 0),
		render.NewButton("Mortgage", 0, 0, 0, 0),
		render.NewButton("Trade", 0, 0, 0, 0),
		render.NewButton("End Turn", 0, 0, 0, 0),
	}
	g.repositionButtons()
}

// currentPlayer returns the active player.
func (g *Game) currentPlayer() *player.Player {
	if len(g.Players) == 0 {
		return nil
	}
	return g.Players[g.Current]
}

// alivePlayers returns all non-bankrupt players.
func (g *Game) alivePlayers() []*player.Player {
	var alive []*player.Player
	for _, p := range g.Players {
		if !p.Bankrupt {
			alive = append(alive, p)
		}
	}
	return alive
}

// nextPlayer advances to the next non-bankrupt player.
func (g *Game) nextPlayer() {
	for {
		g.Current = (g.Current + 1) % len(g.Players)
		if !g.Players[g.Current].Bankrupt {
			break
		}
	}
	g.DoublesCount = 0
}

// Stub methods for states not yet implemented

func (g *Game) updateMenu(dt float64)    {}
func (g *Game) updateSetup(dt float64)   {}
func (g *Game) updateGameOver(dt float64) {
	if g.MouseClicked {
		// Return to menu on any click
	}
}

func (g *Game) drawMenu(canvas *glow.Canvas) {
	// Animated background with Moroccan patterns
	render.DrawMenuBackground(canvas, g.GameTimer)

	cx := canvas.Width() / 2

	// Title with shadow
	render.DrawTextCentered(canvas, "MONOPOLY MAROC", cx+2, 172, glow.Color{R: 0, G: 0, B: 0}, 4)
	render.DrawTextCentered(canvas, "MONOPOLY MAROC", cx, 170, render.ZelligeGreen, 4)

	render.DrawTextCentered(canvas, "~ Moroccan Edition ~", cx, 220, render.ZelligeGold, 1)

	// Decorative line
	canvas.DrawLine(cx-120, 240, cx+120, 240, render.ZelligeGold)

	y := 280
	render.DrawTextCentered(canvas, "ENTER - New Game (1 Human + 1 AI)", cx, y, render.TextLight, 1)
	y += 20
	render.DrawTextCentered(canvas, "2 - Two Players (1H + 1AI)", cx, y, render.TextLight, 1)
	y += 16
	render.DrawTextCentered(canvas, "3 - Three Players (1H + 2AI)", cx, y, render.TextLight, 1)
	y += 16
	render.DrawTextCentered(canvas, "4 - Four Players (1H + 3AI)", cx, y, render.TextLight, 1)

	if save.HasSave() {
		y += 30
		render.DrawTextCentered(canvas, "R - Resume Saved Game", cx, y, render.TextGold, 2)
	}

	y += 50
	render.DrawTextCentered(canvas, "F5 = Save during game", cx, y, glow.Color{R: 120, G: 160, B: 120}, 1)

	// Currency display
	y += 40
	render.DrawTextCentered(canvas, "Currency: MAD (Moroccan Dirham)", cx, y, render.ZelligeGold, 1)
}

func (g *Game) drawSetup(canvas *glow.Canvas) {}

func (g *Game) drawPlaying(canvas *glow.Canvas) {
	// Draw board
	g.BoardRenderer.Draw(canvas, g.Board)
	g.BoardRenderer.DrawOwnershipDots(canvas, g.Board)

	// Draw tokens
	for _, p := range g.Players {
		if !p.Bankrupt {
			render.DrawToken(canvas, p, g.BoardRenderer.SpaceRects[p.Position])
		}
	}

	// Highlight current player's token
	cp := g.currentPlayer()
	if cp != nil && !cp.Bankrupt {
		render.DrawTokenHighlight(canvas, cp, g.BoardRenderer.SpaceRects[cp.Position], g.GameTimer)
	}

	// Draw HUD panel
	render.DrawHUD(canvas, g.hudData(), g.Layout.PanelX, g.Layout.PanelWidth, g.Layout.WinH)

	// Draw buttons
	for i := range g.Buttons {
		g.Buttons[i].Draw(canvas, g.MouseX, g.MouseY)
	}

	// Draw dice if rolled
	if g.Die1 > 0 {
		render.DrawDice(canvas, g.Die1, g.Die2,
			g.Layout.BoardX+g.Layout.BoardSize/2-50,
			g.Layout.BoardY+g.Layout.BoardSize/2+80,
			g.DiceRolling, g.DiceAnimTimer)
	}

	// Board space hover highlighting
	g.drawBoardHover(canvas)

	// Draw dialogs
	g.drawDialogs(canvas)
}

func (g *Game) drawBoardHover(canvas *glow.Canvas) {
	if g.Dialog != DialogNone {
		return
	}
	// Find which space the mouse is hovering over
	for i := 0; i < config.SpaceCount; i++ {
		r := g.BoardRenderer.SpaceRects[i]
		if g.MouseX >= r.X && g.MouseX < r.X+r.W && g.MouseY >= r.Y && g.MouseY < r.Y+r.H {
			// Highlight the space
			for dy := 0; dy < r.H; dy += 2 {
				for dx := 0; dx < r.W; dx += 2 {
					canvas.SetPixel(r.X+dx, r.Y+dy, glow.Color{R: 255, G: 255, B: 200})
				}
			}

			// Show property card if it's a property
			space := g.Board.Spaces[i]
			if g.Board.IsProperty(i) {
				prop := g.Board.Properties[i]
				ownerName := ""
				if prop.OwnerID >= 0 {
					ownerName = g.Players[prop.OwnerID].Name
				}
				cardX := g.Layout.PanelX + 20
				cardW := g.Layout.PanelWidth - 40

				switch space.Type {
				case board.SpaceRailroad:
					ownedCount := 0
					if prop.OwnerID >= 0 {
						ownedCount = g.countOwnedRailroads(prop.OwnerID)
					}
					render.DrawRailroadCard(canvas, cardX, g.Layout.WinH-150, cardW,
						space.Name, space.Price, ownerName, ownedCount, prop.Mortgaged)
				case board.SpaceUtility:
					ownedCount := 0
					if prop.OwnerID >= 0 {
						ownedCount = g.countOwnedUtilities(prop.OwnerID)
					}
					render.DrawUtilityCard(canvas, cardX, g.Layout.WinH-130, cardW,
						space.Name, space.Price, ownerName, ownedCount, prop.Mortgaged)
				default:
					groupCol := render.GroupColor(space.Group)
					render.DrawPropertyCard(canvas, cardX, g.Layout.WinH-180, cardW,
						space.Name, space.Price, space.Rent, space.HouseCost,
						groupCol, ownerName, prop.Houses, prop.Mortgaged)
				}
			}
			break
		}
	}
}

func (g *Game) drawDialogs(canvas *glow.Canvas) {
	switch g.Dialog {
	case DialogBuyProperty:
		p := g.currentPlayer()
		space := g.Board.Spaces[p.Position]
		data := render.DialogData{
			Title: "Buy Property?",
			Lines: []string{
				space.Name,
				fmt.Sprintf("Price: %d MAD", space.Price),
				fmt.Sprintf("Your money: %d MAD", p.Money),
			},
			Buttons: []render.DialogButton{
				{Label: fmt.Sprintf("Buy for %d MAD", space.Price), ID: 0, Enabled: p.Money >= space.Price},
				{Label: "Decline (Auction)", ID: 1, Enabled: true},
			},
		}
		g.DialogHovered = render.DrawDialog(canvas, data, g.MouseX, g.MouseY)

	case DialogIncomeTax:
		p := g.currentPlayer()
		tenPercent := g.PlayerNetWorth(p.ID) / 10
		data := render.DialogData{
			Title: "Impot sur le Revenu",
			Lines: []string{
				fmt.Sprintf("%s must pay income tax.", p.Name),
				fmt.Sprintf("Net worth: %d MAD", g.PlayerNetWorth(p.ID)),
			},
			Buttons: []render.DialogButton{
				{Label: "Pay 200 MAD (flat)", ID: 0, Enabled: true},
				{Label: fmt.Sprintf("Pay 10%% (%d MAD)", tenPercent), ID: 1, Enabled: true},
			},
		}
		g.DialogHovered = render.DrawDialog(canvas, data, g.MouseX, g.MouseY)

	case DialogJailOptions:
		p := g.currentPlayer()
		data := render.DialogData{
			Title: "In Jail!",
			Lines: []string{
				fmt.Sprintf("%s is in jail (turn %d/%d)", p.Name, p.JailTurns+1, config.MaxJailTurns),
			},
			Buttons: []render.DialogButton{
				{Label: fmt.Sprintf("Pay %d MAD fine", config.JailFine), ID: 0, Enabled: p.Money >= config.JailFine},
				{Label: "Use Get Out of Jail card", ID: 1, Enabled: p.GetOutOfJailCards > 0},
				{Label: "Try to roll doubles", ID: 2, Enabled: true},
			},
		}
		g.DialogHovered = render.DrawDialog(canvas, data, g.MouseX, g.MouseY)

	case DialogBuild:
		p := g.currentPlayer()
		var btns []render.DialogButton
		for i, idx := range g.SelectableSpaces {
			space := g.Board.Spaces[idx]
			prop := g.Board.Properties[idx]
			level := "house"
			if prop.Houses == config.HousesPerHotel {
				level = "HOTEL"
			}
			label := fmt.Sprintf("%s (%s, %d MAD)", space.Name, level, space.HouseCost)
			btns = append(btns, render.DialogButton{Label: label, ID: i, Enabled: p.Money >= space.HouseCost})
		}
		btns = append(btns, render.DialogButton{Label: "Cancel", ID: -1, Enabled: true})

		data := render.DialogData{
			Title:   "Build",
			Lines:   []string{"Select a property to build on:"},
			Buttons: btns,
		}
		g.DialogHovered = render.DrawDialog(canvas, data, g.MouseX, g.MouseY)

	case DialogMortgage:
		var btns []render.DialogButton
		for i, idx := range g.SelectableSpaces {
			space := g.Board.Spaces[idx]
			prop := g.Board.Properties[idx]
			if prop.Mortgaged {
				cost := g.UnmortgageCost(idx)
				label := fmt.Sprintf("Unmortgage %s (%d MAD)", space.Name, cost)
				btns = append(btns, render.DialogButton{Label: label, ID: i, Enabled: true})
			} else {
				val := g.MortgageValue(idx)
				label := fmt.Sprintf("Mortgage %s (+%d MAD)", space.Name, val)
				btns = append(btns, render.DialogButton{Label: label, ID: i, Enabled: true})
			}
		}
		btns = append(btns, render.DialogButton{Label: "Cancel", ID: -1, Enabled: true})

		data := render.DialogData{
			Title:   "Mortgage",
			Lines:   []string{"Select a property:"},
			Buttons: btns,
		}
		g.DialogHovered = render.DrawDialog(canvas, data, g.MouseX, g.MouseY)

	case DialogAuction:
		p := g.Players[g.AuctionCurrent]
		space := g.Board.Spaces[g.AuctionSpaceIdx]
		data := render.DialogData{
			Title: "Auction",
			Lines: []string{
				fmt.Sprintf("Property: %s", space.Name),
				fmt.Sprintf("Current bid: %d MAD", g.AuctionHighBid),
				fmt.Sprintf("%s's turn to bid", p.Name),
			},
			Buttons: []render.DialogButton{
				{Label: fmt.Sprintf("Bid %d MAD", g.AuctionHighBid+10), ID: 0, Enabled: p.Money > g.AuctionHighBid+10},
				{Label: "Pass", ID: 1, Enabled: true},
			},
		}
		g.DialogHovered = render.DrawDialog(canvas, data, g.MouseX, g.MouseY)

	case DialogTrade:
		p := g.currentPlayer()
		switch g.TradeStage {
		case TradeSelectPartner:
			var btns []render.DialogButton
			for _, other := range g.Players {
				if other.ID != p.ID && !other.Bankrupt {
					label := fmt.Sprintf("%s (%d MAD, %d props)", other.Name, other.Money, len(other.Properties))
					btns = append(btns, render.DialogButton{Label: label, ID: other.ID, Enabled: true})
				}
			}
			btns = append(btns, render.DialogButton{Label: "Cancel", ID: -1, Enabled: true})
			data := render.DialogData{
				Title:   "Trade - Select Partner",
				Lines:   []string{"Who do you want to trade with?"},
				Buttons: btns,
			}
			g.DialogHovered = render.DrawDialog(canvas, data, g.MouseX, g.MouseY)

		case TradeSelectOffer:
			partner := g.Players[g.TradePartner]
			var btns []render.DialogButton

			// Own properties to offer (toggle)
			for _, idx := range p.Properties {
				space := g.Board.Spaces[idx]
				offered := g.tradePropsContains(g.TradeOfferedProps, idx)
				prefix := "[ ] Offer: "
				if offered {
					prefix = "[X] Offer: "
				}
				btns = append(btns, render.DialogButton{
					Label:   prefix + space.Name,
					ID:      idx,
					Enabled: !g.Board.Properties[idx].Mortgaged && g.Board.Properties[idx].Houses == 0,
				})
			}
			// Partner properties to request (toggle)
			for _, idx := range partner.Properties {
				space := g.Board.Spaces[idx]
				wanted := g.tradePropsContains(g.TradeWantedProps, idx)
				prefix := "[ ] Want:  "
				if wanted {
					prefix = "[X] Want:  "
				}
				btns = append(btns, render.DialogButton{
					Label:   prefix + space.Name,
					ID:      1000 + idx,
					Enabled: !g.Board.Properties[idx].Mortgaged && g.Board.Properties[idx].Houses == 0,
				})
			}

			// Money buttons
			btns = append(btns, render.DialogButton{Label: fmt.Sprintf("Offer +50 MAD (now: %d)", g.TradeOfferedMoney), ID: 2000, Enabled: true})
			btns = append(btns, render.DialogButton{Label: fmt.Sprintf("Offer -50 MAD (now: %d)", g.TradeOfferedMoney), ID: 2001, Enabled: g.TradeOfferedMoney >= 50})
			btns = append(btns, render.DialogButton{Label: fmt.Sprintf("Want +50 MAD (now: %d)", g.TradeWantedMoney), ID: 2002, Enabled: true})
			btns = append(btns, render.DialogButton{Label: fmt.Sprintf("Want -50 MAD (now: %d)", g.TradeWantedMoney), ID: 2003, Enabled: g.TradeWantedMoney >= 50})

			// Jail card toggles
			if p.GetOutOfJailCards > 0 {
				jailLabel := "[ ] Offer Jail Card"
				if g.TradeOfferJailCard {
					jailLabel = "[X] Offer Jail Card"
				}
				btns = append(btns, render.DialogButton{Label: jailLabel, ID: 2004, Enabled: true})
			}
			if partner.GetOutOfJailCards > 0 {
				jailLabel := "[ ] Want Jail Card"
				if g.TradeWantJailCard {
					jailLabel = "[X] Want Jail Card"
				}
				btns = append(btns, render.DialogButton{Label: jailLabel, ID: 2005, Enabled: true})
			}

			// Propose / Cancel
			hasContent := len(g.TradeOfferedProps) > 0 || len(g.TradeWantedProps) > 0 ||
				g.TradeOfferedMoney > 0 || g.TradeWantedMoney > 0 ||
				g.TradeOfferJailCard || g.TradeWantJailCard
			btns = append(btns, render.DialogButton{Label: "Propose Trade >>", ID: 3000, Enabled: hasContent})
			btns = append(btns, render.DialogButton{Label: "Cancel", ID: -1, Enabled: true})

			lines := []string{fmt.Sprintf("Building offer with %s:", partner.Name)}
			data := render.DialogData{
				Title:   "Trade - Build Offer",
				Lines:   lines,
				Buttons: btns,
			}
			g.DialogHovered = render.DrawDialog(canvas, data, g.MouseX, g.MouseY)

		case TradeConfirm:
			partner := g.Players[g.TradePartner]
			var lines []string
			lines = append(lines, fmt.Sprintf("Trade with %s:", partner.Name))
			lines = append(lines, "--- You give ---")
			for _, idx := range g.TradeOfferedProps {
				lines = append(lines, "  "+g.Board.Spaces[idx].Name)
			}
			if g.TradeOfferedMoney > 0 {
				lines = append(lines, fmt.Sprintf("  %d MAD", g.TradeOfferedMoney))
			}
			if g.TradeOfferJailCard {
				lines = append(lines, "  Jail Card")
			}
			lines = append(lines, "--- You get ---")
			for _, idx := range g.TradeWantedProps {
				lines = append(lines, "  "+g.Board.Spaces[idx].Name)
			}
			if g.TradeWantedMoney > 0 {
				lines = append(lines, fmt.Sprintf("  %d MAD", g.TradeWantedMoney))
			}
			if g.TradeWantJailCard {
				lines = append(lines, "  Jail Card")
			}

			data := render.DialogData{
				Title: "Confirm Trade",
				Lines: lines,
				Buttons: []render.DialogButton{
					{Label: "Propose Trade", ID: 0, Enabled: true},
					{Label: "Go Back", ID: 1, Enabled: true},
					{Label: "Cancel", ID: -1, Enabled: true},
				},
			}
			g.DialogHovered = render.DrawDialog(canvas, data, g.MouseX, g.MouseY)
		}

	case DialogTradeReceived:
		// Show incoming trade offer for human-to-human trades
		if g.PendingOffer != nil {
			from := g.Players[g.PendingOffer.FromPlayer]
			var lines []string
			lines = append(lines, fmt.Sprintf("%s offers you:", from.Name))
			for _, idx := range g.PendingOffer.OfferedProps {
				lines = append(lines, "  "+g.Board.Spaces[idx].Name)
			}
			if g.PendingOffer.OfferedMoney > 0 {
				lines = append(lines, fmt.Sprintf("  %d MAD", g.PendingOffer.OfferedMoney))
			}
			if g.PendingOffer.OfferedJailCards > 0 {
				lines = append(lines, "  Jail Card")
			}
			lines = append(lines, "In exchange for:")
			for _, idx := range g.PendingOffer.WantedProps {
				lines = append(lines, "  "+g.Board.Spaces[idx].Name)
			}
			if g.PendingOffer.WantedMoney > 0 {
				lines = append(lines, fmt.Sprintf("  %d MAD", g.PendingOffer.WantedMoney))
			}
			if g.PendingOffer.WantedJailCards > 0 {
				lines = append(lines, "  Jail Card")
			}
			data := render.DialogData{
				Title: "Trade Offer Received",
				Lines: lines,
				Buttons: []render.DialogButton{
					{Label: "Accept", ID: 0, Enabled: true},
					{Label: "Decline", ID: 1, Enabled: true},
				},
			}
			g.DialogHovered = render.DrawDialog(canvas, data, g.MouseX, g.MouseY)
		}
	}
}

func (g *Game) drawGameOver(canvas *glow.Canvas) {
	// Animated zellige background
	render.DrawMenuBackground(canvas, g.GameTimer)

	cx := canvas.Width() / 2

	// Title
	render.DrawTextCentered(canvas, "GAME OVER", cx+2, 102, glow.Color{R: 0, G: 0, B: 0}, 4)
	render.DrawTextCentered(canvas, "GAME OVER", cx, 100, render.TextGold, 4)

	// Winner announcement
	alive := g.alivePlayers()
	if len(alive) == 1 {
		winner := alive[0]
		render.DrawTextCentered(canvas, winner.Name+" WINS!", cx, 160, render.PlayerColors[winner.ID%4], 3)
	}

	// Build rankings sorted by net worth (descending)
	type playerStat struct {
		p        *player.Player
		netWorth int
		houses   int
		hotels   int
	}
	var stats []playerStat
	for _, p := range g.Players {
		nw := 0
		houses := 0
		hotels := 0
		if !p.Bankrupt {
			nw = g.PlayerNetWorth(p.ID)
			for _, idx := range p.Properties {
				h := g.Board.Properties[idx].Houses
				if h == config.HotelLevel {
					hotels++
				} else {
					houses += h
				}
			}
		}
		stats = append(stats, playerStat{p: p, netWorth: nw, houses: houses, hotels: hotels})
	}
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].p.Bankrupt != stats[j].p.Bankrupt {
			return !stats[i].p.Bankrupt // non-bankrupt first
		}
		return stats[i].netWorth > stats[j].netWorth
	})

	// Rankings table
	y := 210
	render.DrawTextCentered(canvas, "--- Final Standings ---", cx, y, render.ZelligeGold, 1)
	y += 20

	// Header
	hx := cx - 180
	render.DrawText(canvas, "#", hx, y, render.TextLight, 1)
	render.DrawText(canvas, "Player", hx+20, y, render.TextLight, 1)
	render.DrawText(canvas, "Cash", hx+160, y, render.TextLight, 1)
	render.DrawText(canvas, "Worth", hx+230, y, render.TextLight, 1)
	render.DrawText(canvas, "Props", hx+300, y, render.TextLight, 1)
	render.DrawText(canvas, "H/Ht", hx+350, y, render.TextLight, 1)
	y += 16

	// Separator
	canvas.DrawLine(hx, y-2, hx+380, y-2, render.PanelBorder)

	for rank, s := range stats {
		col := render.PlayerColors[s.p.ID%4]
		status := ""
		if s.p.Bankrupt {
			col = render.MortgageColor
			status = " [BANKRUPT]"
		}
		render.DrawText(canvas, fmt.Sprintf("%d", rank+1), hx, y, col, 1)
		render.DrawText(canvas, s.p.Name+status, hx+20, y, col, 1)
		render.DrawText(canvas, fmt.Sprintf("%d", s.p.Money), hx+160, y, col, 1)
		render.DrawText(canvas, fmt.Sprintf("%d", s.netWorth), hx+230, y, col, 1)
		render.DrawText(canvas, fmt.Sprintf("%d", len(s.p.Properties)), hx+300, y, col, 1)
		render.DrawText(canvas, fmt.Sprintf("%d/%d", s.houses, s.hotels), hx+350, y, col, 1)
		y += 16
	}

	y += 20
	render.DrawTextCentered(canvas, "Press ENTER to return to menu", cx, y, render.TextLight, 1)
}

func (g *Game) keyMenu(key glow.Key) {
	switch key {
	case glow.KeyEnter:
		// Default: 1 human + 1 AI
		players := []*player.Player{
			player.NewPlayer(0, "Player 1", false),
			player.NewPlayer(1, "AI Player", true),
		}
		g.StartGame(players)
	case glow.KeyR:
		if save.HasSave() {
			g.loadGame()
		}
	case glow.Key2:
		players := []*player.Player{
			player.NewPlayer(0, "Player 1", false),
			player.NewPlayer(1, "AI Player", true),
		}
		g.StartGame(players)
	case glow.Key3:
		players := []*player.Player{
			player.NewPlayer(0, "Player 1", false),
			player.NewPlayer(1, "AI Player 1", true),
			player.NewPlayer(2, "AI Player 2", true),
		}
		g.StartGame(players)
	case glow.Key4:
		players := []*player.Player{
			player.NewPlayer(0, "Player 1", false),
			player.NewPlayer(1, "AI Player 1", true),
			player.NewPlayer(2, "AI Player 2", true),
			player.NewPlayer(3, "AI Player 3", true),
		}
		g.StartGame(players)
	}
}

func (g *Game) keySetup(key glow.Key) {}

func (g *Game) keyPlaying(key glow.Key) {
	if key == glow.KeyF5 {
		g.saveGame()
	}
}

// saveGame serialises current game state to disk.
func (g *Game) saveGame() {
	data := &save.SaveData{
		Current:    g.Current,
		Properties: save.BoardToPropertyData(g.Board),
		HousePool:  g.Board.HousePool,
		HotelPool:  g.Board.HotelPool,
		Die1:       g.Die1,
		Die2:       g.Die2,
		Messages:   g.Messages,
	}

	for _, p := range g.Players {
		data.Players = append(data.Players, save.PlayerData{
			ID:                p.ID,
			Name:              p.Name,
			IsAI:              p.IsAI,
			Money:             p.Money,
			Position:          p.Position,
			InJail:            p.InJail,
			JailTurns:         p.JailTurns,
			Bankrupt:          p.Bankrupt,
			Properties:        p.Properties,
			GetOutOfJailCards: p.GetOutOfJailCards,
		})
	}

	if err := save.Save(data); err != nil {
		g.AddMessage("Save failed: " + err.Error())
	} else {
		g.AddMessage("Game saved!")
	}
}

// loadGame restores game state from disk.
func (g *Game) loadGame() bool {
	data, err := save.Load()
	if err != nil {
		g.AddMessage("Load failed: " + err.Error())
		return false
	}

	g.Board = board.NewBoard()
	save.PropertyDataToBoard(g.Board, data.Properties)
	g.Board.HousePool = data.HousePool
	g.Board.HotelPool = data.HotelPool

	g.Players = nil
	for _, pd := range data.Players {
		p := player.NewPlayer(pd.ID, pd.Name, pd.IsAI)
		p.Money = pd.Money
		p.Position = pd.Position
		p.InJail = pd.InJail
		p.JailTurns = pd.JailTurns
		p.Bankrupt = pd.Bankrupt
		p.Properties = pd.Properties
		p.GetOutOfJailCards = pd.GetOutOfJailCards
		g.Players = append(g.Players, p)
	}

	g.Current = data.Current
	g.Die1 = data.Die1
	g.Die2 = data.Die2
	g.Messages = data.Messages
	g.State = StatePlaying
	g.Phase = PhasePreRoll
	g.Dialog = DialogNone
	g.setupButtons()
	g.updateButtonStates()
	g.AddMessage("Game loaded!")
	return true
}

// hudData builds an HUDData struct for the renderer.
func (g *Game) hudData() render.HUDData {
	data := render.HUDData{
		CurrentPlayerID: g.currentPlayer().ID,
		Messages:        g.Messages,
		Die1:            g.Die1,
		Die2:            g.Die2,
		Phase:           g.phaseString(),
	}
	for _, p := range g.Players {
		data.Players = append(data.Players, render.PlayerInfo{
			ID:       p.ID,
			Name:     p.Name,
			Money:    p.Money,
			Bankrupt: p.Bankrupt,
			InJail:   p.InJail,
			IsAI:     p.IsAI,
		})
	}
	return data
}

// tradePropsContains checks if a space index is in a property list.
func (g *Game) tradePropsContains(props []int, idx int) bool {
	for _, p := range props {
		if p == idx {
			return true
		}
	}
	return false
}

// tradePropsToggle adds or removes a space index from a property list.
func (g *Game) tradePropsToggle(props []int, idx int) []int {
	for i, p := range props {
		if p == idx {
			return append(props[:i], props[i+1:]...)
		}
	}
	return append(props, idx)
}

func (g *Game) phaseString() string {
	switch g.Phase {
	case PhasePreRoll:
		return "Click [Roll Dice]"
	case PhaseRolling:
		return "Rolling..."
	case PhaseMoving:
		return "Moving..."
	case PhaseLanded:
		return "Landed!"
	case PhaseDialog:
		return "Decision time"
	case PhaseAuction:
		return "Auction"
	case PhasePostAction:
		return "Post-action"
	case PhaseTurnEnd:
		return "Turn ending"
	case PhaseJailDecision:
		return "Jail decision"
	default:
		return ""
	}
}
