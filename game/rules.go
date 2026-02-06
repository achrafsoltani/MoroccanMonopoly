package game

import (
	"github.com/AchrafSoltani/MoroccanMonopoly/board"
	"github.com/AchrafSoltani/MoroccanMonopoly/config"
)

// CanBuildOnGroup checks if a player can build on a colour group.
// Requires: owns all properties in group, none mortgaged, even building rule.
func (g *Game) CanBuildOnGroup(playerID int, group board.ColorGroup) bool {
	if group == board.GroupNone {
		return false
	}
	if !g.hasMonopoly(playerID, group) {
		return false
	}
	spaces := g.Board.SpacesInGroup(group)
	for _, idx := range spaces {
		if g.Board.Properties[idx].Mortgaged {
			return false
		}
	}
	return true
}

// CanBuildOnSpace checks if a house/hotel can be built on a specific property.
// Even building rule: difference between min and max houses in group <= 1.
func (g *Game) CanBuildOnSpace(spaceIndex int) bool {
	space := g.Board.Spaces[spaceIndex]
	prop := g.Board.Properties[spaceIndex]

	if space.Type != board.SpaceProperty {
		return false
	}
	if prop.OwnerID < 0 || prop.Mortgaged {
		return false
	}
	if !g.CanBuildOnGroup(prop.OwnerID, space.Group) {
		return false
	}
	if prop.Houses >= config.HotelLevel {
		return false
	}

	// Check pool availability
	if prop.Houses == config.HousesPerHotel {
		// Upgrading to hotel
		if g.Board.HotelPool <= 0 {
			return false
		}
	} else {
		if g.Board.HousePool <= 0 {
			return false
		}
	}

	// Even building rule: this property must have the fewest houses in its group
	minHouses := g.minHousesInGroup(space.Group)
	return prop.Houses <= minHouses
}

// BuildHouse adds a house (or hotel) to a property.
func (g *Game) BuildHouse(spaceIndex int) int {
	space := g.Board.Spaces[spaceIndex]
	prop := &g.Board.Properties[spaceIndex]

	if prop.Houses == config.HousesPerHotel {
		// Upgrade to hotel: return 4 houses to pool, take 1 hotel
		g.Board.HousePool += config.HousesPerHotel
		g.Board.HotelPool--
		prop.Houses = config.HotelLevel
	} else {
		g.Board.HousePool--
		prop.Houses++
	}

	return space.HouseCost
}

// SellHouse removes a house from a property (even sell-down rule).
func (g *Game) CanSellHouseOnSpace(spaceIndex int) bool {
	space := g.Board.Spaces[spaceIndex]
	prop := g.Board.Properties[spaceIndex]

	if space.Type != board.SpaceProperty {
		return false
	}
	if prop.Houses <= 0 {
		return false
	}

	// Even selling rule: this property must have the most houses in its group
	maxHouses := g.maxHousesInGroup(space.Group)
	return prop.Houses >= maxHouses
}

// SellHouse removes a house and returns half the house cost.
func (g *Game) SellHouse(spaceIndex int) int {
	space := g.Board.Spaces[spaceIndex]
	prop := &g.Board.Properties[spaceIndex]

	if prop.Houses == config.HotelLevel {
		// Downgrade from hotel: return 1 hotel, take 4 houses
		g.Board.HotelPool++
		g.Board.HousePool -= config.HousesPerHotel
		prop.Houses = config.HousesPerHotel
	} else {
		g.Board.HousePool++
		prop.Houses--
	}

	return space.HouseCost / 2
}

// BuildableProperties returns all space indices where the player can build.
func (g *Game) BuildableProperties(playerID int) []int {
	var result []int
	for _, idx := range g.Players[playerID].Properties {
		if g.CanBuildOnSpace(idx) {
			result = append(result, idx)
		}
	}
	return result
}

// MortgageableProperties returns properties that can be mortgaged.
func (g *Game) MortgageableProperties(playerID int) []int {
	var result []int
	for _, idx := range g.Players[playerID].Properties {
		prop := g.Board.Properties[idx]
		if !prop.Mortgaged && prop.Houses == 0 {
			result = append(result, idx)
		}
	}
	return result
}

// UnmortgageableProperties returns mortgaged properties that can be unmortgaged.
func (g *Game) UnmortgageableProperties(playerID int) []int {
	var result []int
	for _, idx := range g.Players[playerID].Properties {
		prop := g.Board.Properties[idx]
		if prop.Mortgaged {
			cost := g.UnmortgageCost(idx)
			if g.Players[playerID].Money >= cost {
				result = append(result, idx)
			}
		}
	}
	return result
}

// MortgageValue returns cash received for mortgaging.
func (g *Game) MortgageValue(spaceIndex int) int {
	return g.Board.Spaces[spaceIndex].Price * config.MortgageRate / 100
}

// UnmortgageCost returns cost to unmortgage.
func (g *Game) UnmortgageCost(spaceIndex int) int {
	mortgageVal := g.MortgageValue(spaceIndex)
	return mortgageVal * config.UnmortgageRate / 100
}

// MortgageProperty mortgages a property.
func (g *Game) MortgageProperty(spaceIndex int) int {
	g.Board.Properties[spaceIndex].Mortgaged = true
	return g.MortgageValue(spaceIndex)
}

// UnmortgageProperty unmortgages a property.
func (g *Game) UnmortgageProperty(spaceIndex int) int {
	g.Board.Properties[spaceIndex].Mortgaged = false
	return g.UnmortgageCost(spaceIndex)
}

func (g *Game) minHousesInGroup(group board.ColorGroup) int {
	spaces := g.Board.SpacesInGroup(group)
	min := 999
	for _, idx := range spaces {
		h := g.Board.Properties[idx].Houses
		if h < min {
			min = h
		}
	}
	return min
}

func (g *Game) maxHousesInGroup(group board.ColorGroup) int {
	spaces := g.Board.SpacesInGroup(group)
	max := -1
	for _, idx := range spaces {
		h := g.Board.Properties[idx].Houses
		if h > max {
			max = h
		}
	}
	return max
}

// PlayerNetWorth calculates a player's total asset value.
func (g *Game) PlayerNetWorth(playerID int) int {
	p := g.Players[playerID]
	total := p.Money
	for _, idx := range p.Properties {
		space := g.Board.Spaces[idx]
		prop := g.Board.Properties[idx]
		if prop.Mortgaged {
			total += g.MortgageValue(idx)
		} else {
			total += space.Price
		}
		if prop.Houses > 0 && prop.Houses <= config.HousesPerHotel {
			total += prop.Houses * space.HouseCost / 2
		} else if prop.Houses == config.HotelLevel {
			total += config.HousesPerHotel * space.HouseCost / 2
		}
	}
	return total
}
