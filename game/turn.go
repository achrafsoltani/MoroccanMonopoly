package game

import (
	"fmt"
	"math/rand"

	"github.com/AchrafSoltani/MoroccanMonopoly/board"
	"github.com/AchrafSoltani/MoroccanMonopoly/config"
	"github.com/AchrafSoltani/MoroccanMonopoly/player"
)

// updatePlaying handles the main gameplay loop.
func (g *Game) updatePlaying(dt float64) {
	p := g.currentPlayer()
	if p == nil {
		return
	}

	// AI auto-actions
	if p.IsAI {
		g.updateAI(dt)
		return
	}

	switch g.Phase {
	case PhasePreRoll:
		// If in jail, show jail options
		if p.InJail {
			g.Phase = PhaseJailDecision
			g.Dialog = DialogJailOptions
		}
	case PhaseJailDecision:
		// Waiting for player input (handled in button clicks)
	case PhaseRolling:
		g.updateDiceAnim(dt)
	case PhaseMoving:
		g.updateTokenMove(dt)
	case PhaseLanded:
		g.resolveLanding()
	case PhaseTurnEnd:
		g.endTurn()
	}

	// Handle button clicks
	if g.MouseClicked {
		g.handleDialogClicks()
		g.handleButtonClicks()
	}
}

// handleButtonClicks processes clicks on the action buttons.
func (g *Game) handleButtonClicks() {
	for i, btn := range g.Buttons {
		if !btn.Visible || !btn.Enabled || !btn.Contains(g.MouseX, g.MouseY) {
			continue
		}
		switch i {
		case 0: // Roll Dice
			if g.Phase == PhasePreRoll {
				g.startDiceRoll()
			}
		case 1: // Buy
			if g.Phase == PhaseDialog && g.Dialog == DialogBuyProperty {
				g.buyProperty()
			}
		case 2: // Auction
			if g.Phase == PhaseDialog && g.Dialog == DialogBuyProperty {
				g.declineBuy()
			}
		case 3: // Build
			if g.Phase == PhasePreRoll || g.Phase == PhasePostAction {
				g.openBuildDialog()
			}
		case 4: // Mortgage
			if g.Phase == PhasePreRoll || g.Phase == PhasePostAction {
				g.openMortgageDialog()
			}
		case 5: // Trade
			if g.Phase == PhasePreRoll || g.Phase == PhasePostAction {
				g.openTradeDialog()
			}
		case 6: // End Turn
			if g.Phase == PhasePostAction {
				g.endTurn()
			}
		}
	}
}

