package tests

import (
	"image/color"

	"github.com/keshon/screen-tester/internal/core"
)

type SolidBlack struct {
	opts core.TestOptions
}

func (t *SolidBlack) Name() string        { return "Black" }
func (t *SolidBlack) Description() string { return "Solid black screen" }
func (t *SolidBlack) Order() int          { return 5 }

func (t *SolidBlack) Options() core.TestOptions {
	t.opts.Brightness = 1.0
	return t.opts
}

func (t *SolidBlack) HandleKeys(ctx *core.WindowContext) {
	core.AdjustBrightnessWithKeys(ctx, 0.1)
}

func (t *SolidBlack) Run(ctx *core.WindowContext) {
	t.HandleKeys(ctx)
	ctx.Win.Clear(core.AdjustBrightness(color.RGBA{0, 0, 0, 255}, ctx.Brightness))
}

func init() {
	core.RegisterTest(&SolidBlack{})
}
