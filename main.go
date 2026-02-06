package main

import (
	"log"
	"time"

	"github.com/AchrafSoltani/MoroccanMonopoly/config"
	"github.com/AchrafSoltani/MoroccanMonopoly/game"
	"github.com/AchrafSoltani/glow"
)

func main() {
	win, err := glow.NewWindow(config.WindowTitle, config.WindowWidth, config.WindowHeight)
	if err != nil {
		log.Fatal(err)
	}
	defer win.Close()

	g := game.NewGame()
	canvas := win.Canvas()
	running := true
	lastTime := time.Now()

	for running {
		now := time.Now()
		dt := now.Sub(lastTime).Seconds()
		lastTime = now

		if dt > 0.05 {
			dt = 0.05
		}

		for {
			event := win.PollEvent()
			if event == nil {
				break
			}
			switch event.Type {
			case glow.EventQuit:
				running = false
			case glow.EventKeyDown:
				g.KeyDown(event.Key)
			case glow.EventMouseButtonDown:
				g.MouseDown(event.X, event.Y, event.Button)
			case glow.EventMouseMotion:
				g.MouseMove(event.X, event.Y)
			case glow.EventWindowResize:
				g.OnResize(event.Width, event.Height)
			}
		}

		g.Update(dt)

		canvas.Clear(glow.Black)
		g.Draw(canvas)
		win.Present()

		elapsed := time.Since(now)
		target := time.Second / 60
		if elapsed < target {
			time.Sleep(target - elapsed)
		}
	}
}
