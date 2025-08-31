package tests

import (
	"image/color"
	"time"

	"app/internal/core"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type MotionBalls struct {
	opts         core.TestOptions
	state        *motionState
	speed        time.Duration
	defaultSpeed time.Duration
	minSpeed     time.Duration
	maxSpeed     time.Duration
	step         time.Duration
}

type motionState struct {
	lastUpdate   time.Time
	balls        []ball
	blackBall    ball
	bgIndex      int
	backgrounds  []bgState
	imd          *imdraw.IMDraw
	frameElapsed time.Duration
}

type ball struct {
	pos       pixel.Vec
	vel       pixel.Vec
	color     color.RGBA
	radius    float64
	origColor color.RGBA
}

type bgState struct {
	color color.RGBA
	name  string
}

func (t *MotionBalls) Name() string { return "Motion Balls" }
func (t *MotionBalls) Description() string {
	return "Bouncing balls with background cycling (Shift+Up/Down: background, Up/Down: speed)"
}
func (t *MotionBalls) Order() int { return 51 }

func (t *MotionBalls) Options() core.TestOptions {
	return core.TestOptions{
		Brightness: 1.0,
		Extra: map[string]interface{}{
			"speed": t.getSpeed(),
		},
	}
}

func (t *MotionBalls) HandleKeys(ctx *core.WindowContext) {
	if ctx.Win.Pressed(pixelgl.KeyLeftShift) || ctx.Win.Pressed(pixelgl.KeyRightShift) {

		if ctx.Win.JustPressed(pixelgl.KeyUp) {
			t.setSpeed(t.getSpeed() + t.step)
			t.rescaleVelocities()
		}
		if ctx.Win.JustPressed(pixelgl.KeyDown) {
			t.setSpeed(t.getSpeed() - t.step)
			t.rescaleVelocities()
		}
	} else {

		if ctx.Win.JustPressed(pixelgl.KeyUp) {
			t.state.bgIndex = (t.state.bgIndex + 1) % len(t.state.backgrounds)
		} else if ctx.Win.JustPressed(pixelgl.KeyDown) {
			t.state.bgIndex--
			if t.state.bgIndex < 0 {
				t.state.bgIndex = len(t.state.backgrounds) - 1
			}
		}
	}
}

func (t *MotionBalls) Run(ctx *core.WindowContext) {
	t.HandleKeys(ctx)

	if t.state == nil {
		t.init(ctx)
	}

	now := time.Now()
	dt := now.Sub(t.state.lastUpdate)
	t.state.lastUpdate = now

	win := ctx.Win
	bounds := win.Bounds()
	win.Clear(t.state.backgrounds[t.state.bgIndex].color)
	t.state.imd.Clear()

	t.updateBalls(dt.Seconds(), bounds)
	t.state.imd.Draw(win)
}

func (t *MotionBalls) init(ctx *core.WindowContext) {
	t.state = &motionState{
		lastUpdate:  time.Now(),
		imd:         imdraw.New(nil),
		backgrounds: []bgState{{colornames.Black, "Black"}, {colornames.White, "White"}, {colornames.Red, "Red"}, {colornames.Green, "Green"}, {colornames.Blue, "Blue"}},
		bgIndex:     0,
	}

	speed := float64(t.getSpeed().Milliseconds()) // пиксели/сек по сути
	t.state.balls = []ball{
		{pixel.V(100, 100), pixel.V(1, 1).Scaled(speed), colornames.Red, 50, colornames.Red},
		{pixel.V(300, 300), pixel.V(-1, 1).Scaled(speed), colornames.Green, 50, colornames.Green},
		{pixel.V(500, 200), pixel.V(1, -1).Scaled(speed), colornames.Blue, 50, colornames.Blue},
		{pixel.V(700, 400), pixel.V(-1, -1).Scaled(speed), colornames.White, 50, colornames.White},
	}
	t.state.blackBall = ball{pixel.V(750, 450), pixel.V(-1, -1).Scaled(speed), colornames.Black, 50, colornames.Black}
}

func (t *MotionBalls) rescaleVelocities() {
	speed := float64(t.getSpeed().Milliseconds())
	for i := range t.state.balls {
		dir := t.state.balls[i].vel.Unit()
		t.state.balls[i].vel = dir.Scaled(speed)
	}
	dir := t.state.blackBall.vel.Unit()
	t.state.blackBall.vel = dir.Scaled(speed)
}

func (t *MotionBalls) updateBalls(dt float64, bounds pixel.Rect) {
	bgName := t.state.backgrounds[t.state.bgIndex].name
	switch bgName {
	case "White":
		for i := range t.state.balls {
			if t.state.balls[i].origColor == colornames.White {
				t.state.balls[i].color = colornames.Black
			} else {
				t.state.balls[i].color = t.state.balls[i].origColor
			}
		}
	default:
		for i := range t.state.balls {
			t.state.balls[i].color = t.state.balls[i].origColor
		}
	}

	for i := range t.state.balls {
		t.moveAndBounce(&t.state.balls[i], dt, bounds)
		t.state.imd.Color = core.AdjustBrightness(t.state.balls[i].color, 1.0)
		t.state.imd.Push(t.state.balls[i].pos)
		t.state.imd.Circle(t.state.balls[i].radius, 0)
	}

	if bgName == "Red" || bgName == "Green" || bgName == "Blue" {
		t.moveAndBounce(&t.state.blackBall, dt, bounds)
		t.state.imd.Color = colornames.Black
		t.state.imd.Push(t.state.blackBall.pos)
		t.state.imd.Circle(t.state.blackBall.radius, 0)
	}
}

func (t *MotionBalls) moveAndBounce(b *ball, dt float64, bounds pixel.Rect) {
	b.pos = b.pos.Add(b.vel.Scaled(dt))
	if b.pos.X-b.radius < 0 {
		b.pos.X = b.radius
		b.vel.X = -b.vel.X
	} else if b.pos.X+b.radius > bounds.W() {
		b.pos.X = bounds.W() - b.radius
		b.vel.X = -b.vel.X
	}
	if b.pos.Y-b.radius < 0 {
		b.pos.Y = b.radius
		b.vel.Y = -b.vel.Y
	} else if b.pos.Y+b.radius > bounds.H() {
		b.pos.Y = bounds.H() - b.radius
		b.vel.Y = -b.vel.Y
	}
}

func (t *MotionBalls) getSpeed() time.Duration {
	if t.speed == 0 {
		t.speed = t.defaultSpeed
	}
	return t.speed
}

func (t *MotionBalls) setSpeed(d time.Duration) {
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
	core.Register(&MotionBalls{
		defaultSpeed: 500 * time.Millisecond,
		minSpeed:     50 * time.Millisecond,
		maxSpeed:     2000 * time.Millisecond,
		step:         100 * time.Millisecond,
	})
}