// handleDialogClicks processes clicks on dialog buttons.
func (g *Game) handleDialogClicks() {
	if g.DialogHovered < 0 {
		return
	}

	switch g.Dialog {
	case DialogBuyProperty:
		switch g.DialogHovered {
		case 0: // Buy
			g.buyProperty()
		case 1: // Decline/Auction
			g.declineBuy()
		}

	case DialogJailOptions:
		p := g.currentPlayer()
		switch g.DialogHovered {
		case 0: // Pay fine
			if p.Money >= config.JailFine {
				p.Pay(config.JailFine)
				p.InJail = false
				p.JailTurns = 0
				g.Dialog = DialogNone
				g.Phase = PhasePreRoll
				g.AddMessage(fmt.Sprintf("%s paid %d MAD to get out of jail", p.Name, config.JailFine))
			}
		case 1: // Use card
			if p.GetOutOfJailCards > 0 {
				p.GetOutOfJailCards--
				p.InJail = false
				p.JailTurns = 0
				g.Dialog = DialogNone
				g.Phase = PhasePreRoll
				g.AddMessage(fmt.Sprintf("%s used Get Out of Jail Free card", p.Name))
			}
		case 2: // Roll doubles
			g.Dialog = DialogNone
			g.startDiceRoll()
		}

	case DialogBuild:
		if g.DialogHovered == -1 {
			// Cancel
			g.Dialog = DialogNone
			g.Phase = PhasePostAction
			g.updateButtonStates()
		} else if g.DialogHovered >= 0 && g.DialogHovered < len(g.SelectableSpaces) {
			idx := g.SelectableSpaces[g.DialogHovered]
			p := g.currentPlayer()
			space := g.Board.Spaces[idx]
			cost := g.BuildHouse(idx)
			p.Pay(cost)
			level := g.Board.Properties[idx].Houses
			levelName := fmt.Sprintf("%d house(s)", level)
			if level == config.HotelLevel {
				levelName = "hotel"
			}
			g.AddMessage(fmt.Sprintf("%s built on %s (%s, -%d MAD)", p.Name, space.Name, levelName, cost))
			// Refresh buildable list
			buildable := g.BuildableProperties(p.ID)
			if len(buildable) == 0 {
				g.Dialog = DialogNone
				g.Phase = PhasePostAction
				g.updateButtonStates()
			} else {
				g.SelectableSpaces = buildable
			}
		}

	case DialogMortgage:
		if g.DialogHovered == -1 {
			g.Dialog = DialogNone
			g.Phase = PhasePostAction
			g.updateButtonStates()
		} else if g.DialogHovered >= 0 && g.DialogHovered < len(g.SelectableSpaces) {
			idx := g.SelectableSpaces[g.DialogHovered]
			p := g.currentPlayer()
			space := g.Board.Spaces[idx]
			prop := g.Board.Properties[idx]
			if prop.Mortgaged {
				cost := g.UnmortgageProperty(idx)
				p.Pay(cost)
				g.AddMessage(fmt.Sprintf("%s unmortgaged %s (-%d MAD)", p.Name, space.Name, cost))
			} else {
				val := g.MortgageProperty(idx)
				p.Receive(val)
				g.AddMessage(fmt.Sprintf("%s mortgaged %s (+%d MAD)", p.Name, space.Name, val))
			}
			// Refresh list
			mortgageable := g.MortgageableProperties(p.ID)
			unmortgageable := g.UnmortgageableProperties(p.ID)
			all := append(mortgageable, unmortgageable...)
			if len(all) == 0 {
				g.Dialog = DialogNone
				g.Phase = PhasePostAction
				g.updateButtonStates()
			} else {
				g.SelectableSpaces = all
			}
		}

	case DialogAuction:
		g.handleAuctionClick()

	case DialogTrade:
		if g.DialogHovered == -1 {
			g.Dialog = DialogNone
			g.Phase = PhasePostAction
			g.TradePartner = -1
			g.updateButtonStates()
		} else if g.TradePartner < 0 {
			// Selected a partner
			g.TradePartner = g.DialogHovered
		} else {
			// Selected a property to trade
			idx := g.DialogHovered
			if idx >= 1000 {
				// Want a property from partner
				actualIdx := idx - 1000
				offer := TradeOffer{
					FromPlayer:  g.currentPlayer().ID,
					ToPlayer:    g.TradePartner,
					WantedProps: []int{actualIdx},
				}
				partner := g.Players[g.TradePartner]
				if partner.IsAI {
					if g.aiEvaluateTrade(offer) {
						g.executeTrade(offer)
					} else {
						g.AddMessage(fmt.Sprintf("%s declined the trade", partner.Name))
					}
				} else {
					g.executeTrade(offer)
				}
				g.Dialog = DialogNone
				g.Phase = PhasePostAction
				g.TradePartner = -1
				g.updateButtonStates()
			} else {
				// Offer one of own properties
				offer := TradeOffer{
					FromPlayer:   g.currentPlayer().ID,
					ToPlayer:     g.TradePartner,
					OfferedProps:  []int{idx},
				}
				partner := g.Players[g.TradePartner]
				if partner.IsAI {
					if g.aiEvaluateTrade(offer) {
						g.executeTrade(offer)
					} else {
						g.AddMessage(fmt.Sprintf("%s declined the trade", partner.Name))
					}
				} else {
					g.executeTrade(offer)
				}
				g.Dialog = DialogNone
				g.Phase = PhasePostAction
				g.TradePartner = -1
				g.updateButtonStates()
			}
		}
	}
}

// startDiceRoll begins the dice animation.
func (g *Game) startDiceRoll() {
	g.Phase = PhaseRolling
	g.DiceRolling = true
	g.DiceAnimTimer = 0
	g.Die1 = rand.Intn(6) + 1
	g.Die2 = rand.Intn(6) + 1
	g.Audio.PlayDiceRoll()
}

