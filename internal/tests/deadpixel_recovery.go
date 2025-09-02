package tests

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/keshon/screen-tester/internal/core"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type DeadPixelRecovery struct {
	opts         core.TestOptions
	state        *flickerState
	speed        time.Duration
	defaultSpeed time.Duration
	minSpeed     time.Duration
	maxSpeed     time.Duration
	step         time.Duration
}

type flickerState struct {
	lastUpdate time.Time
	pic        *pixel.PictureData
	colors     []color.RGBA
}

func (t *DeadPixelRecovery) Name() string { return "Dead Pixel Recovery" }
func (t *DeadPixelRecovery) Description() string {
	return "Flashes colors to exercise dead pixels (Shift+Up/Down to adjust speed)"
}
func (t *DeadPixelRecovery) Order() int { return 61 }

func (t *DeadPixelRecovery) Options() core.TestOptions {
	return core.TestOptions{
		Brightness: 1.0,
		Extra: map[string]interface{}{
			"speed": t.getSpeed(),
		},
	}
}

func (t *DeadPixelRecovery) HandleKeys(ctx *core.WindowContext) {
	if ctx.Win.Pressed(pixelgl.KeyLeftShift) || ctx.Win.Pressed(pixelgl.KeyRightShift) {
		if ctx.Win.JustPressed(pixelgl.KeyUp) {
			t.setSpeed(t.getSpeed() - t.step)
		}
		if ctx.Win.JustPressed(pixelgl.KeyDown) {
			t.setSpeed(t.getSpeed() + t.step)
		}
	} else {
		core.AdjustBrightnessWithKeys(ctx, 0.1)
	}
}

func (t *DeadPixelRecovery) Run(ctx *core.WindowContext) {
	t.HandleKeys(ctx)

	if t.state == nil {
		t.state = &flickerState{
			lastUpdate: time.Now(),
			colors: []color.RGBA{
				{R: 255, G: 0, B: 0, A: 255},
				{R: 0, G: 255, B: 0, A: 255},
				{R: 0, G: 0, B: 255, A: 255},
				{R: 255, G: 255, B: 255, A: 255},
				{R: 0, G: 0, B: 0, A: 255},
			},
		}
	}

	bounds := ctx.Win.Bounds()

	if t.state.pic == nil || t.state.pic.Bounds().W() != bounds.W() || t.state.pic.Bounds().H() != bounds.H() {
		t.state.pic = pixel.MakePictureData(bounds)
	}

	if time.Since(t.state.lastUpdate) >= t.speed {
		t.state.lastUpdate = time.Now()
		for i := range t.state.pic.Pix {
			c := t.state.colors[rand.Intn(len(t.state.colors))]
			t.state.pic.Pix[i] = core.AdjustBrightness(c, ctx.Brightness)
		}
	}

	sprite := pixel.NewSprite(t.state.pic, t.state.pic.Bounds())
	ctx.Win.Clear(color.Black)
	sprite.Draw(ctx.Win, pixel.IM.Moved(bounds.Center()))
}

func (t *DeadPixelRecovery) getSpeed() time.Duration {
	if t.speed == 0 {
		t.speed = t.defaultSpeed
	}
	return t.speed
}

func (t *DeadPixelRecovery) setSpeed(d time.Duration) {
	if d < t.minSpeed {
		d = t.minSpeed
	}
	if d > t.maxSpeed {
		d = t.maxSpeed
	}
	t.speed = d
	if t.opts.Extra == nil {
		t.opts.Extra = map[string]interface{}{}
	}
	t.opts.Extra["speed"] = t.speed
}

func init() {
	core.RegisterTest(&DeadPixelRecovery{
		defaultSpeed: 50 * time.Millisecond,
		minSpeed:     10 * time.Millisecond,
		maxSpeed:     100 * time.Millisecond,
		step:         10 * time.Millisecond,
	})
}
