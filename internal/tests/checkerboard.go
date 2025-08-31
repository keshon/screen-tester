package tests

import (
	"image/color"

	"github.com/keshon/screen-tester/internal/core"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Checkerboard struct {
	opts        core.TestOptions
	defaultSize int
	minSize     int
	maxSize     int
	step        int
}

func (t *Checkerboard) Name() string { return "Small Checkerboard" }
func (t *Checkerboard) Description() string {
	return "Black & white checkerboard with adjustable square size (Shift+Up/Down)"
}
func (t *Checkerboard) Order() int { return 30 }

func (t *Checkerboard) Options() core.TestOptions {
	if t.opts.Brightness == 0 {
		t.opts.Brightness = 1.0
	}
	if _, ok := t.opts.Extra["size"]; !ok {
		t.opts.Extra = map[string]interface{}{"size": t.defaultSize}
	}
	return t.opts
}

func (t *Checkerboard) HandleKeys(ctx *core.WindowContext) {
	if !(ctx.Win.Pressed(pixelgl.KeyLeftShift) || ctx.Win.Pressed(pixelgl.KeyRightShift)) {
		core.AdjustBrightnessWithKeys(ctx, 0.1)
	}

	size := t.getSize()

	if ctx.Win.Pressed(pixelgl.KeyLeftShift) || ctx.Win.Pressed(pixelgl.KeyRightShift) {
		if ctx.Win.JustPressed(pixelgl.KeyUp) {
			size += t.step
		}
		if ctx.Win.JustPressed(pixelgl.KeyDown) {
			size -= t.step
		}
		t.setSize(size)
	}
}

func (t *Checkerboard) Run(ctx *core.WindowContext) {
	t.HandleKeys(ctx)

	size := t.opts.Extra["size"].(int)
	bounds := ctx.Win.Bounds()
	width := int(bounds.W())
	height := int(bounds.H())
	pic := pixel.MakePictureData(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var col color.RGBA
			if (x/size+y/size)%2 == 0 {
				col = core.AdjustBrightness(color.RGBA{255, 255, 255, 255}, ctx.Brightness)
			} else {
				col = core.AdjustBrightness(color.RGBA{0, 0, 0, 255}, ctx.Brightness)
			}
			pic.Pix[y*width+x] = col
		}
	}

	sprite := pixel.NewSprite(pic, bounds)
	sprite.Draw(ctx.Win, pixel.IM.Moved(bounds.Center()))
}

func (t *Checkerboard) getSize() int {
	if t.opts.Extra == nil {
		t.opts.Extra = map[string]interface{}{}
	}

	sizeIface, ok := t.opts.Extra["size"]
	if !ok || sizeIface == nil {
		t.opts.Extra["size"] = t.defaultSize
		return t.defaultSize
	}
	return sizeIface.(int)
}

func (t *Checkerboard) setSize(size int) {
	if size < t.minSize {
		size = t.minSize
	}
	if size > t.maxSize {
		size = t.maxSize
	}
	t.opts.Extra["size"] = size
}

func init() {
	core.Register(&Checkerboard{
		defaultSize: 20,
		minSize:     2,
		maxSize:     50,
		step:        5,
	})
}
