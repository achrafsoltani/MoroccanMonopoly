package config

// Window dimensions
const (
	WindowWidth  = 1200
	WindowHeight = 800
	WindowTitle  = "Monopoly Maroc"
)

// Board geometry
const (
	BoardX      = 20  // board top-left X
	BoardY      = 50  // board top-left Y
	BoardSize   = 698 // total board size (square)
	CornerSize  = 88  // corner space width/height
	SpaceWidth  = 58  // non-corner space width along edge
	SpaceDepth  = 88  // non-corner space depth
	SpaceCount  = 40  // total board spaces
	SidesSpaces = 9   // non-corner spaces per side
)

// Info panel
const (
	PanelX     = 740
	PanelY     = 0
	PanelWidth = 440
)

// Game constants
const (
	StartingMoney  = 1500
	GoSalary       = 200
	MaxPlayers     = 4
	MinPlayers     = 2
	MaxHouses      = 32
	MaxHotels      = 12
	HousesPerHotel = 4
	HotelLevel     = 5 // 4 houses + 1 hotel
	JailPosition   = 10
	GoToJailPos    = 30
	GoPosition     = 0
	JailFine       = 50
	MaxJailTurns   = 3
	IncomeTax      = 200
	LuxuryTax      = 100
	MortgageRate   = 50  // percent of price
	UnmortgageRate = 110 // percent of mortgage value
)

// Layout holds geometry values computed from window dimensions.
type Layout struct {
	WinW, WinH                         int
	BoardX, BoardY, BoardSize          int
	CornerSize, SpaceWidth, SpaceDepth int
	PanelX, PanelWidth                 int
}

// NewLayout computes a responsive layout from window dimensions.
// The board stays square, sized to fit the available height (with margins).
// The info panel fills the remaining width on the right.
func NewLayout(winW, winH int) Layout {
	margin := 20
	topMargin := 50

	// Available height for the board
	availH := winH - topMargin - margin
	if availH < 400 {
		availH = 400
	}

	// Board is square, must also leave room for the panel (min 200px)
	maxBoardFromWidth := winW - margin - 200 - margin
	boardSize := availH
	if maxBoardFromWidth < boardSize {
		boardSize = maxBoardFromWidth
	}
	if boardSize < 400 {
		boardSize = 400
	}

	// Scale corner and space widths from original ratios:
	// Original: BoardSize=698, CornerSize=88, SpaceWidth=58
	// cornerSize = boardSize * 88/698, spaceWidth = (boardSize - 2*cornerSize) / 9
	cornerSize := boardSize * CornerSize / BoardSize
	spaceWidth := (boardSize - 2*cornerSize) / SidesSpaces
	// Recompute board size to eliminate rounding gaps
	boardSize = 2*cornerSize + SidesSpaces*spaceWidth

	boardX := margin
	boardY := topMargin

	panelX := boardX + boardSize + margin
	panelWidth := winW - panelX
	if panelWidth < 200 {
		panelWidth = 200
	}

	return Layout{
		WinW:       winW,
		WinH:       winH,
		BoardX:     boardX,
		BoardY:     boardY,
		BoardSize:  boardSize,
		CornerSize: cornerSize,
		SpaceWidth: spaceWidth,
		SpaceDepth: cornerSize,
		PanelX:     panelX,
		PanelWidth: panelWidth,
	}
}

// Timing
const (
	DiceAnimDuration  = 0.8  // seconds
	TokenMoveDuration = 0.15 // seconds per space
	AITurnDelay       = 1.0  // seconds between AI actions
	MessageDuration   = 3.0  // seconds to show messages
)
