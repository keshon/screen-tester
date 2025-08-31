package tests

import (
	"image/color"

	"github.com/keshon/screen-tester/internal/core"
)

type SolidBlue struct {
	opts core.TestOptions
}

func (t *SolidBlue) Name() string        { return "Blue" }
func (t *SolidBlue) Description() string { return "Solid blue screen" }
func (t *SolidBlue) Order() int          { return 3 }

func (t *SolidBlue) Options() core.TestOptions {
	if t.opts.Brightness == 0 {
		t.opts.Brightness = 1.0
	}
	return t.opts
}

func (t *SolidBlue) HandleKeys(ctx *core.WindowContext) {
	core.AdjustBrightnessWithKeys(ctx, 0.1)
}

func (t *SolidBlue) Run(ctx *core.WindowContext) {
	t.HandleKeys(ctx)
	ctx.Win.Clear(core.AdjustBrightness(color.RGBA{0, 0, 255, 255}, ctx.Brightness))
}

func init() {
	core.Register(&SolidBlue{})
}
