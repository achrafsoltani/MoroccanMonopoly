package board

import "math/rand"

// CardEffect types
type CardEffectType int

const (
	EffectCollect      CardEffectType = iota // Collect money from bank
	EffectPay                               // Pay money to bank
	EffectMoveTo                            // Move to specific space
	EffectMoveSteps                         // Move forward N steps
	EffectGoToJail                          // Go directly to jail
	EffectGetOutOfJail                      // Get out of jail free card
	EffectPayPerHouse                       // Pay per house/hotel owned
	EffectCollectAll                        // Collect from each player
	EffectPayAll                            // Pay each player
	EffectMoveNearest                       // Move to nearest railroad/utility
)

// Card represents a Chance or Community Chest card.
type Card struct {
	Text       string
	Effect     CardEffectType
	Amount     int // money amount or target space index
	AmountHotel int // per-hotel amount (for EffectPayPerHouse)
}

// Deck represents a shuffled card deck.
type Deck struct {
	Cards   []Card
	Current int
}

// NewDeck creates a deck from the given cards and shuffles it.
func NewDeck(cards []Card) *Deck {
	d := &Deck{
		Cards: make([]Card, len(cards)),
	}
	copy(d.Cards, cards)
	d.Shuffle()
	return d
}

// Shuffle randomises the deck order.
func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
	d.Current = 0
}

// Draw returns the next card from the deck.
func (d *Deck) Draw() Card {
	card := d.Cards[d.Current]
	d.Current++
	if d.Current >= len(d.Cards) {
		d.Shuffle()
	}
	return card
}

// ChanceCards returns the Moroccan-themed Chance deck.
func ChanceCards() []Card {
	return []Card{
		{Text: "Avancez jusqu'a DEPART. Recevez 200 MAD.", Effect: EffectMoveTo, Amount: 0},
		{Text: "Allez a Jemaa el-Fna. Si vous passez par DEPART, recevez 200 MAD.", Effect: EffectMoveTo, Amount: 21},
		{Text: "Allez a Av. Mohammed V (Rabat). Si vous passez par DEPART, recevez 200 MAD.", Effect: EffectMoveTo, Amount: 11},
		{Text: "Allez a Gare Marrakech. Si vous passez par DEPART, recevez 200 MAD.", Effect: EffectMoveTo, Amount: 25},
		{Text: "La banque vous verse 50 MAD de dividendes.", Effect: EffectCollect, Amount: 50},
		{Text: "Vous avez gagne le prix du Festival de Fes! Recevez 150 MAD.", Effect: EffectCollect, Amount: 150},
		{Text: "Carte de sortie de prison gratuite.", Effect: EffectGetOutOfJail},
		{Text: "Reculez de 3 cases.", Effect: EffectMoveSteps, Amount: -3},
		{Text: "Allez en prison. Ne passez pas par DEPART.", Effect: EffectGoToJail},
		{Text: "Faites des reparations: payez 25 MAD par maison et 100 MAD par hotel.", Effect: EffectPayPerHouse, Amount: 25, AmountHotel: 100},
		{Text: "Amende pour exces de vitesse: 15 MAD.", Effect: EffectPay, Amount: 15},
		{Text: "Voyage au souk! Allez a Gare Casa-Voyageurs.", Effect: EffectMoveTo, Amount: 5},
		{Text: "Elu president du conseil communal. Payez 50 MAD a chaque joueur.", Effect: EffectPayAll, Amount: 50},
		{Text: "Votre investissement immobilier rapporte: recevez 100 MAD.", Effect: EffectCollect, Amount: 100},
		{Text: "Allez a Chefchaouen. Si vous passez par DEPART, recevez 200 MAD.", Effect: EffectMoveTo, Amount: 37},
		{Text: "Allez a Derb Sultan.", Effect: EffectMoveTo, Amount: 1},
	}
}

// CommunityChestCards returns the Moroccan-themed Caisse Commune deck.
func CommunityChestCards() []Card {
	return []Card{
		{Text: "Avancez jusqu'a DEPART. Recevez 200 MAD.", Effect: EffectMoveTo, Amount: 0},
		{Text: "Erreur bancaire en votre faveur. Recevez 200 MAD.", Effect: EffectCollect, Amount: 200},
		{Text: "Frais medicaux. Payez 50 MAD.", Effect: EffectPay, Amount: 50},
		{Text: "Vente de votre huile d'argan. Recevez 50 MAD.", Effect: EffectCollect, Amount: 50},
		{Text: "Carte de sortie de prison gratuite.", Effect: EffectGetOutOfJail},
		{Text: "Allez en prison. Ne passez pas par DEPART.", Effect: EffectGoToJail},
		{Text: "Fete de l'Aid! Recevez 100 MAD de chaque joueur.", Effect: EffectCollectAll, Amount: 100},
		{Text: "Remboursement d'impots. Recevez 20 MAD.", Effect: EffectCollect, Amount: 20},
		{Text: "C'est votre anniversaire! Recevez 10 MAD de chaque joueur.", Effect: EffectCollectAll, Amount: 10},
		{Text: "Assurance vie. Recevez 100 MAD.", Effect: EffectCollect, Amount: 100},
		{Text: "Frais de scolarite. Payez 150 MAD.", Effect: EffectPay, Amount: 150},
		{Text: "Recevez votre allocation vacances. Recevez 100 MAD.", Effect: EffectCollect, Amount: 100},
		{Text: "Heritage familial. Recevez 100 MAD.", Effect: EffectCollect, Amount: 100},
		{Text: "Reparations de votre riad: payez 40 MAD par maison et 115 MAD par hotel.", Effect: EffectPayPerHouse, Amount: 40, AmountHotel: 115},
		{Text: "Frais d'hospitalisation. Payez 100 MAD.", Effect: EffectPay, Amount: 100},
		{Text: "Deuxieme prix au concours de beaute. Recevez 10 MAD.", Effect: EffectCollect, Amount: 10},
	}
}
