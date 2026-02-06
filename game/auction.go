package game

import (
	"fmt"

	"github.com/AchrafSoltani/MoroccanMonopoly/config"
)

// startAuction begins an auction for the given property.
func (g *Game) startAuction(spaceIndex int) {
	g.AuctionSpaceIdx = spaceIndex
	g.AuctionHighBid = 0
	g.AuctionHighBidder = -1
	g.AuctionCurrent = g.Current

	// Mark all non-bankrupt players as active
	for i, p := range g.Players {
		g.AuctionActive[i] = !p.Bankrupt
	}

	g.Phase = PhaseAuction
	g.Dialog = DialogAuction
	space := g.Board.Spaces[spaceIndex]
	g.AddMessage(fmt.Sprintf("Auction started for %s!", space.Name))
	g.advanceAuction()
}

// advanceAuction moves to the next active bidder.
func (g *Game) advanceAuction() {
	// Find next active player
	for {
		g.AuctionCurrent = (g.AuctionCurrent + 1) % len(g.Players)
		if g.AuctionActive[g.AuctionCurrent] {
			break
		}
		// Check if only one active remains
		activeCount := 0
		lastActive := -1
		for i := 0; i < len(g.Players); i++ {
			if g.AuctionActive[i] {
				activeCount++
				lastActive = i
			}
		}
		if activeCount <= 1 {
			g.endAuction(lastActive)
			return
		}
	}

	// Check if auction should end (everyone passed except high bidder)
	activeCount := 0
	for i := 0; i < len(g.Players); i++ {
		if g.AuctionActive[i] {
			activeCount++
		}
	}
	if activeCount <= 1 {
		g.endAuction(g.AuctionHighBidder)
		return
	}

	// If current player is AI, auto-bid
	if g.Players[g.AuctionCurrent].IsAI {
		g.aiBid()
	}
}

// handleAuctionClick processes auction dialog button clicks.
func (g *Game) handleAuctionClick() {
	if g.DialogHovered < 0 {
		return
	}

	p := g.Players[g.AuctionCurrent]
	bidAmount := g.AuctionHighBid + 10

	switch g.DialogHovered {
	case 0: // Bid
		if p.Money >= bidAmount {
			g.AuctionHighBid = bidAmount
			g.AuctionHighBidder = g.AuctionCurrent
			g.AddMessage(fmt.Sprintf("%s bids %d MAD", p.Name, bidAmount))
		}
		g.advanceAuction()
	case 1: // Pass
		g.AuctionActive[g.AuctionCurrent] = false
		g.AddMessage(fmt.Sprintf("%s passes", p.Name))
		g.advanceAuction()
	}
}

// aiBid handles AI auction bidding.
func (g *Game) aiBid() {
	p := g.Players[g.AuctionCurrent]
	space := g.Board.Spaces[g.AuctionSpaceIdx]
	bidAmount := g.AuctionHighBid + 10

	// AI bids up to 80% of property value if it has enough money
	maxBid := space.Price * 80 / 100
	if bidAmount <= maxBid && p.Money >= bidAmount+100 {
		g.AuctionHighBid = bidAmount
		g.AuctionHighBidder = g.AuctionCurrent
		g.AddMessage(fmt.Sprintf("%s (AI) bids %d MAD", p.Name, bidAmount))
	} else {
		g.AuctionActive[g.AuctionCurrent] = false
		g.AddMessage(fmt.Sprintf("%s (AI) passes", p.Name))
	}
	g.advanceAuction()
}

// endAuction finishes the auction, transferring property to the winner.
func (g *Game) endAuction(winnerIdx int) {
	space := g.Board.Spaces[g.AuctionSpaceIdx]

	if winnerIdx < 0 || g.AuctionHighBid <= 0 {
		g.AddMessage(fmt.Sprintf("No bids! %s remains unowned.", space.Name))
	} else {
		winner := g.Players[winnerIdx]
		winner.Pay(g.AuctionHighBid)
		winner.AddProperty(g.AuctionSpaceIdx)
		g.Board.Properties[g.AuctionSpaceIdx].OwnerID = winnerIdx
		g.AddMessage(fmt.Sprintf("%s wins auction for %s at %d MAD!", winner.Name, space.Name, g.AuctionHighBid))
	}

	// Reset auction state
	for i := range g.AuctionActive {
		g.AuctionActive[i] = false
		g.AuctionBids[i] = 0
	}

	g.Dialog = DialogNone
	g.Phase = PhasePostAction
	g.updateButtonStates()
}

// Ensure config is used
var _ = config.MaxPlayers
