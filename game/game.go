package game

import (
	"fmt"

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
	TradePartner     int // target player index
	TradeOfferedProps []int
	TradeWantedProps  []int
	TradeOfferedMoney int
	TradeWantedMoney  int

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
	bx := g.Layout.PanelX + 20
	by := 360
	bw := 90
	bh := 28
	gap := 6

	positions := [][4]int{
		{bx, by, bw*2 + gap, bh},                 // Roll Dice
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
	bx := g.Layout.PanelX + 20
	by := 360
	bw := 90
	bh := 28
	gap := 6

	g.Buttons = []render.Button{
		render.NewButton("Roll Dice", bx, by, bw*2+gap, bh),
		render.NewButton("Buy", bx, by+bh+gap, bw, bh),
		render.NewButton("Auction", bx+bw+gap, by+bh+gap, bw, bh),
		render.NewButton("Build", bx, by+2*(bh+gap), bw, bh),
		render.NewButton("Mortgage", bx+bw+gap, by+2*(bh+gap), bw, bh),
		render.NewButton("Trade", bx, by+3*(bh+gap), bw, bh),
		render.NewButton("End Turn", bx+bw+gap, by+3*(bh+gap), bw, bh),
	}
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
				groupCol := render.SpaceBg
				if space.Type == board.SpaceProperty {
					groupCol = render.GroupColor(space.Group)
				}
				// Draw card on the HUD panel
				render.DrawPropertyCard(canvas,
					g.Layout.PanelX+20, g.Layout.WinH-180, g.Layout.PanelWidth-40,
					space.Name, space.Price, space.Rent, space.HouseCost,
					groupCol, ownerName, prop.Houses, prop.Mortgaged)
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
		if g.TradePartner < 0 {
			// Select trade partner
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
		} else {
			// Simple trade: select one of your properties to offer
			partner := g.Players[g.TradePartner]
			var btns []render.DialogButton
			for _, idx := range p.Properties {
				space := g.Board.Spaces[idx]
				label := fmt.Sprintf("Offer: %s", space.Name)
				btns = append(btns, render.DialogButton{Label: label, ID: idx, Enabled: !g.Board.Properties[idx].Mortgaged && g.Board.Properties[idx].Houses == 0})
			}
			// Also show partner's properties we could request
			for _, idx := range partner.Properties {
				space := g.Board.Spaces[idx]
				label := fmt.Sprintf("Want: %s", space.Name)
				btns = append(btns, render.DialogButton{Label: label, ID: 1000 + idx, Enabled: !g.Board.Properties[idx].Mortgaged && g.Board.Properties[idx].Houses == 0})
			}
			btns = append(btns, render.DialogButton{Label: "Cancel", ID: -1, Enabled: true})
			data := render.DialogData{
				Title:   fmt.Sprintf("Trade with %s", partner.Name),
				Lines:   []string{"Select properties to trade:"},
				Buttons: btns,
			}
			g.DialogHovered = render.DrawDialog(canvas, data, g.MouseX, g.MouseY)
		}
	}
}

func (g *Game) drawGameOver(canvas *glow.Canvas) {
	cx := canvas.Width() / 2
	render.DrawTextCentered(canvas, "GAME OVER", cx, 300, render.TextGold, 4)
	if len(g.alivePlayers()) == 1 {
		winner := g.alivePlayers()[0]
		render.DrawTextCentered(canvas, winner.Name+" WINS!", cx, 380, render.PlayerColors[winner.ID%4], 3)
	}
	render.DrawTextCentered(canvas, "Press ENTER to return to menu", cx, 460, render.TextLight, 1)
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
