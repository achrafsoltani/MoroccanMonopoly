package render

import "github.com/AchrafSoltani/glow"

// Moroccan colour palette
var (
	// Board background and borders
	BoardBg     = glow.Color{R: 205, G: 230, B: 208} // soft mint green
	BoardBorder = glow.Color{R: 0, G: 80, B: 0}      // dark green border
	SpaceBg     = glow.Color{R: 215, G: 235, B: 215}  // light green space bg

	// Property colour group strips
	ColorBrown    = glow.Color{R: 139, G: 90, B: 43}
	ColorLightBlue = glow.Color{R: 135, G: 206, B: 235}
	ColorPink     = glow.Color{R: 219, G: 112, B: 147}
	ColorOrange   = glow.Color{R: 237, G: 145, B: 33}
	ColorRed      = glow.Color{R: 205, G: 43, B: 30}
	ColorYellow   = glow.Color{R: 252, G: 209, B: 22}
	ColorGreen    = glow.Color{R: 30, G: 130, B: 76}
	ColorDarkBlue = glow.Color{R: 0, G: 56, B: 168}

	// Special spaces
	ColorGo        = glow.Color{R: 255, G: 240, B: 220}
	ColorJail      = glow.Color{R: 240, G: 200, B: 160}
	ColorParking   = glow.Color{R: 255, G: 240, B: 220}
	ColorGoToJail  = glow.Color{R: 240, G: 200, B: 160}
	ColorChance    = glow.Color{R: 255, G: 165, B: 0}
	ColorCommunity = glow.Color{R: 30, G: 144, B: 255}
	ColorTax       = glow.Color{R: 200, G: 200, B: 200}
	ColorRailroad  = glow.Color{R: 200, G: 200, B: 200}
	ColorUtility   = glow.Color{R: 200, G: 200, B: 200}

	// UI colours
	TextDark    = glow.Color{R: 20, G: 20, B: 20}
	TextLight   = glow.Color{R: 240, G: 240, B: 240}
	TextGold    = glow.Color{R: 218, G: 165, B: 32}
	PanelBg     = glow.Color{R: 35, G: 60, B: 35}
	PanelBorder = glow.Color{R: 100, G: 140, B: 100}
	ButtonBg    = glow.Color{R: 60, G: 100, B: 60}
	ButtonHover = glow.Color{R: 80, G: 130, B: 80}
	ButtonText  = glow.Color{R: 240, G: 240, B: 240}
	DialogBg    = glow.Color{R: 40, G: 40, B: 40}
	DialogBorder = glow.Color{R: 180, G: 140, B: 60}

	// Player token colours
	PlayerColors = [4]glow.Color{
		{R: 220, G: 50, B: 50},   // red
		{R: 50, G: 120, B: 220},  // blue
		{R: 50, G: 180, B: 50},   // green
		{R: 220, G: 180, B: 50},  // gold
	}

	// Moroccan decorative colours
	ZelligeGreen = glow.Color{R: 0, G: 120, B: 60}
	ZelligeBlue  = glow.Color{R: 20, G: 70, B: 140}
	ZelligeGold  = glow.Color{R: 200, G: 160, B: 40}

	// House/hotel indicators
	HouseColor = glow.Color{R: 0, G: 150, B: 0}
	HotelColor = glow.Color{R: 200, G: 0, B: 0}

	// Mortgage indicator
	MortgageColor = glow.Color{R: 150, G: 150, B: 150}
)