// updateDiceAnim advances the dice rolling animation.
func (g *Game) updateDiceAnim(dt float64) {
	g.DiceAnimTimer += dt
	if g.DiceAnimTimer >= config.DiceAnimDuration {
		g.DiceRolling = false
		g.Doubles = g.Die1 == g.Die2
		if g.Doubles {
			g.DoublesCount++
		}

		p := g.currentPlayer()
		total := g.Die1 + g.Die2
		g.AddMessage(fmt.Sprintf("%s rolled %d + %d = %d", p.Name, g.Die1, g.Die2, total))

		if g.Doubles {
			g.AddMessage("Doubles!")
		}

		// Check for 3 consecutive doubles
		if g.DoublesCount >= 3 {
			g.AddMessage(fmt.Sprintf("%s: 3 doubles! Go to jail!", p.Name))
			g.sendToJail(p)
			g.Phase = PhasePostAction
			return
		}

		// Handle jail
		if p.InJail {
			if g.Doubles {
				p.InJail = false
				p.JailTurns = 0
				g.AddMessage(fmt.Sprintf("%s rolled doubles and is free!", p.Name))
			} else {
				p.JailTurns++
				if p.JailTurns >= config.MaxJailTurns {
					p.Pay(config.JailFine)
					p.InJail = false
					p.JailTurns = 0
					g.AddMessage(fmt.Sprintf("%s paid %d MAD jail fine (forced)", p.Name, config.JailFine))
				} else {
					g.AddMessage(fmt.Sprintf("%s stays in jail (%d/3 turns)", p.Name, p.JailTurns))
					g.Phase = PhasePostAction
					return
				}
			}
		}

		// Start movement
		g.startMove(p.Position, total)
	}
}

// startMove initiates token movement animation.
func (g *Game) startMove(from, steps int) {
	g.Phase = PhaseMoving
	g.MoveFrom = from
	g.MoveSteps = steps
	g.MoveCurrent = 0
	g.MoveTimer = 0
}

// updateTokenMove advances the token movement animation.
func (g *Game) updateTokenMove(dt float64) {
	g.MoveTimer += dt
	if g.MoveTimer >= config.TokenMoveDuration {
		g.MoveTimer -= config.TokenMoveDuration
		g.MoveCurrent++

		p := g.currentPlayer()
		prevPos := (g.MoveFrom + g.MoveCurrent - 1) % config.SpaceCount
		newPos := (g.MoveFrom + g.MoveCurrent) % config.SpaceCount
		p.Position = newPos

		// Check for passing GO (crossed from 39 to 0)
		if newPos < prevPos && g.MoveCurrent < g.MoveSteps {
			p.Receive(config.GoSalary)
			g.AddMessage(fmt.Sprintf("%s passed GO! +%d MAD", p.Name, config.GoSalary))
			g.Audio.PlayPassGo()
		}

		if g.MoveCurrent >= g.MoveSteps {
			g.MoveTo = newPos

			// Check if passed GO on final step
			if newPos < prevPos {
				p.Receive(config.GoSalary)
				g.AddMessage(fmt.Sprintf("%s passed GO! +%d MAD", p.Name, config.GoSalary))
				g.Audio.PlayPassGo()
			}

			g.Phase = PhaseLanded
		}
	}
}

