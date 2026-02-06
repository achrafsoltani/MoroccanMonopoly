package player

import "github.com/AchrafSoltani/MoroccanMonopoly/config"

// Player represents a game participant.
type Player struct {
	ID       int
	Name     string
	IsAI     bool
	Money    int
	Position int
	InJail   bool
	JailTurns int
	Bankrupt  bool
	Properties []int // space indices owned
	GetOutOfJailCards int
}

// NewPlayer creates a player with starting money.
func NewPlayer(id int, name string, isAI bool) *Player {
	return &Player{
		ID:    id,
		Name:  name,
		IsAI:  isAI,
		Money: config.StartingMoney,
	}
}

// AddProperty records ownership of a property.
func (p *Player) AddProperty(spaceIndex int) {
	p.Properties = append(p.Properties, spaceIndex)
}

// RemoveProperty removes a property from ownership.
func (p *Player) RemoveProperty(spaceIndex int) {
	for i, idx := range p.Properties {
		if idx == spaceIndex {
			p.Properties = append(p.Properties[:i], p.Properties[i+1:]...)
			return
		}
	}
}

// OwnsProperty checks if the player owns a specific space.
func (p *Player) OwnsProperty(spaceIndex int) bool {
	for _, idx := range p.Properties {
		if idx == spaceIndex {
			return true
		}
	}
	return false
}

// Pay deducts money from the player. Returns false if insufficient funds.
func (p *Player) Pay(amount int) bool {
	if p.Money >= amount {
		p.Money -= amount
		return true
	}
	return false
}

// Receive adds money to the player.
func (p *Player) Receive(amount int) {
	p.Money += amount
}

// NetWorth calculates total assets (money + property values + house values).
// Requires board access, so takes property prices and house costs as params.
func (p *Player) TotalProperties() int {
	return len(p.Properties)
}
