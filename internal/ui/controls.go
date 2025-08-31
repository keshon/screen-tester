package ui

import (
	"os"

	"app/internal/core"

	"github.com/faiface/pixel/pixelgl"
)

type Controls struct {
	Current int
}

func (c *Controls) HandleInput(ctx *core.WindowContext, tests []core.ScreenTest) {
	win := ctx.Win

	if win.JustPressed(pixelgl.KeyF1) {
		ctx.ShowInfo = !ctx.ShowInfo
	}
	if win.Pressed(pixelgl.KeyEscape) {
		os.Exit(0)
	}
	if win.JustPressed(pixelgl.KeyRight) {
		c.Current = (c.Current + 1) % len(tests)
	}
	if win.JustPressed(pixelgl.KeyLeft) {
		c.Current = (c.Current - 1 + len(tests)) % len(tests)
	}
}
