package board

// PropertyState tracks the ownership and development state of a property.
type PropertyState struct {
	OwnerID   int  // -1 = unowned
	Houses    int  // 0-4 houses, 5 = hotel
	Mortgaged bool
}

// NewPropertyState returns a fresh unowned property state.
func NewPropertyState() PropertyState {
	return PropertyState{OwnerID: -1}
}
