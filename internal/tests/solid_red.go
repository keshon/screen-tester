package tests

import (
	"image/color"

	"github.com/keshon/screen-tester/internal/core"
)

type SolidRed struct {
	opts core.TestOptions
}

func (t *SolidRed) Name() string        { return "Red" }
func (t *SolidRed) Description() string { return "Solid red screen" }
func (t *SolidRed) Order() int          { return 1 }

func (t *SolidRed) Options() core.TestOptions {
	if t.opts.Brightness == 0 {
		t.opts.Brightness = 1.0
	}
	return t.opts
}

func (t *SolidRed) HandleKeys(ctx *core.WindowContext) {
	core.AdjustBrightnessWithKeys(ctx, 0.1)
}

func (t *SolidRed) Run(ctx *core.WindowContext) {
	t.HandleKeys(ctx)
	ctx.Win.Clear(core.AdjustBrightness(color.RGBA{255, 0, 0, 255}, ctx.Brightness))
}

func init() {
	core.RegisterTest(&SolidRed{})
}
