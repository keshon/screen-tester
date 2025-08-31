package tests

import (
	"app/internal/core"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type PixelGrid struct {
	opts        core.TestOptions
	defaultSize int
	minSize     int
	maxSize     int
	step        int
}

func (t *PixelGrid) Name() string { return "Large Pixel Grid" }
func (t *PixelGrid) Description() string {
	return "Grid overlay with adjustable cells size (Shift+Up/Down)"
}
func (t *PixelGrid) Order() int { return 40 }

func (t *PixelGrid) Options() core.TestOptions {
	if t.opts.Brightness == 0 {
		t.opts.Brightness = 1.0
	}
	if _, ok := t.opts.Extra["size"]; !ok {
		t.opts.Extra = map[string]interface{}{"size": t.defaultSize}
	}
	return t.opts
}

func (t *PixelGrid) HandleKeys(ctx *core.WindowContext) {
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

func (t *PixelGrid) Run(ctx *core.WindowContext) {
	t.HandleKeys(ctx)

	size := t.getSize()
	ctx.Win.Clear(core.AdjustBrightness(color.RGBA{0, 0, 0, 255}, ctx.Brightness)) // clear to black

	bounds := ctx.Win.Bounds()
	imd := imdraw.New(nil)
	imd.Color = core.AdjustBrightness(color.RGBA{255, 255, 255, 255}, ctx.Brightness)

	for x := float64(0); x <= bounds.W(); x += float64(size) {
		imd.Push(pixel.V(x, 0), pixel.V(x, bounds.H()))
		imd.Line(1)
	}
	for y := float64(0); y <= bounds.H(); y += float64(size) {
		imd.Push(pixel.V(0, y), pixel.V(bounds.W(), y))
		imd.Line(1)
	}
	imd.Draw(ctx.Win)
}

func (t *PixelGrid) getSize() int {
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

func (t *PixelGrid) setSize(size int) {
	if size < t.minSize {
		size = t.minSize
	}
	if size > t.maxSize {
		size = t.maxSize
	}
	t.opts.Extra["size"] = size
}

func init() {
	core.Register(&PixelGrid{
		defaultSize: 20,
		minSize:     2,
		maxSize:     50,
		step:        5,
	})
}
