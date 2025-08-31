package tests

import (
	"image/color"

	"github.com/keshon/screen-tester/internal/core"
)

type SolidWhite struct {
	opts core.TestOptions
}

func (t *SolidWhite) Name() string        { return "White" }
func (t *SolidWhite) Description() string { return "Solid white screen" }
func (t *SolidWhite) Order() int          { return 4 }

func (t *SolidWhite) Options() core.TestOptions {
	if t.opts.Brightness == 0 {
		t.opts.Brightness = 1.0
	}
	return t.opts
}

func (t *SolidWhite) HandleKeys(ctx *core.WindowContext) {
	core.AdjustBrightnessWithKeys(ctx, 0.1)
}

func (t *SolidWhite) Run(ctx *core.WindowContext) {
	t.HandleKeys(ctx)
	ctx.Win.Clear(core.AdjustBrightness(color.RGBA{255, 255, 255, 255}, ctx.Brightness))
}

func init() {
	core.Register(&SolidWhite{})
}