// resolveLanding handles what happens when a player lands on a space.
func (g *Game) resolveLanding() {
	p := g.currentPlayer()
	space := g.Board.Spaces[p.Position]
	g.AddMessage(fmt.Sprintf("%s landed on %s", p.Name, space.Name))

	switch space.Type {
	case board.SpaceGo:
		// Already collected when passing
		g.Phase = PhasePostAction

	case board.SpaceProperty, board.SpaceRailroad, board.SpaceUtility:
		prop := g.Board.Properties[p.Position]
		if prop.OwnerID < 0 {
			// Unowned — offer to buy
			g.Dialog = DialogBuyProperty
			g.Phase = PhaseDialog
			g.AddMessage(fmt.Sprintf("Buy %s for %d MAD?", space.Name, space.Price))
			g.updateButtonStates()
		} else if prop.OwnerID != p.ID {
			// Owned by someone else — pay rent
			rent := g.calculateRent(p.Position)
			if prop.Mortgaged {
				g.AddMessage(fmt.Sprintf("%s is mortgaged - no rent", space.Name))
			} else {
				owner := g.Players[prop.OwnerID]
				g.AddMessage(fmt.Sprintf("%s pays %d MAD rent to %s", p.Name, rent, owner.Name))
				g.Audio.PlayRent()
				g.payDebt(p, owner, rent)
			}
			g.Phase = PhasePostAction
		} else {
			// Own property
			g.Phase = PhasePostAction
		}

	case board.SpaceChance:
		g.drawChanceCard()

	case board.SpaceCommunityChest:
		g.drawCommunityCard()

	case board.SpaceTax:
		g.AddMessage(fmt.Sprintf("%s pays %d MAD tax", p.Name, space.TaxAmount))
		g.payDebt(p, nil, space.TaxAmount) // nil = bank
		g.Phase = PhasePostAction

	case board.SpaceJail:
		g.AddMessage(fmt.Sprintf("%s is just visiting", p.Name))
		g.Phase = PhasePostAction

	case board.SpaceFreeParking:
		g.AddMessage("Free parking - nothing happens")
		g.Phase = PhasePostAction

	case board.SpaceGoToJail:
		g.AddMessage(fmt.Sprintf("%s goes to jail!", p.Name))
		g.sendToJail(p)
		g.Phase = PhasePostAction
	}
}

// buyProperty handles the player buying the current property.
func (g *Game) buyProperty() {
	p := g.currentPlayer()
	space := g.Board.Spaces[p.Position]

	if p.Money < space.Price {
		g.AddMessage("Not enough money!")
		return
	}

	p.Pay(space.Price)
	p.AddProperty(p.Position)
	g.Board.Properties[p.Position].OwnerID = p.ID
	g.AddMessage(fmt.Sprintf("%s bought %s for %d MAD", p.Name, space.Name, space.Price))
	g.Audio.PlayPurchase()
	g.Dialog = DialogNone
	g.Phase = PhasePostAction
	g.updateButtonStates()
}

// declineBuy handles the player declining to buy — triggers an auction.
func (g *Game) declineBuy() {
	p := g.currentPlayer()
	g.AddMessage(fmt.Sprintf("%s declined to buy", p.Name))
	g.Dialog = DialogNone
	g.startAuction(p.Position)
}

// sendToJail moves a player to jail.
func (g *Game) sendToJail(p *player.Player) {
	p.Position = config.JailPosition
	p.InJail = true
	p.JailTurns = 0
	g.Audio.PlayJail()
}

// endTurn finishes the current turn and advances to the next player.
func (g *Game) endTurn() {
	// Check if doubles — roll again
	if g.Doubles && !g.currentPlayer().InJail {
		g.Phase = PhasePreRoll
		g.AddMessage(fmt.Sprintf("%s rolls again (doubles)", g.currentPlayer().Name))
		g.Die1 = 0
		g.Die2 = 0
		g.updateButtonStates()
		return
	}

	g.nextPlayer()
	g.Phase = PhasePreRoll
	g.Die1 = 0
	g.Die2 = 0
	g.Doubles = false
	g.AddMessage(fmt.Sprintf("--- %s's turn ---", g.currentPlayer().Name))
	g.updateButtonStates()

	// Check win condition
	if len(g.alivePlayers()) <= 1 {
		g.State = StateGameOver
	}
}

// calculateRent computes rent for a property.
func (g *Game) calculateRent(spaceIndex int) int {
	space := g.Board.Spaces[spaceIndex]
	prop := g.Board.Properties[spaceIndex]

	if prop.Mortgaged {
		return 0
	}

	switch space.Type {
	case board.SpaceProperty:
		houses := prop.Houses
		if houses > 0 {
			return space.Rent[houses]
		}
		// Check for monopoly (doubles base rent)
		if g.hasMonopoly(prop.OwnerID, space.Group) {
			return space.Rent[0] * 2
		}
		return space.Rent[0]

	case board.SpaceRailroad:
		count := g.countOwnedRailroads(prop.OwnerID)
		rents := [4]int{25, 50, 100, 200}
		if count >= 1 && count <= 4 {
			return rents[count-1]
		}
		return 25

	case board.SpaceUtility:
		count := g.countOwnedUtilities(prop.OwnerID)
		diceTotal := g.Die1 + g.Die2
		if count == 2 {
			return diceTotal * 10
		}
		return diceTotal * 4
	}

	return 0
}

