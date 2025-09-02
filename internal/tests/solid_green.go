package tests

import (
	"image/color"

	"github.com/keshon/screen-tester/internal/core"
)

type SolidGreen struct {
	opts core.TestOptions
}

func (t *SolidGreen) Name() string        { return "Green" }
func (t *SolidGreen) Description() string { return "Solid green screen" }
func (t *SolidGreen) Order() int          { return 2 }

func (t *SolidGreen) Options() core.TestOptions {
	if t.opts.Brightness == 0 {
		t.opts.Brightness = 1.0
	}
	return t.opts
}

func (t *SolidGreen) HandleKeys(ctx *core.WindowContext) {
	core.AdjustBrightnessWithKeys(ctx, 0.1)
}

func (t *SolidGreen) Run(ctx *core.WindowContext) {
	t.HandleKeys(ctx)
	ctx.Win.Clear(core.AdjustBrightness(color.RGBA{0, 255, 0, 255}, ctx.Brightness))
}

func init() {
	core.RegisterTest(&SolidGreen{})
}
