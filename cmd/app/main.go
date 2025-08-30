package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var (
	brightness      = 1.0
	showInfo        = true
	atlas           = text.NewAtlas(basicfont.Face7x13, text.ASCII)
	tests           []ScreenTest
	current         = 0
	screenWidth     = 1920
	screenHeight    = 1080
	flickerInterval = 100 * time.Millisecond
)

type ScreenTest struct {
	name               string
	draw               func(win *pixelgl.Window)
	supportsBrightness bool
	supportsSpeed      bool
	helpText           string
}

func main() {
	pixelgl.Run(run)
}

func run() {
	monitor := pixelgl.PrimaryMonitor()
	width, height := monitor.Size()
	screenWidth = int(width)
	screenHeight = int(height)

	cfg := pixelgl.WindowConfig{
		Title:       "Monitor Tester — Señor Mega's Wrath",
		Bounds:      pixel.R(0, 0, width, height),
		Monitor:     monitor,
		Undecorated: true,
		Maximized:   true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	loadTests()

	for !win.Closed() {
		// Handle input for current test capabilities
		test := tests[current]

		if test.supportsBrightness {
			if win.JustPressed(pixelgl.KeyUp) {
				brightness = clamp(brightness+0.1, 0.0, 1.0)
			}
			if win.JustPressed(pixelgl.KeyDown) {
				brightness = clamp(brightness-0.1, 0.0, 1.0)
			}
		}

		if test.supportsSpeed {
			if win.JustPressed(pixelgl.KeyUp) {
				flickerInterval -= 10 * time.Millisecond
				if flickerInterval < 10*time.Millisecond {
					flickerInterval = 10 * time.Millisecond
				}
			}
			if win.JustPressed(pixelgl.KeyDown) {
				flickerInterval += 10 * time.Millisecond
				if flickerInterval > 500*time.Millisecond {
					flickerInterval = 500 * time.Millisecond
				}
			}
		}

		if win.JustPressed(pixelgl.KeyF1) {
			showInfo = !showInfo
		}
		if win.Pressed(pixelgl.KeyEscape) {
			os.Exit(0)
		}
		if win.JustPressed(pixelgl.KeyRight) {
			current = (current + 1) % len(tests)
		}
		if win.JustPressed(pixelgl.KeyLeft) {
			current = (current - 1 + len(tests)) % len(tests)
		}

		win.Clear(color.Black)
		test.draw(win)
		if showInfo {
			drawInfo(win, test)
		}
		win.Update()
	}
}

func drawInfo(win *pixelgl.Window, test ScreenTest) {
	lines := []string{
		fmt.Sprintf("%s", test.name),
		fmt.Sprintf("Resolution: %dx%d", screenWidth, screenHeight),
	}

	if test.supportsBrightness {
		lines = append(lines, fmt.Sprintf("Brightness: %.1f (UP/DOWN)", brightness))
	}

	if test.supportsSpeed {
		lines = append(lines, fmt.Sprintf("Speed: %.0f ms (UP/DOWN)", flickerInterval.Seconds()*1000))
	}

	if test.helpText != "" {
		lines = append(lines, "") // blank line before help
		// Wrap help text in multiple lines if too long
		helpLines := wrapText(test.helpText, 60)
		lines = append(lines, helpLines...)
	}

	imd := imdraw.New(nil)
	imd.Color = color.RGBA{0, 0, 0, 180}

	x := 10.0
	y := win.Bounds().H() - 10.0
	padding := 6.0
	lineHeight := 14.0

	// Calculate box width based on longest line
	var boxWidth float64
	for _, line := range lines {
		txt := text.New(pixel.V(0, 0), atlas)
		fmt.Fprint(txt, line)
		if w := txt.Bounds().W(); w > boxWidth {
			boxWidth = w
		}
	}

	boxHeight := float64(len(lines))*lineHeight + padding
	imd.Push(
		pixel.V(x, y-boxHeight),
		pixel.V(x+boxWidth+padding, y),
	)
	imd.Rectangle(0)
	imd.Draw(win)

	// Draw each line of text
	for i, line := range lines {
		txt := text.New(pixel.V(x+padding/2, y-lineHeight*float64(i+1)), atlas)
		txt.Color = colornames.White
		fmt.Fprint(txt, line)
		txt.Draw(win, pixel.IM)
	}
}

func wrapText(text string, limit int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	var lines []string
	line := ""

	for _, word := range words {
		if len(line)+len(word)+1 > limit {
			lines = append(lines, line)
			line = word
		} else {
			if line != "" {
				line += " "
			}
			line += word
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func adjustBrightness(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	f := brightness
	return color.RGBA{
		R: uint8(clamp(float64(r>>8)*f, 0, 255)),
		G: uint8(clamp(float64(g>>8)*f, 0, 255)),
		B: uint8(clamp(float64(b>>8)*f, 0, 255)),
		A: uint8(clamp(float64(a>>8), 0, 255)),
	}
}

func loadTests() {
	tests = []ScreenTest{
		{
			name:               "Red",
			draw:               solidColor(colornames.Red),
			supportsBrightness: true,
			helpText:           "Adjust brightness with UP/DOWN keys.",
		},
		{
			name:               "Green",
			draw:               solidColor(colornames.Green),
			supportsBrightness: true,
			helpText:           "Adjust brightness with UP/DOWN keys.",
		},
		{
			name:               "Blue",
			draw:               solidColor(colornames.Blue),
			supportsBrightness: true,
			helpText:           "Adjust brightness with UP/DOWN keys.",
		},
		{
			name:               "White",
			draw:               solidColor(colornames.White),
			supportsBrightness: true,
			helpText:           "Adjust brightness with UP/DOWN keys.",
		},
		{
			name:               "Black",
			draw:               solidColor(colornames.Black),
			supportsBrightness: true,
			helpText:           "Adjust brightness with UP/DOWN keys.",
		},
		{
			name:               "Gray",
			draw:               solidColor(colornames.Gray),
			supportsBrightness: true,
			helpText:           "Adjust brightness with UP/DOWN keys.",
		},
		{
			name:               "Horizontal Gradient",
			draw:               gradient(true),
			supportsBrightness: true,
			helpText:           "Adjust brightness with UP/DOWN keys.",
		},
		{
			name:               "Vertical Gradient",
			draw:               gradient(false),
			supportsBrightness: true,
			helpText:           "Adjust brightness with UP/DOWN keys.",
		},
		{
			name:               "Checkerboard",
			draw:               checkerboard(32),
			supportsBrightness: true,
			helpText:           "Adjust brightness with UP/DOWN keys.",
		},
		{
			name:               "Pixel Grid",
			draw:               pixelGrid(8),
			supportsBrightness: true,
			helpText:           "Adjust brightness with UP/DOWN keys.",
		},
		{
			name:          "Dead Pixel Recovery",
			draw:          deadPixelRecovery(),
			supportsSpeed: true,
			helpText:      "Adjust flicker speed with UP/DOWN keys.",
		},
		{
			name:          "Motion Test",
			draw:          motionTest(),
			supportsSpeed: true,
			helpText:      "Bouncing ball. Adjust speed with UP/DOWN keys.",
		},
		{
			name:               "Subpixel Test",
			draw:               subpixelTest(),
			supportsBrightness: true,
			helpText:           "Subpixel movement. Adjust speed with UP/DOWN keys.",
		},
	}
}

func solidColor(c color.Color) func(win *pixelgl.Window) {
	return func(win *pixelgl.Window) {
		win.Clear(adjustBrightness(c))
	}
}

func gradient(horizontal bool) func(win *pixelgl.Window) {
	return func(win *pixelgl.Window) {
		bounds := win.Bounds()
		width := int(bounds.W())
		height := int(bounds.H())
		pic := pixel.MakePictureData(bounds)

		for y := range height {
			for x := range width {
				var val uint8
				if horizontal {
					val = uint8((x * 255) / width)
				} else {
					val = uint8((y * 255) / height)
				}
				c := adjustBrightness(color.RGBA{val, val, val, 255})
				pic.Pix[y*width+x] = c
			}
		}

		sprite := pixel.NewSprite(pic, bounds)
		sprite.Draw(win, pixel.IM.Moved(bounds.Center()))
	}
}

func checkerboard(size int) func(win *pixelgl.Window) {
	return func(win *pixelgl.Window) {
		bounds := win.Bounds()
		width := int(bounds.W())
		height := int(bounds.H())
		pic := pixel.MakePictureData(bounds)

		for y := range height {
			for x := range width {
				if (x/size+y/size)%2 == 0 {
					pic.Pix[y*width+x] = adjustBrightness(colornames.White)
				} else {
					pic.Pix[y*width+x] = adjustBrightness(colornames.Black)
				}
			}
		}

		sprite := pixel.NewSprite(pic, bounds)
		sprite.Draw(win, pixel.IM.Moved(bounds.Center()))
	}
}

func pixelGrid(size int) func(win *pixelgl.Window) {
	return func(win *pixelgl.Window) {
		bounds := win.Bounds()
		imd := imdraw.New(nil)
		imd.Color = adjustBrightness(colornames.White)

		for x := float64(0); x <= bounds.W(); x += float64(size) {
			imd.Push(pixel.V(x, 0), pixel.V(x, bounds.H()))
			imd.Line(1)
		}
		for y := float64(0); y <= bounds.H(); y += float64(size) {
			imd.Push(pixel.V(0, y), pixel.V(bounds.W(), y))
			imd.Line(1)
		}
		imd.Draw(win)
	}
}

func deadPixelRecovery() func(win *pixelgl.Window) {
	type flickerState struct {
		lastUpdate time.Time
		pic        *pixel.PictureData
		colors     []color.RGBA
	}

	state := &flickerState{
		lastUpdate: time.Now(),
		colors: []color.RGBA{
			colornames.Red,
			colornames.Green,
			colornames.Blue,
			colornames.White,
			colornames.Black,
		},
	}

	return func(win *pixelgl.Window) {
		bounds := win.Bounds()

		if state.pic == nil || state.pic.Bounds().W() != bounds.W() || state.pic.Bounds().H() != bounds.H() {
			state.pic = pixel.MakePictureData(bounds)
		}

		if time.Since(state.lastUpdate) >= flickerInterval {
			state.lastUpdate = time.Now()

			for i := range state.pic.Pix {
				state.pic.Pix[i] = adjustBrightness(state.colors[rand.Intn(len(state.colors))])
			}
		}

		sprite := pixel.NewSprite(state.pic, state.pic.Bounds())
		win.Clear(color.Black)
		sprite.Draw(win, pixel.IM.Moved(bounds.Center()))
	}
}

func motionTest() func(win *pixelgl.Window) {
	type ball struct {
		pos       pixel.Vec
		vel       pixel.Vec
		color     color.RGBA
		radius    float64
		origColor color.RGBA // original color for toggling
	}

	speed := 300.0
	lastTime := time.Now()

	// Define backgrounds we'll cycle through
	type bgState struct {
		color color.RGBA
		name  string
	}

	backgrounds := []bgState{
		{color: colornames.Black, name: "Black"},
		{color: colornames.White, name: "White"},
		{color: colornames.Red, name: "Red"},
		{color: colornames.Green, name: "Green"},
		{color: colornames.Blue, name: "Blue"},
	}

	bgIndex := 0 // start with black background

	// Balls setup
	balls := []ball{
		{pos: pixel.V(100, 100), vel: pixel.V(1, 1).Scaled(speed), color: colornames.Red, radius: 50, origColor: colornames.Red},
		{pos: pixel.V(300, 300), vel: pixel.V(-1, 1).Scaled(speed), color: colornames.Green, radius: 50, origColor: colornames.Green},
		{pos: pixel.V(500, 200), vel: pixel.V(1, -1).Scaled(speed), color: colornames.Blue, radius: 50, origColor: colornames.Blue},
		{pos: pixel.V(700, 400), vel: pixel.V(-1, -1).Scaled(speed), color: colornames.White, radius: 50, origColor: colornames.White},
	}

	// Extra black ball copy for colored backgrounds
	blackBall := ball{
		pos:       pixel.V(750, 450),
		vel:       pixel.V(-1, -1).Scaled(speed),
		color:     colornames.Black,
		radius:    50,
		origColor: colornames.Black,
	}

	imd := imdraw.New(nil)

	return func(win *pixelgl.Window) {
		now := time.Now()
		dt := now.Sub(lastTime).Seconds()
		lastTime = now

		// Background cycling with Shift+Up / Shift+Down
		if win.Pressed(pixelgl.KeyLeftShift) || win.Pressed(pixelgl.KeyRightShift) {
			if win.JustPressed(pixelgl.KeyUp) {
				bgIndex = (bgIndex + 1) % len(backgrounds)
			} else if win.JustPressed(pixelgl.KeyDown) {
				bgIndex--
				if bgIndex < 0 {
					bgIndex = len(backgrounds) - 1
				}
			}
		} else {
			// Speed controls without shift
			if win.JustPressed(pixelgl.KeyUp) {
				speed += 50
				if speed > 1000 {
					speed = 1000
				}
				for i := range balls {
					dir := balls[i].vel.Unit()
					balls[i].vel = dir.Scaled(speed)
				}
				// Also for black ball
				dir := blackBall.vel.Unit()
				blackBall.vel = dir.Scaled(speed)
			}
			if win.JustPressed(pixelgl.KeyDown) {
				speed -= 50
				if speed < 50 {
					speed = 50
				}
				for i := range balls {
					dir := balls[i].vel.Unit()
					balls[i].vel = dir.Scaled(speed)
				}
				dir := blackBall.vel.Unit()
				blackBall.vel = dir.Scaled(speed)
			}
		}

		bounds := win.Bounds()
		win.Clear(backgrounds[bgIndex].color)
		imd.Clear()

		// Update colors of balls depending on background
		switch backgrounds[bgIndex].name {
		case "White":
			// Background white: white ball becomes black, others stay colored
			for i := range balls {
				if balls[i].origColor == colornames.White {
					balls[i].color = colornames.Black
				} else {
					balls[i].color = balls[i].origColor
				}
			}
			// No black ball here
		case "Red", "Green", "Blue":
			// On colored backgrounds: have both white and black balls + colored balls
			for i := range balls {
				balls[i].color = balls[i].origColor
			}
		case "Black":
			// Background black: white ball stays white, others colored, no black ball
			for i := range balls {
				balls[i].color = balls[i].origColor
			}
		}

		// Update and draw balls
		for i := range balls {
			balls[i].pos = balls[i].pos.Add(balls[i].vel.Scaled(dt))

			if balls[i].pos.X-balls[i].radius < 0 {
				balls[i].pos.X = balls[i].radius
				balls[i].vel.X = -balls[i].vel.X
			} else if balls[i].pos.X+balls[i].radius > bounds.W() {
				balls[i].pos.X = bounds.W() - balls[i].radius
				balls[i].vel.X = -balls[i].vel.X
			}
			if balls[i].pos.Y-balls[i].radius < 0 {
				balls[i].pos.Y = balls[i].radius
				balls[i].vel.Y = -balls[i].vel.Y
			} else if balls[i].pos.Y+balls[i].radius > bounds.H() {
				balls[i].pos.Y = bounds.H() - balls[i].radius
				balls[i].vel.Y = -balls[i].vel.Y
			}

			imd.Color = adjustBrightness(balls[i].color)
			imd.Push(balls[i].pos)
			imd.Circle(balls[i].radius, 0)
		}

		// If background is colored (red, green, blue), draw the black ball as well
		if backgrounds[bgIndex].name == "Red" || backgrounds[bgIndex].name == "Green" || backgrounds[bgIndex].name == "Blue" {
			blackBall.pos = blackBall.pos.Add(blackBall.vel.Scaled(dt))

			if blackBall.pos.X-blackBall.radius < 0 {
				blackBall.pos.X = blackBall.radius
				blackBall.vel.X = -blackBall.vel.X
			} else if blackBall.pos.X+blackBall.radius > bounds.W() {
				blackBall.pos.X = bounds.W() - blackBall.radius
				blackBall.vel.X = -blackBall.vel.X
			}
			if blackBall.pos.Y-blackBall.radius < 0 {
				blackBall.pos.Y = blackBall.radius
				blackBall.vel.Y = -blackBall.vel.Y
			} else if blackBall.pos.Y+blackBall.radius > bounds.H() {
				blackBall.pos.Y = bounds.H() - blackBall.radius
				blackBall.vel.Y = -blackBall.vel.Y
			}

			imd.Color = colornames.Black
			imd.Push(blackBall.pos)
			imd.Circle(blackBall.radius, 0)
		}

		imd.Draw(win)
	}
}

func subpixelTest() func(win *pixelgl.Window) {
	adjustBrightness := func(c color.RGBA) color.RGBA {
		factor := 0.7
		return color.RGBA{
			R: uint8(float64(c.R) * factor),
			G: uint8(float64(c.G) * factor),
			B: uint8(float64(c.B) * factor),
			A: c.A,
		}
	}

	shift := 0 // shifts RGB channels: 0,1,2 cycling

	return func(win *pixelgl.Window) {
		// Shift RGB channels with Shift+Up / Shift+Down
		if win.Pressed(pixelgl.KeyLeftShift) || win.Pressed(pixelgl.KeyRightShift) {
			if win.JustPressed(pixelgl.KeyUp) {
				shift = (shift + 1) % 3
			} else if win.JustPressed(pixelgl.KeyDown) {
				shift = (shift + 2) % 3 // +2 mod 3 is like -1 mod 3
			}
		}

		bounds := win.Bounds()
		width := int(bounds.W())
		height := int(bounds.H())
		pic := pixel.MakePictureData(bounds)

		getColor := func(x int) color.RGBA {
			switch (x + shift) % 3 {
			case 0:
				return adjustBrightness(color.RGBA{255, 0, 0, 255}) // Red
			case 1:
				return adjustBrightness(color.RGBA{0, 255, 0, 255}) // Green
			default:
				return adjustBrightness(color.RGBA{0, 0, 255, 255}) // Blue
			}
		}

		for y := range height {
			for x := range width {
				pic.Pix[y*width+x] = getColor(x)
			}
		}

		sprite := pixel.NewSprite(pic, bounds)
		win.Clear(color.Black)
		sprite.Draw(win, pixel.IM.Moved(bounds.Center()))
	}
}