// hasMonopoly checks if a player owns all properties in a colour group.
func (g *Game) hasMonopoly(playerID int, group board.ColorGroup) bool {
	if group == board.GroupNone {
		return false
	}
	spaces := g.Board.SpacesInGroup(group)
	for _, idx := range spaces {
		if g.Board.Properties[idx].OwnerID != playerID {
			return false
		}
	}
	return true
}

func (g *Game) countOwnedRailroads(playerID int) int {
	count := 0
	for _, idx := range g.Board.RailroadSpaces() {
		if g.Board.Properties[idx].OwnerID == playerID {
			count++
		}
	}
	return count
}

func (g *Game) countOwnedUtilities(playerID int) int {
	count := 0
	for _, idx := range g.Board.UtilitySpaces() {
		if g.Board.Properties[idx].OwnerID == playerID {
			count++
		}
	}
	return count
}

// updateButtonStates enables/disables buttons based on current phase.
func (g *Game) updateButtonStates() {
	if len(g.Buttons) < 7 {
		return
	}

	// Roll Dice
	g.Buttons[0].Enabled = g.Phase == PhasePreRoll
	// Buy
	g.Buttons[1].Enabled = g.Phase == PhaseDialog && g.Dialog == DialogBuyProperty
	g.Buttons[1].Visible = g.Phase == PhaseDialog && g.Dialog == DialogBuyProperty
	// Auction
	g.Buttons[2].Enabled = g.Phase == PhaseDialog && g.Dialog == DialogBuyProperty
	g.Buttons[2].Visible = g.Phase == PhaseDialog && g.Dialog == DialogBuyProperty
	// Build
	g.Buttons[3].Enabled = g.Phase == PhasePreRoll || g.Phase == PhasePostAction
	// Mortgage
	g.Buttons[4].Enabled = g.Phase == PhasePreRoll || g.Phase == PhasePostAction
	// Trade
	g.Buttons[5].Enabled = g.Phase == PhasePreRoll || g.Phase == PhasePostAction
	// End Turn
	g.Buttons[6].Enabled = g.Phase == PhasePostAction
}

// drawChanceCard draws and executes a Chance card.
func (g *Game) drawChanceCard() {
	card := g.Board.ChanceDeck.Draw()
	g.Audio.PlayCardDraw()
	g.AddMessage("Chance: " + card.Text)
	g.executeCard(card)
}

// drawCommunityCard draws and executes a Community Chest card.
func (g *Game) drawCommunityCard() {
	card := g.Board.CommunityDeck.Draw()
	g.Audio.PlayCardDraw()
	g.AddMessage("Caisse: " + card.Text)
	g.executeCard(card)
}

