package board

// SpaceType identifies what kind of board space this is.
type SpaceType int

const (
	SpaceGo             SpaceType = iota
	SpaceProperty                 // colour-group property
	SpaceCommunityChest           // Caisse Commune
	SpaceChance                   // Chance
	SpaceTax                      // Income tax or luxury tax
	SpaceRailroad                 // Gare
	SpaceUtility                  // ONEE, LYDEC
	SpaceJail                     // En Visite / jail
	SpaceFreeParking              // Parking Gratuit
	SpaceGoToJail                 // Allez en Prison
)

// Space represents a single board space.
type Space struct {
	Index     int
	Name      string
	ShortName string // short display name for the board (max ~7 chars)
	Type      SpaceType
	Group     ColorGroup
	Price     int
	Rent      [6]int // base, 1 house, 2 houses, 3 houses, 4 houses, hotel
	HouseCost int
	TaxAmount int // only for SpaceTax
}
