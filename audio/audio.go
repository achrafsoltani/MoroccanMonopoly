package audio

import (
	"bytes"
	"log"

	"github.com/AchrafSoltani/glow"
)

// Engine manages audio playback using Glow's PulseAudio backend.
type Engine struct {
	ctx *glow.AudioContext

	diceRollBuf   []byte
	purchaseBuf   []byte
	rentBuf       []byte
	cardDrawBuf   []byte
	jailBuf       []byte
	passGoBuf     []byte
	buildBuf      []byte
	bankruptcyBuf []byte
	winBuf        []byte
	menuSelBuf    []byte
}

// NewEngine initialises the audio subsystem.
func NewEngine() *Engine {
	ctx, err := glow.NewAudioContext(sampleRate, 1, 2)
	if err != nil {
		log.Printf("audio: failed to init: %v", err)
		return &Engine{}
	}

	return &Engine{
		ctx:           ctx,
		diceRollBuf:   GenerateDiceRoll(),
		purchaseBuf:   GeneratePurchase(),
		rentBuf:       GenerateRent(),
		cardDrawBuf:   GenerateCardDraw(),
		jailBuf:       GenerateJail(),
		passGoBuf:     GeneratePassGo(),
		buildBuf:      GenerateBuild(),
		bankruptcyBuf: GenerateBankruptcy(),
		winBuf:        GenerateWin(),
		menuSelBuf:    GenerateMenuSelect(),
	}
}

func (e *Engine) play(buf []byte) {
	if e.ctx == nil || len(buf) == 0 {
		return
	}
	p := e.ctx.NewPlayer(bytes.NewReader(buf))
	p.Play()
}

func (e *Engine) PlayDiceRoll()   { e.play(e.diceRollBuf) }
func (e *Engine) PlayPurchase()   { e.play(e.purchaseBuf) }
func (e *Engine) PlayRent()       { e.play(e.rentBuf) }
func (e *Engine) PlayCardDraw()   { e.play(e.cardDrawBuf) }
func (e *Engine) PlayJail()       { e.play(e.jailBuf) }
func (e *Engine) PlayPassGo()     { e.play(e.passGoBuf) }
func (e *Engine) PlayBuild()      { e.play(e.buildBuf) }
func (e *Engine) PlayBankruptcy() { e.play(e.bankruptcyBuf) }
func (e *Engine) PlayWin()        { e.play(e.winBuf) }
func (e *Engine) PlayMenuSelect() { e.play(e.menuSelBuf) }