// executeCard applies a card's effect.
func (g *Game) executeCard(card board.Card) {
	p := g.currentPlayer()

	switch card.Effect {
	case board.EffectCollect:
		p.Receive(card.Amount)
		g.AddMessage(fmt.Sprintf("%s receives %d MAD", p.Name, card.Amount))

	case board.EffectPay:
		p.Pay(card.Amount)
		g.AddMessage(fmt.Sprintf("%s pays %d MAD", p.Name, card.Amount))

	case board.EffectMoveTo:
		target := card.Amount
		// Check if passing GO
		if target < p.Position && target != 0 {
			p.Receive(config.GoSalary)
			g.AddMessage(fmt.Sprintf("%s passed GO! +%d MAD", p.Name, config.GoSalary))
		} else if target == 0 {
			p.Receive(config.GoSalary)
			g.AddMessage(fmt.Sprintf("%s collects %d MAD from DEPART", p.Name, config.GoSalary))
		}
		p.Position = target
		g.Phase = PhaseLanded
		return

	case board.EffectMoveSteps:
		steps := card.Amount
		newPos := (p.Position + steps + config.SpaceCount) % config.SpaceCount
		p.Position = newPos
		g.Phase = PhaseLanded
		return

	case board.EffectGoToJail:
		g.sendToJail(p)

	case board.EffectGetOutOfJail:
		p.GetOutOfJailCards++
		g.AddMessage(fmt.Sprintf("%s gets a Get Out of Jail Free card!", p.Name))

	case board.EffectPayPerHouse:
		totalHouses := 0
		totalHotels := 0
		for _, idx := range p.Properties {
			h := g.Board.Properties[idx].Houses
			if h == 5 {
				totalHotels++
			} else {
				totalHouses += h
			}
		}
		cost := totalHouses*card.Amount + totalHotels*card.AmountHotel
		p.Pay(cost)
		g.AddMessage(fmt.Sprintf("%s pays %d MAD (%d houses, %d hotels)", p.Name, cost, totalHouses, totalHotels))

	case board.EffectCollectAll:
		total := 0
		for _, other := range g.Players {
			if other.ID != p.ID && !other.Bankrupt {
				other.Pay(card.Amount)
				total += card.Amount
			}
		}
		p.Receive(total)
		g.AddMessage(fmt.Sprintf("%s collects %d MAD from all players", p.Name, total))

	case board.EffectPayAll:
		for _, other := range g.Players {
			if other.ID != p.ID && !other.Bankrupt {
				p.Pay(card.Amount)
				other.Receive(card.Amount)
			}
		}
		total := card.Amount * (len(g.alivePlayers()) - 1)
		g.AddMessage(fmt.Sprintf("%s pays %d MAD total to all players", p.Name, total))
	}

	g.Phase = PhasePostAction
}

// payDebt attempts to pay a debt. If unable, triggers bankruptcy.
// creditor is nil when paying the bank.
func (g *Game) payDebt(debtor *player.Player, creditor *player.Player, amount int) {
	if debtor.Money >= amount {
		debtor.Pay(amount)
		if creditor != nil {
			creditor.Receive(amount)
		}
		return
	}

	// Try to auto-liquidate: sell houses, then mortgage
	g.autoLiquidate(debtor)

	if debtor.Money >= amount {
		debtor.Pay(amount)
		if creditor != nil {
			creditor.Receive(amount)
		}
		return
	}

	// Bankrupt
	g.declareBankruptcy(debtor, creditor)
}

// autoLiquidate sells houses and mortgages properties to raise cash.
func (g *Game) autoLiquidate(p *player.Player) {
	// First: sell all houses/hotels
	for {
		sold := false
		for _, idx := range p.Properties {
			if g.Board.Properties[idx].Houses > 0 && g.CanSellHouseOnSpace(idx) {
				refund := g.SellHouse(idx)
				p.Receive(refund)
				space := g.Board.Spaces[idx]
				g.AddMessage(fmt.Sprintf("%s sold house on %s (+%d MAD)", p.Name, space.Name, refund))
				sold = true
			}
		}
		if !sold {
			break
		}
	}

	// Then: mortgage unimproved properties
	for _, idx := range p.Properties {
		prop := g.Board.Properties[idx]
		if !prop.Mortgaged && prop.Houses == 0 {
			val := g.MortgageProperty(idx)
			p.Receive(val)
			space := g.Board.Spaces[idx]
			g.AddMessage(fmt.Sprintf("%s mortgaged %s (+%d MAD)", p.Name, space.Name, val))
		}
	}
}

// declareBankruptcy eliminates a player and transfers assets.
func (g *Game) declareBankruptcy(debtor *player.Player, creditor *player.Player) {
	g.AddMessage(fmt.Sprintf("%s is BANKRUPT!", debtor.Name))
	g.Audio.PlayBankruptcy()
	debtor.Bankrupt = true

	if creditor != nil {
		// Transfer all assets to creditor
		creditor.Receive(debtor.Money)
		for _, idx := range debtor.Properties {
			creditor.AddProperty(idx)
			g.Board.Properties[idx].OwnerID = creditor.ID
		}
		creditor.GetOutOfJailCards += debtor.GetOutOfJailCards
		g.AddMessage(fmt.Sprintf("%s receives all of %s's assets", creditor.Name, debtor.Name))
	} else {
		// Owed to bank — return properties to bank (unowned), auction them
		for _, idx := range debtor.Properties {
			g.Board.Properties[idx].OwnerID = -1
			g.Board.Properties[idx].Mortgaged = false
			g.Board.Properties[idx].Houses = 0
		}
	}

	debtor.Money = 0
	debtor.Properties = nil
	debtor.GetOutOfJailCards = 0

	// Check win condition
	alive := g.alivePlayers()
	if len(alive) <= 1 {
		g.State = StateGameOver
	}
}

