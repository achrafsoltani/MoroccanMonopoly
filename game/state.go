package game

// GameState represents the top-level game state.
type GameState int

const (
	StateMenu    GameState = iota
	StateSetup             // player setup screen
	StatePlaying           // main gameplay
	StateGameOver
)

// TurnPhase represents the phase within a player's turn.
type TurnPhase int

const (
	PhasePreRoll      TurnPhase = iota
	PhaseJailDecision           // choosing how to exit jail
	PhaseRolling                // dice animation
	PhaseMoving                 // token moving animation
	PhaseLanded                 // just landed, evaluate space
	PhaseDialog                 // showing buy/rent/card/tax dialog
	PhaseAuction                // auction in progress
	PhasePostAction             // post-landing, may roll again (doubles)
	PhaseTurnEnd                // turn ending, advance to next player
	PhaseTrade                  // trade dialog open
	PhaseBuild                  // build dialog open
	PhaseMortgage               // mortgage dialog open
)

// DialogType identifies which dialog is showing.
type DialogType int

const (
	DialogNone         DialogType = iota
	DialogBuyProperty
	DialogPayRent
	DialogChanceCard
	DialogCommunityCard
	DialogPayTax
	DialogJailOptions
	DialogBuild
	DialogMortgage
	DialogTrade
	DialogAuction
	DialogBankruptcy
	DialogGameOver
)
