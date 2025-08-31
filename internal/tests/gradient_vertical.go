package tests

import (
	"image/color"

	"github.com/keshon/screen-tester/internal/core"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type GradientVertical struct {
	opts      core.TestOptions
	direction string // "black to white" or "white to black"
}

func (t *GradientVertical) Name() string { return "Vertical Gradient" }
func (t *GradientVertical) Description() string {
	return "Black to white gradient (Shift+Up/Down to invert)"
}
func (t *GradientVertical) Order() int { return 20 }

func (t *GradientVertical) Options() core.TestOptions {
	if t.opts.Brightness == 0 {
		t.opts.Brightness = 1.0
	}
	if t.opts.Extra == nil {
		t.opts.Extra = map[string]interface{}{}
	}
	if _, ok := t.opts.Extra["direction"]; !ok {
		t.opts.Extra["direction"] = "black to white"
	}
	t.direction = t.getDirection()
	return t.opts
}

func (t *GradientVertical) HandleKeys(ctx *core.WindowContext) {
	if ctx.Win.Pressed(pixelgl.KeyLeftShift) || ctx.Win.Pressed(pixelgl.KeyRightShift) {
		if ctx.Win.JustPressed(pixelgl.KeyUp) || ctx.Win.JustPressed(pixelgl.KeyDown) {
			if t.getDirection() == "black to white" {
				t.setDirection("white to black")
			} else {
				t.setDirection("black to white")
			}
		}
	} else {
		core.AdjustBrightnessWithKeys(ctx, 0.1)
	}
}

func (t *GradientVertical) Run(ctx *core.WindowContext) {
	t.HandleKeys(ctx)

	bounds := ctx.Win.Bounds()
	width := int(bounds.W())
	height := int(bounds.H())

	pic := pixel.MakePictureData(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			val := uint8((y * 255) / height)
			if t.getDirection() == "black to white" {
				val = 255 - val
			}
			c := core.AdjustBrightness(color.RGBA{val, val, val, 255}, ctx.Brightness)
			pic.Pix[y*width+x] = c
		}
	}

	sprite := pixel.NewSprite(pic, bounds)
	sprite.Draw(ctx.Win, pixel.IM.Moved(bounds.Center()))
}

func (t *GradientVertical) getDirection() string {
	if t.opts.Extra == nil {
		t.opts.Extra = map[string]interface{}{}
	}
	if dir, ok := t.opts.Extra["direction"].(string); ok {
		t.direction = dir
	} else {
		t.direction = "black to white"
		t.opts.Extra["direction"] = t.direction
	}
	return t.direction
}

func (t *GradientVertical) setDirection(dir string) {
	if t.opts.Extra == nil {
		t.opts.Extra = map[string]interface{}{}
	}
	t.direction = dir
	t.opts.Extra["direction"] = dir
}

func init() {
	core.Register(&GradientVertical{})
}