// openBuildDialog opens the build house dialog.
func (g *Game) openBuildDialog() {
	p := g.currentPlayer()
	buildable := g.BuildableProperties(p.ID)
	if len(buildable) == 0 {
		g.AddMessage("No properties available to build on")
		return
	}
	g.SelectableSpaces = buildable
	g.SelectedSpace = -1
	g.Dialog = DialogBuild
	g.Phase = PhaseDialog
}

// openMortgageDialog opens the mortgage dialog.
func (g *Game) openMortgageDialog() {
	p := g.currentPlayer()
	mortgageable := g.MortgageableProperties(p.ID)
	unmortgageable := g.UnmortgageableProperties(p.ID)
	all := append(mortgageable, unmortgageable...)
	if len(all) == 0 {
		g.AddMessage("No properties available to mortgage/unmortgage")
		return
	}
	g.SelectableSpaces = all
	g.SelectedSpace = -1
	g.Dialog = DialogMortgage
	g.Phase = PhaseDialog
}

// updateAI handles AI player turns with automatic actions.
func (g *Game) updateAI(dt float64) {
	p := g.currentPlayer()

	// Add delay for AI actions so humans can follow
	if g.Phase == PhasePreRoll || g.Phase == PhasePostAction || g.Phase == PhaseJailDecision || g.Phase == PhaseDialog {
		g.AITimer += dt
		if g.AITimer < 0.5 {
			return
		}
		g.AITimer = 0
	}

	switch g.Phase {
	case PhasePreRoll:
		// AI tries to build before rolling
		g.aiBuildIfPossible()

		if p.InJail {
			g.Phase = PhaseJailDecision
			g.Dialog = DialogJailOptions
		} else {
			g.startDiceRoll()
		}
	case PhaseJailDecision:
		if p.GetOutOfJailCards > 0 {
			p.GetOutOfJailCards--
			p.InJail = false
			p.JailTurns = 0
			g.Dialog = DialogNone
			g.Phase = PhasePreRoll
			g.AddMessage(fmt.Sprintf("%s (AI) used Get Out of Jail Free card", p.Name))
		} else if p.Money >= config.JailFine+200 {
			p.Pay(config.JailFine)
			p.InJail = false
			p.JailTurns = 0
			g.Dialog = DialogNone
			g.Phase = PhasePreRoll
			g.AddMessage(fmt.Sprintf("%s (AI) paid %d MAD jail fine", p.Name, config.JailFine))
		} else {
			g.Dialog = DialogNone
			g.startDiceRoll()
		}
	case PhaseRolling:
		g.updateDiceAnim(dt)
	case PhaseMoving:
		g.updateTokenMove(dt)
	case PhaseLanded:
		g.resolveLanding()
	case PhaseDialog:
		if g.Dialog == DialogBuyProperty {
			space := g.Board.Spaces[p.Position]
			// Count total owned to adjust strategy
			totalOwned := 0
			for _, pl := range g.Players {
				totalOwned += len(pl.Properties)
			}
			if p.ShouldBuy(space.Price, totalOwned) {
				g.buyProperty()
			} else {
				g.declineBuy()
			}
		}
	case PhasePostAction:
		g.endTurn()
	case PhaseTurnEnd:
		g.endTurn()
	}
}

// aiBuildIfPossible has the AI build houses if it can.
func (g *Game) aiBuildIfPossible() {
	p := g.currentPlayer()
	buildable := g.BuildableProperties(p.ID)
	for _, idx := range buildable {
		space := g.Board.Spaces[idx]
		if p.ShouldBuild(space.HouseCost) {
			cost := g.BuildHouse(idx)
			p.Pay(cost)
			level := g.Board.Properties[idx].Houses
			levelName := fmt.Sprintf("%d house(s)", level)
			if level == config.HotelLevel {
				levelName = "hotel"
			}
			g.AddMessage(fmt.Sprintf("%s (AI) built on %s (%s)", p.Name, space.Name, levelName))
		}
	}
}

