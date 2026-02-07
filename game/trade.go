package game

import (
	"fmt"

	"github.com/AchrafSoltani/MoroccanMonopoly/board"
)

// TradeOffer represents a trade proposal.
type TradeOffer struct {
	FromPlayer int
	ToPlayer   int
	OfferedProps []int // space indices offered
	WantedProps  []int // space indices wanted
	OfferedMoney int
	WantedMoney  int
	OfferedJailCards int
	WantedJailCards  int
}

// TradeState tracks the trading UI state.
type TradeState int

const (
	TradeSelectPartner TradeState = iota
	TradeSelectOffer
	TradeConfirm
)

// openTradeDialog starts the trade flow.
func (g *Game) openTradeDialog() {
	p := g.currentPlayer()
	// Find other alive players
	var partners []int
	for _, other := range g.Players {
		if other.ID != p.ID && !other.Bankrupt {
			partners = append(partners, other.ID)
		}
	}
	if len(partners) == 0 {
		g.AddMessage("No players to trade with")
		return
	}

	g.TradePartner = -1
	g.TradeOfferedProps = nil
	g.TradeWantedProps = nil
	g.TradeOfferedMoney = 0
	g.TradeWantedMoney = 0
	g.TradeOfferJailCard = false
	g.TradeWantJailCard = false
	g.TradeStage = TradeSelectPartner
	g.PendingOffer = nil
	g.Dialog = DialogTrade
	g.Phase = PhaseTrade

	g.SelectableSpaces = partners
	g.SelectedSpace = -1
}

// executeTrade performs the trade between two players.
func (g *Game) executeTrade(offer TradeOffer) {
	from := g.Players[offer.FromPlayer]
	to := g.Players[offer.ToPlayer]

	// Transfer offered properties
	for _, idx := range offer.OfferedProps {
		from.RemoveProperty(idx)
		to.AddProperty(idx)
		g.Board.Properties[idx].OwnerID = offer.ToPlayer
	}

	// Transfer wanted properties
	for _, idx := range offer.WantedProps {
		to.RemoveProperty(idx)
		from.AddProperty(idx)
		g.Board.Properties[idx].OwnerID = offer.FromPlayer
	}

	// Transfer money
	if offer.OfferedMoney > 0 {
		from.Pay(offer.OfferedMoney)
		to.Receive(offer.OfferedMoney)
	}
	if offer.WantedMoney > 0 {
		to.Pay(offer.WantedMoney)
		from.Receive(offer.WantedMoney)
	}

	// Transfer jail cards
	from.GetOutOfJailCards -= offer.OfferedJailCards
	to.GetOutOfJailCards += offer.OfferedJailCards
	to.GetOutOfJailCards -= offer.WantedJailCards
	from.GetOutOfJailCards += offer.WantedJailCards

	g.AddMessage(fmt.Sprintf("Trade completed between %s and %s", from.Name, to.Name))
}

// aiEvaluateTrade decides if the AI should accept a trade offer.
func (g *Game) aiEvaluateTrade(offer TradeOffer) bool {
	aiID := offer.ToPlayer
	received := offer.OfferedMoney
	given := offer.WantedMoney

	for _, idx := range offer.OfferedProps {
		space := g.Board.Spaces[idx]
		value := space.Price
		// Weight higher if receiving this property would complete AI's monopoly
		if g.almostMonopoly(aiID, space.Group) {
			value = value * 18 / 10 // 1.8x
		}
		received += value
	}
	for _, idx := range offer.WantedProps {
		space := g.Board.Spaces[idx]
		value := space.Price

		// Reject if this would complete opponent's monopoly
		if g.wouldCompleteMonopoly(offer.FromPlayer, space.Group) {
			return false
		}
		// Weight higher if AI almost has a monopoly in that group (reluctant to give up)
		if g.almostMonopoly(aiID, space.Group) {
			value = value * 15 / 10 // 1.5x
		}
		given += value
	}

	return received >= given
}

// almostMonopoly returns true if the player owns all but one property in a group.
func (g *Game) almostMonopoly(playerID int, group board.ColorGroup) bool {
	if group == board.GroupNone {
		return false
	}
	spaces := g.Board.SpacesInGroup(group)
	if len(spaces) == 0 {
		return false
	}
	owned := 0
	for _, idx := range spaces {
		if g.Board.Properties[idx].OwnerID == playerID {
			owned++
		}
	}
	return owned == len(spaces)-1
}

// wouldCompleteMonopoly checks if giving a player a property would complete their monopoly.
func (g *Game) wouldCompleteMonopoly(playerID int, group board.ColorGroup) bool {
	if group == board.GroupNone {
		return false
	}
	spaces := g.Board.SpacesInGroup(group)
	owned := 0
	for _, idx := range spaces {
		if g.Board.Properties[idx].OwnerID == playerID {
			owned++
		}
	}
	return owned >= len(spaces)-1
}
