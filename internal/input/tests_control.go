package input

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/keshon/screen-tester/internal/core"
)

type TestInput struct {
	Current int
}

func (ti *TestInput) HandleTestInput(ctx *core.WindowContext, tests []core.ScreenTest) {
	win := ctx.Win

	if win.JustPressed(pixelgl.KeyF1) {
		ctx.ShowInfo = !ctx.ShowInfo
	}
	if win.JustPressed(pixelgl.KeyRight) {
		ti.Current = (ti.Current + 1) % len(tests)
	}
	if win.JustPressed(pixelgl.KeyLeft) {
		ti.Current = (ti.Current - 1 + len(tests)) % len(tests)
	}
}
