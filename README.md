# MoroccanMonopoly

A Moroccan-themed Monopoly board game built entirely in Go with the [Glow](https://github.com/AchrafSoltani/glow) engine — a pure Go 2D software renderer that communicates directly with X11. No CGo, no images, no external assets. Everything is drawn with primitives and generated procedurally.

## Features

- **40-space board** with real Moroccan cities and landmarks — from Derb Sultan to Mosquee Hassan II
- **Currency: MAD** (Moroccan Dirham) — prices, rents, and taxes all in Dirhams
- **1–4 players** — one human with up to 3 AI opponents
- **Full Monopoly rules** — properties, houses/hotels, rent, mortgages, auctions, trading, jail, bankruptcy
- **Adaptive AI** — buys strategically with cash buffers, builds when profitable, handles jail decisions
- **Procedural audio** — 10 synthesised sound effects (dice roll, purchase, rent, jail, victory fanfare, etc.)
- **Responsive window** — board and HUD scale proportionally when the window is resized
- **Save/load** — game state persisted to `~/.config/moroccan-monopoly/save.json`
- **Animated menu** — zellige-inspired geometric pattern background
- **8x8 bitmap font** — full printable ASCII set, scaleable

## Controls

| Key | Action |
|-----|--------|
| Enter | New game (1 Human + 1 AI) |
| 2 / 3 / 4 | New game with 2 / 3 / 4 players |
| R | Resume saved game |
| F5 | Save game (during play) |
| Mouse | Click buttons, hover spaces for property cards |

## Building from Source

### Prerequisites

- Go 1.21+
- X11 (or XWayland)

```bash
# Clone both repositories
git clone https://github.com/AchrafSoltani/glow.git
git clone https://github.com/AchrafSoltani/MoroccanMonopoly.git

# Build and run
cd MoroccanMonopoly
go run .
```

> The `go.mod` expects `glow` to be cloned alongside `MoroccanMonopoly` (i.e. `../glow`).

### Build a standalone binary

```bash
go build -ldflags="-s -w" -o moroccan-monopoly .
./moroccan-monopoly
```

## Board

| Colour | Properties | Price |
|--------|-----------|-------|
| Brown | Derb Sultan, Bab Marrakech | 60 MAD |
| Light Blue | Av. Hassan II, Bab Bou Jeloud, Talaa Kebira | 100–120 MAD |
| Pink | Av. Mohammed V, Rue des Consuls, Kasbah Oudayas | 140–160 MAD |
| Orange | Av. de la Liberte, Rue de la Liberte, Grand Socco | 180–200 MAD |
| Red | Jemaa el-Fna, Rue Bab Agnaou, Koutoubia | 220–240 MAD |
| Yellow | Av. Mohammed VI, Corniche Ain Diab, Bd. de la Corniche | 260–280 MAD |
| Green | Vallee du Dades, Gorges du Todra, Merzouga (Sahara) | 300–320 MAD |
| Dark Blue | Chefchaouen, Mosquee Hassan II | 350–400 MAD |

**Railroads**: Gare Casa-Voyageurs, Gare Rabat-Ville, Gare Marrakech, Gare Tanger-Ville (200 MAD each)
**Utilities**: ONEE (Electricite), LYDEC (Eau) (150 MAD each)

## Project Structure

```
MoroccanMonopoly/
├── main.go                      # Entry point, game loop, resize handling
├── config/config.go             # Constants + responsive Layout struct
├── board/                       # Board data model
│   ├── board.go                 # 40 spaces with Moroccan property names
│   ├── cards.go                 # Chance / Caisse Commune decks
│   ├── color_group.go           # 8 colour groups
│   ├── property.go              # Ownership, houses, mortgages
│   └── space.go                 # Space types
├── game/                        # Game logic
│   ├── game.go                  # Game struct, drawing, resize, save/load
│   ├── state.go                 # State and phase enums
│   ├── turn.go                  # Turn management, dice, movement, AI
│   ├── rules.go                 # Buy, rent, build, mortgage, bankruptcy
│   ├── auction.go               # Property auction system
│   └── trade.go                 # Player-to-player trading
├── player/                      # Player model
│   ├── player.go                # Player struct
│   └── ai.go                    # AI decision-making
├── render/                      # Rendering
│   ├── board_renderer.go        # Board with responsive layout
│   ├── hud_renderer.go          # Right-side info panel
│   ├── dialog_renderer.go       # Modal dialogs
│   ├── menu_renderer.go         # Animated zellige menu
│   ├── dice_renderer.go         # Dice display
│   ├── token_renderer.go        # Player tokens
│   ├── ui.go                    # Button component
│   ├── colors.go                # Colour palette
│   ├── font.go                  # Bitmap font renderer
│   └── font_data.go             # 5×7 glyph data
├── audio/                       # Procedural audio
│   ├── audio.go                 # Glow PulseAudio backend
│   └── synth.go                 # 10 synthesised sound effects
├── save/save.go                 # JSON save/load
└── go.mod
```

## Engine

Built with [Glow](https://github.com/AchrafSoltani/glow) — a pure Go 2D graphics library that talks directly to X11 via Unix sockets and to PulseAudio for audio. No CGo, no SDL, no OpenGL. Just Go and Unix sockets.

## Licence

MIT
