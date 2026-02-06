package save

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/AchrafSoltani/MoroccanMonopoly/board"
	"github.com/AchrafSoltani/MoroccanMonopoly/config"
)

const saveDir = ".config/moroccan-monopoly"
const saveFile = "save.json"

// SaveData represents the full game state for serialisation.
type SaveData struct {
	Players    []PlayerData                         `json:"players"`
	Current    int                                  `json:"current"`
	Properties [config.SpaceCount]PropertyData      `json:"properties"`
	HousePool  int                                  `json:"house_pool"`
	HotelPool  int                                  `json:"hotel_pool"`
	Die1       int                                  `json:"die1"`
	Die2       int                                  `json:"die2"`
	Messages   []string                             `json:"messages"`
}

// PlayerData is the serialisable player state.
type PlayerData struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	IsAI              bool   `json:"is_ai"`
	Money             int    `json:"money"`
	Position          int    `json:"position"`
	InJail            bool   `json:"in_jail"`
	JailTurns         int    `json:"jail_turns"`
	Bankrupt          bool   `json:"bankrupt"`
	Properties        []int  `json:"properties"`
	GetOutOfJailCards int    `json:"get_out_of_jail_cards"`
}

// PropertyData is the serialisable property state.
type PropertyData struct {
	OwnerID   int  `json:"owner_id"`
	Houses    int  `json:"houses"`
	Mortgaged bool `json:"mortgaged"`
}

func savePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, saveDir, saveFile)
}

// Save writes game state to disk.
func Save(data *SaveData) error {
	path := savePath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonData, 0644)
}

// Load reads game state from disk.
func Load() (*SaveData, error) {
	path := savePath()
	jsonData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	data := &SaveData{}
	if err := json.Unmarshal(jsonData, data); err != nil {
		return nil, err
	}

	return data, nil
}

// HasSave checks if a save file exists.
func HasSave() bool {
	_, err := os.Stat(savePath())
	return err == nil
}

// DeleteSave removes the save file.
func DeleteSave() {
	os.Remove(savePath())
}

// BoardToPropertyData converts board properties to saveable format.
func BoardToPropertyData(b *board.Board) [config.SpaceCount]PropertyData {
	var data [config.SpaceCount]PropertyData
	for i := 0; i < config.SpaceCount; i++ {
		p := b.Properties[i]
		data[i] = PropertyData{
			OwnerID:   p.OwnerID,
			Houses:    p.Houses,
			Mortgaged: p.Mortgaged,
		}
	}
	return data
}

// PropertyDataToBoard restores board properties from saved data.
func PropertyDataToBoard(b *board.Board, data [config.SpaceCount]PropertyData) {
	for i := 0; i < config.SpaceCount; i++ {
		b.Properties[i] = board.PropertyState{
			OwnerID:   data[i].OwnerID,
			Houses:    data[i].Houses,
			Mortgaged: data[i].Mortgaged,
		}
	}
}
