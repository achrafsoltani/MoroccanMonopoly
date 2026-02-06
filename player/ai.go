package player

// AI decision weights and thresholds.
const (
	AIBuyBufferEarly = 200 // keep at least this much after buying (early game)
	AIBuyBufferLate  = 100 // keep at least this much after buying (late game)
	AIBuildBuffer    = 150 // keep after building
)

// ShouldBuy decides if an AI player should buy a property.
func (p *Player) ShouldBuy(price, totalProps int) bool {
	buffer := AIBuyBufferEarly
	if totalProps > 10 { // late game
		buffer = AIBuyBufferLate
	}
	return p.Money >= price+buffer
}

// ShouldBuild decides if an AI player should build a house.
func (p *Player) ShouldBuild(houseCost int) bool {
	return p.Money >= houseCost+AIBuildBuffer
}
