package board

// ColorGroup represents property colour groups.
type ColorGroup int

const (
	GroupNone     ColorGroup = iota
	GroupBrown               // Derb Sultan, Bab Marrakech
	GroupLightBlue           // Av. Hassan II, Bab Bou Jeloud, Talaa Kebira
	GroupPink                // Av. Mohammed V, Rue des Consuls, Kasbah des Oudayas
	GroupOrange              // Av. de la Liberte, Rue de la Liberte, Grand Socco
	GroupRed                 // Jemaa el-Fna, Rue Bab Agnaou, Koutoubia
	GroupYellow              // Av. Mohammed VI, Corniche Ain Diab, Bd. de la Corniche
	GroupGreen               // Vallee du Dades, Gorges du Todra, Merzouga
	GroupDarkBlue            // Chefchaouen, Mosquee Hassan II
)

// GroupSize returns the number of properties in a colour group.
func GroupSize(g ColorGroup) int {
	switch g {
	case GroupBrown, GroupDarkBlue:
		return 2
	case GroupLightBlue, GroupPink, GroupOrange, GroupRed, GroupYellow, GroupGreen:
		return 3
	default:
		return 0
	}
}
