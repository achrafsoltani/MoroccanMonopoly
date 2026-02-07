package board

import "github.com/AchrafSoltani/MoroccanMonopoly/config"

// Board holds all spaces, property states, card decks, and house/hotel pools.
type Board struct {
	Spaces         [config.SpaceCount]Space
	Properties     [config.SpaceCount]PropertyState
	HousePool      int
	HotelPool      int
	ChanceDeck     *Deck
	CommunityDeck  *Deck
}

// NewBoard creates and initialises the full 40-space Moroccan Monopoly board.
func NewBoard() *Board {
	b := &Board{
		HousePool:     config.MaxHouses,
		HotelPool:     config.MaxHotels,
		ChanceDeck:    NewDeck(ChanceCards()),
		CommunityDeck: NewDeck(CommunityChestCards()),
	}

	for i := 0; i < config.SpaceCount; i++ {
		b.Properties[i] = NewPropertyState()
	}

	// Bottom row (spaces 0-10)
	b.Spaces[0] = Space{Index: 0, Name: "DEPART", Type: SpaceGo}
	b.Spaces[1] = Space{Index: 1, Name: "Derb Sultan", ShortName: "Derb S.", Type: SpaceProperty, Group: GroupBrown,
		Price: 60, Rent: [6]int{2, 10, 30, 90, 160, 250}, HouseCost: 50}
	b.Spaces[2] = Space{Index: 2, Name: "Caisse Commune", Type: SpaceCommunityChest}
	b.Spaces[3] = Space{Index: 3, Name: "Bab Marrakech", ShortName: "Bab Mk.", Type: SpaceProperty, Group: GroupBrown,
		Price: 60, Rent: [6]int{4, 20, 60, 180, 320, 450}, HouseCost: 50}
	b.Spaces[4] = Space{Index: 4, Name: "Impot sur le Revenu", ShortName: "IMPOT", Type: SpaceTax, TaxAmount: config.IncomeTax}
	b.Spaces[5] = Space{Index: 5, Name: "Gare Casa-Voyageurs", ShortName: "G.Casa", Type: SpaceRailroad, Price: 200}
	b.Spaces[6] = Space{Index: 6, Name: "Av. Hassan II", ShortName: "HassnII", Type: SpaceProperty, Group: GroupLightBlue,
		Price: 100, Rent: [6]int{6, 30, 90, 270, 400, 550}, HouseCost: 50}
	b.Spaces[7] = Space{Index: 7, Name: "Chance", Type: SpaceChance}
	b.Spaces[8] = Space{Index: 8, Name: "Bab Bou Jeloud", ShortName: "B.Jelud", Type: SpaceProperty, Group: GroupLightBlue,
		Price: 100, Rent: [6]int{6, 30, 90, 270, 400, 550}, HouseCost: 50}
	b.Spaces[9] = Space{Index: 9, Name: "Talaa Kebira", ShortName: "Talaa K", Type: SpaceProperty, Group: GroupLightBlue,
		Price: 120, Rent: [6]int{8, 40, 100, 300, 450, 600}, HouseCost: 50}
	b.Spaces[10] = Space{Index: 10, Name: "EN VISITE", Type: SpaceJail}

	// Left column (spaces 11-20)
	b.Spaces[11] = Space{Index: 11, Name: "Av. Mohammed V", ShortName: "Mohd V", Type: SpaceProperty, Group: GroupPink,
		Price: 140, Rent: [6]int{10, 50, 150, 450, 625, 750}, HouseCost: 100}
	b.Spaces[12] = Space{Index: 12, Name: "ONEE (Electricite)", ShortName: "ONEE", Type: SpaceUtility, Price: 150}
	b.Spaces[13] = Space{Index: 13, Name: "Rue des Consuls", ShortName: "Consuls", Type: SpaceProperty, Group: GroupPink,
		Price: 140, Rent: [6]int{10, 50, 150, 450, 625, 750}, HouseCost: 100}
	b.Spaces[14] = Space{Index: 14, Name: "Kasbah Oudayas", ShortName: "Kasbah", Type: SpaceProperty, Group: GroupPink,
		Price: 160, Rent: [6]int{12, 60, 180, 500, 700, 900}, HouseCost: 100}
	b.Spaces[15] = Space{Index: 15, Name: "Gare Rabat-Ville", ShortName: "G.Rabat", Type: SpaceRailroad, Price: 200}
	b.Spaces[16] = Space{Index: 16, Name: "Av. de la Liberte", ShortName: "Av Lbrt", Type: SpaceProperty, Group: GroupOrange,
		Price: 180, Rent: [6]int{14, 70, 200, 550, 750, 950}, HouseCost: 100}
	b.Spaces[17] = Space{Index: 17, Name: "Caisse Commune", Type: SpaceCommunityChest}
	b.Spaces[18] = Space{Index: 18, Name: "Rue de la Liberte", ShortName: "Ru Lbrt", Type: SpaceProperty, Group: GroupOrange,
		Price: 180, Rent: [6]int{14, 70, 200, 550, 750, 950}, HouseCost: 100}
	b.Spaces[19] = Space{Index: 19, Name: "Grand Socco", ShortName: "Socco", Type: SpaceProperty, Group: GroupOrange,
		Price: 200, Rent: [6]int{16, 80, 220, 600, 800, 1000}, HouseCost: 100}
	b.Spaces[20] = Space{Index: 20, Name: "PARKING GRATUIT", Type: SpaceFreeParking}

	// Top row (spaces 21-30)
	b.Spaces[21] = Space{Index: 21, Name: "Jemaa el-Fna", ShortName: "Jemaa", Type: SpaceProperty, Group: GroupRed,
		Price: 220, Rent: [6]int{18, 90, 250, 700, 875, 1050}, HouseCost: 150}
	b.Spaces[22] = Space{Index: 22, Name: "Chance", Type: SpaceChance}
	b.Spaces[23] = Space{Index: 23, Name: "Rue Bab Agnaou", ShortName: "Agnaou", Type: SpaceProperty, Group: GroupRed,
		Price: 220, Rent: [6]int{18, 90, 250, 700, 875, 1050}, HouseCost: 150}
	b.Spaces[24] = Space{Index: 24, Name: "Koutoubia", ShortName: "Koutbia", Type: SpaceProperty, Group: GroupRed,
		Price: 240, Rent: [6]int{20, 100, 300, 750, 925, 1100}, HouseCost: 150}
	b.Spaces[25] = Space{Index: 25, Name: "Gare Marrakech", ShortName: "G.Mrkch", Type: SpaceRailroad, Price: 200}
	b.Spaces[26] = Space{Index: 26, Name: "Av. Mohammed VI", ShortName: "Mohd VI", Type: SpaceProperty, Group: GroupYellow,
		Price: 260, Rent: [6]int{22, 110, 330, 800, 975, 1150}, HouseCost: 150}
	b.Spaces[27] = Space{Index: 27, Name: "Corniche Ain Diab", ShortName: "AinDiab", Type: SpaceProperty, Group: GroupYellow,
		Price: 260, Rent: [6]int{22, 110, 330, 800, 975, 1150}, HouseCost: 150}
	b.Spaces[28] = Space{Index: 28, Name: "LYDEC (Eau)", ShortName: "LYDEC", Type: SpaceUtility, Price: 150}
	b.Spaces[29] = Space{Index: 29, Name: "Bd. de la Corniche", ShortName: "Cornich", Type: SpaceProperty, Group: GroupYellow,
		Price: 280, Rent: [6]int{24, 120, 360, 850, 1025, 1200}, HouseCost: 150}
	b.Spaces[30] = Space{Index: 30, Name: "ALLEZ EN PRISON", Type: SpaceGoToJail}

	// Right column (spaces 31-39)
	b.Spaces[31] = Space{Index: 31, Name: "Vallee du Dades", ShortName: "Dades", Type: SpaceProperty, Group: GroupGreen,
		Price: 300, Rent: [6]int{26, 130, 390, 900, 1100, 1275}, HouseCost: 200}
	b.Spaces[32] = Space{Index: 32, Name: "Gorges du Todra", ShortName: "Todra", Type: SpaceProperty, Group: GroupGreen,
		Price: 300, Rent: [6]int{26, 130, 390, 900, 1100, 1275}, HouseCost: 200}
	b.Spaces[33] = Space{Index: 33, Name: "Caisse Commune", Type: SpaceCommunityChest}
	b.Spaces[34] = Space{Index: 34, Name: "Merzouga (Sahara)", ShortName: "Merzoug", Type: SpaceProperty, Group: GroupGreen,
		Price: 320, Rent: [6]int{28, 150, 450, 1000, 1200, 1400}, HouseCost: 200}
	b.Spaces[35] = Space{Index: 35, Name: "Gare Tanger-Ville", ShortName: "G.Tangr", Type: SpaceRailroad, Price: 200}
	b.Spaces[36] = Space{Index: 36, Name: "Chance", Type: SpaceChance}
	b.Spaces[37] = Space{Index: 37, Name: "Chefchaouen", ShortName: "Chefch.", Type: SpaceProperty, Group: GroupDarkBlue,
		Price: 350, Rent: [6]int{35, 175, 500, 1100, 1300, 1500}, HouseCost: 200}
	b.Spaces[38] = Space{Index: 38, Name: "Taxe de Luxe", ShortName: "T.LUXE", Type: SpaceTax, TaxAmount: config.LuxuryTax}
	b.Spaces[39] = Space{Index: 39, Name: "Mosquee Hassan II", ShortName: "Msq.HII", Type: SpaceProperty, Group: GroupDarkBlue,
		Price: 400, Rent: [6]int{50, 200, 600, 1400, 1700, 2000}, HouseCost: 200}

	return b
}

// IsProperty returns true if the space can be owned (property, railroad, or utility).
func (b *Board) IsProperty(index int) bool {
	t := b.Spaces[index].Type
	return t == SpaceProperty || t == SpaceRailroad || t == SpaceUtility
}

// SpacesInGroup returns all space indices belonging to a colour group.
func (b *Board) SpacesInGroup(g ColorGroup) []int {
	var result []int
	for i, s := range b.Spaces {
		if s.Group == g {
			result = append(result, i)
		}
	}
	return result
}

// RailroadSpaces returns all railroad space indices.
func (b *Board) RailroadSpaces() []int {
	var result []int
	for i, s := range b.Spaces {
		if s.Type == SpaceRailroad {
			result = append(result, i)
		}
	}
	return result
}

// UtilitySpaces returns all utility space indices.
func (b *Board) UtilitySpaces() []int {
	var result []int
	for i, s := range b.Spaces {
		if s.Type == SpaceUtility {
			result = append(result, i)
		}
	}
	return result
}
