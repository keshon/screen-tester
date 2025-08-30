package main

import (
	"fmt"
	"image/color"
	"math"
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

// Application state
type App struct {
	brightness      float64
	showInfo        bool
	atlas           *text.Atlas
	tests           []ScreenTest
	current         int
	screenWidth     int
	screenHeight    int
	flickerInterval time.Duration
	lastFlicker     time.Time

	// Performance optimizations
	cachedPictures map[string]*pixel.PictureData
	needsRedraw    bool

	// Pattern parameters
	gridSize      int
	checkerSize   int
	gradientSteps int
}

type ScreenTest struct {
	name               string
	draw               func(app *App, win *pixelgl.Window)
	supportsBrightness bool
	supportsSpeed      bool
	supportsSize       bool
	helpText           string
	category           string
}

func NewApp() *App {
	return &App{
		brightness:      1.0,
		showInfo:        true,
		atlas:           text.NewAtlas(basicfont.Face7x13, text.ASCII),
		flickerInterval: 100 * time.Millisecond,
		cachedPictures:  make(map[string]*pixel.PictureData),
		needsRedraw:     true,
		gridSize:        8,
		checkerSize:     32,
		gradientSteps:   256,
	}
}

func main() {
	pixelgl.Run(run)
}

func run() {
	app := NewApp()

	monitor := pixelgl.PrimaryMonitor()
	width, height := monitor.Size()
	app.screenWidth = int(width)
	app.screenHeight = int(height)

	cfg := pixelgl.WindowConfig{
		Title:       "Enhanced Monitor Tester v2.0",
		Bounds:      pixel.R(0, 0, width, height),
		Monitor:     monitor,
		Undecorated: true,
		Maximized:   true,
		VSync:       true, // Better performance
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	app.loadTests()

	// Main loop
	for !win.Closed() {
		app.handleInput(win)

		if app.needsRedraw {
			win.Clear(color.Black)
			app.tests[app.current].draw(app, win)
			app.needsRedraw = false
		}

		if app.showInfo {
			app.drawInfo(win)
		}

		win.Update()
	}
}

func (app *App) handleInput(win *pixelgl.Window) {
	test := app.tests[app.current]

	// Navigation
	if win.JustPressed(pixelgl.KeyRight) || win.JustPressed(pixelgl.KeySpace) {
		app.current = (app.current + 1) % len(app.tests)
		app.needsRedraw = true
		app.clearCache() // Clear cache when switching tests
	}
	if win.JustPressed(pixelgl.KeyLeft) {
		app.current = (app.current - 1 + len(app.tests)) % len(app.tests)
		app.needsRedraw = true
		app.clearCache()
	}

	// Jump to test by number
	for i := pixelgl.Key1; i <= pixelgl.Key9; i++ {
		if win.JustPressed(i) {
			testIndex := int(i - pixelgl.Key1)
			if testIndex < len(app.tests) {
				app.current = testIndex
				app.needsRedraw = true
				app.clearCache()
			}
		}
	}

	// Brightness control
	if test.supportsBrightness {
		if win.JustPressed(pixelgl.KeyUp) || win.JustPressed(pixelgl.KeyEqual) {
			app.brightness = clamp(app.brightness+0.1, 0.0, 1.0)
			app.needsRedraw = true
		}
		if win.JustPressed(pixelgl.KeyDown) || win.JustPressed(pixelgl.KeyMinus) {
			app.brightness = clamp(app.brightness-0.1, 0.0, 1.0)
			app.needsRedraw = true
		}
	}

	// Speed control
	if test.supportsSpeed {
		if win.JustPressed(pixelgl.KeyUp) || win.JustPressed(pixelgl.KeyEqual) {
			app.flickerInterval -= 10 * time.Millisecond
			if app.flickerInterval < 10*time.Millisecond {
				app.flickerInterval = 10 * time.Millisecond
			}
		}
		if win.JustPressed(pixelgl.KeyDown) || win.JustPressed(pixelgl.KeyMinus) {
			app.flickerInterval += 10 * time.Millisecond
			if app.flickerInterval > 1000*time.Millisecond {
				app.flickerInterval = 1000 * time.Millisecond
			}
		}
	}

	// Size control
	if test.supportsSize {
		if win.JustPressed(pixelgl.KeyUp) || win.JustPressed(pixelgl.KeyEqual) {
			if strings.Contains(test.name, "Grid") {
				app.gridSize = clampInt(app.gridSize+2, 2, 64)
			} else if strings.Contains(test.name, "Checker") {
				app.checkerSize = clampInt(app.checkerSize+8, 8, 128)
			}
			app.needsRedraw = true
			app.clearCache()
		}
		if win.JustPressed(pixelgl.KeyDown) || win.JustPressed(pixelgl.KeyMinus) {
			if strings.Contains(test.name, "Grid") {
				app.gridSize = clampInt(app.gridSize-2, 2, 64)
			} else if strings.Contains(test.name, "Checker") {
				app.checkerSize = clampInt(app.checkerSize-8, 8, 128)
			}
			app.needsRedraw = true
			app.clearCache()
		}
	}

	// Toggle info
	if win.JustPressed(pixelgl.KeyF1) || win.JustPressed(pixelgl.KeyH) {
		app.showInfo = !app.showInfo
	}

	// Reset to defaults
	if win.JustPressed(pixelgl.KeyR) {
		app.brightness = 1.0
		app.flickerInterval = 100 * time.Millisecond
		app.gridSize = 8
		app.checkerSize = 32
		app.needsRedraw = true
		app.clearCache()
	}

	// Exit
	if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
		os.Exit(0)
	}
}

func (app *App) drawInfo(win *pixelgl.Window) {
	test := app.tests[app.current]

	lines := []string{
		fmt.Sprintf("Test: %s (%d/%d)", test.name, app.current+1, len(app.tests)),
		fmt.Sprintf("Resolution: %dx%d", app.screenWidth, app.screenHeight),
		"",
		"Controls:",
		"← → / SPACE: Switch tests",
		"1-9: Jump to test",
		"F1/H: Toggle info",
		"R: Reset settings",
		"ESC/Q: Exit",
	}

	if test.supportsBrightness {
		lines = append(lines, fmt.Sprintf("↑↓ / +−: Brightness %.1f", app.brightness))
	}

	if test.supportsSpeed {
		lines = append(lines, fmt.Sprintf("↑↓ / +−: Speed %.0fms", app.flickerInterval.Seconds()*1000))
	}

	if test.supportsSize {
		if strings.Contains(test.name, "Grid") {
			lines = append(lines, fmt.Sprintf("↑↓ / +−: Grid size %dpx", app.gridSize))
		} else if strings.Contains(test.name, "Checker") {
			lines = append(lines, fmt.Sprintf("↑↓ / +−: Checker size %dpx", app.checkerSize))
		}
	}

	if test.helpText != "" {
		lines = append(lines, "")
		helpLines := wrapText(test.helpText, 50)
		lines = append(lines, helpLines...)
	}

	app.drawInfoBox(win, lines)
}

func (app *App) drawInfoBox(win *pixelgl.Window, lines []string) {
	imd := imdraw.New(nil)
	imd.Color = color.RGBA{0, 0, 0, 200}

	x := 15.0
	y := win.Bounds().H() - 15.0
	padding := 10.0
	lineHeight := 15.0

	// Calculate box dimensions
	var boxWidth float64
	for _, line := range lines {
		txt := text.New(pixel.V(0, 0), app.atlas)
		fmt.Fprint(txt, line)
		if w := txt.Bounds().W(); w > boxWidth {
			boxWidth = w
		}
	}

	boxHeight := float64(len(lines))*lineHeight + padding*2

	// Draw rounded background
	imd.Push(
		pixel.V(x, y-boxHeight),
		pixel.V(x+boxWidth+padding*2, y),
	)
	imd.Rectangle(0)
	imd.Draw(win)

	// Draw text
	for i, line := range lines {
		txt := text.New(pixel.V(x+padding, y-lineHeight*float64(i+1)-padding/2), app.atlas)
		if strings.HasPrefix(line, "Test:") {
			txt.Color = colornames.Yellow
		} else if strings.HasPrefix(line, "Controls:") || strings.Contains(line, ":") {
			txt.Color = colornames.Lightblue
		} else {
			txt.Color = colornames.White
		}
		fmt.Fprint(txt, line)
		txt.Draw(win, pixel.IM)
	}
}

func (app *App) clearCache() {
	app.cachedPictures = make(map[string]*pixel.PictureData)
}

func (app *App) loadTests() {
	app.tests = []ScreenTest{
		// Basic colors
		{
			name:               "Pure Red",
			draw:               app.solidColor(colornames.Red),
			supportsBrightness: true,
			helpText:           "Test red subpixels and color accuracy.",
			category:           "Color",
		},
		{
			name:               "Pure Green",
			draw:               app.solidColor(colornames.Green),
			supportsBrightness: true,
			helpText:           "Test green subpixels and color accuracy.",
			category:           "Color",
		},
		{
			name:               "Pure Blue",
			draw:               app.solidColor(colornames.Blue),
			supportsBrightness: true,
			helpText:           "Test blue subpixels and color accuracy.",
			category:           "Color",
		},
		{
			name:               "Pure White",
			draw:               app.solidColor(colornames.White),
			supportsBrightness: true,
			helpText:           "Test maximum brightness and white balance.",
			category:           "Color",
		},
		{
			name:               "Pure Black",
			draw:               app.solidColor(colornames.Black),
			supportsBrightness: true,
			helpText:           "Test black levels and backlight bleeding.",
			category:           "Color",
		},

		// Gradients and patterns
		{
			name:               "Horizontal Gradient",
			draw:               app.gradient(true),
			supportsBrightness: true,
			helpText:           "Test smooth color transitions and banding.",
			category:           "Gradient",
		},
		{
			name:               "Vertical Gradient",
			draw:               app.gradient(false),
			supportsBrightness: true,
			helpText:           "Test smooth color transitions and banding.",
			category:           "Gradient",
		},
		{
			name:               "Checkerboard",
			draw:               app.checkerboard(),
			supportsBrightness: true,
			supportsSize:       true,
			helpText:           "Test sharpness and pixel alignment.",
			category:           "Pattern",
		},
		{
			name:               "Pixel Grid",
			draw:               app.pixelGrid(),
			supportsBrightness: true,
			supportsSize:       true,
			helpText:           "Test pixel accuracy and grid alignment.",
			category:           "Pattern",
		},

		// Motion and recovery
		{
			name:          "Dead Pixel Recovery",
			draw:          app.deadPixelRecovery(),
			supportsSpeed: true,
			helpText:      "Rapidly cycle colors to potentially recover stuck pixels.",
			category:      "Recovery",
		},
		{
			name:          "Motion Test",
			draw:          app.motionTest(),
			supportsSpeed: true,
			helpText:      "Test motion blur and response time with moving patterns.",
			category:      "Motion",
		},

		// New tests
		{
			name:               "RGB Subpixel",
			draw:               app.subpixelTest(),
			supportsBrightness: true,
			helpText:           "Test individual RGB subpixel arrangement.",
			category:           "Advanced",
		},
	}
}

// Enhanced drawing functions
func (app *App) solidColor(c color.Color) func(app *App, win *pixelgl.Window) {
	return func(app *App, win *pixelgl.Window) {
		win.Clear(app.adjustBrightness(c))
	}
}

func (app *App) gradient(horizontal bool) func(app *App, win *pixelgl.Window) {
	return func(app *App, win *pixelgl.Window) {
		cacheKey := fmt.Sprintf("gradient_%t_%d_%d_%.1f", horizontal, app.screenWidth, app.screenHeight, app.brightness)

		pic, exists := app.cachedPictures[cacheKey]
		if !exists {
			bounds := win.Bounds()
			width := int(bounds.W())
			height := int(bounds.H())
			pic = pixel.MakePictureData(bounds)

			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					var val uint8
					if horizontal {
						val = uint8((x * 255) / width)
					} else {
						val = uint8((y * 255) / height)
					}
					c := app.adjustBrightness(color.RGBA{val, val, val, 255})
					pic.Pix[y*width+x] = c
				}
			}
			app.cachedPictures[cacheKey] = pic
		}

		sprite := pixel.NewSprite(pic, win.Bounds())
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	}
}

func (app *App) checkerboard() func(app *App, win *pixelgl.Window) {
	return func(app *App, win *pixelgl.Window) {
		cacheKey := fmt.Sprintf("checker_%d_%d_%d_%.1f", app.checkerSize, app.screenWidth, app.screenHeight, app.brightness)

		pic, exists := app.cachedPictures[cacheKey]
		if !exists {
			bounds := win.Bounds()
			width := int(bounds.W())
			height := int(bounds.H())
			pic = pixel.MakePictureData(bounds)

			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					if (x/app.checkerSize+y/app.checkerSize)%2 == 0 {
						pic.Pix[y*width+x] = app.adjustBrightness(colornames.White)
					} else {
						pic.Pix[y*width+x] = app.adjustBrightness(colornames.Black)
					}
				}
			}
			app.cachedPictures[cacheKey] = pic
		}

		sprite := pixel.NewSprite(pic, win.Bounds())
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	}
}

func (app *App) pixelGrid() func(app *App, win *pixelgl.Window) {
	return func(app *App, win *pixelgl.Window) {
		bounds := win.Bounds()
		imd := imdraw.New(nil)
		imd.Color = app.adjustBrightness(colornames.White)

		// Vertical lines
		for x := float64(0); x <= bounds.W(); x += float64(app.gridSize) {
			imd.Push(pixel.V(x, 0), pixel.V(x, bounds.H()))
			imd.Line(1)
		}

		// Horizontal lines
		for y := float64(0); y <= bounds.H(); y += float64(app.gridSize) {
			imd.Push(pixel.V(0, y), pixel.V(bounds.W(), y))
			imd.Line(1)
		}

		imd.Draw(win)
	}
}

func (app *App) deadPixelRecovery() func(app *App, win *pixelgl.Window) {
	colors := []color.RGBA{
		{255, 0, 0, 255},     // Red
		{0, 255, 0, 255},     // Green
		{0, 0, 255, 255},     // Blue
		{255, 255, 255, 255}, // White
		{0, 0, 0, 255},       // Black
		{255, 255, 0, 255},   // Yellow
		{255, 0, 255, 255},   // Magenta
		{0, 255, 255, 255},   // Cyan
	}

	return func(app *App, win *pixelgl.Window) {
		if time.Since(app.lastFlicker) >= app.flickerInterval {
			app.lastFlicker = time.Now()
			app.needsRedraw = true
		}

		colorIndex := int(time.Since(app.lastFlicker) / (app.flickerInterval / time.Duration(len(colors))))
		if colorIndex >= len(colors) {
			colorIndex = len(colors) - 1
		}

		win.Clear(colors[colorIndex])
	}
}

func (app *App) motionTest() func(app *App, win *pixelgl.Window) {
	return func(app *App, win *pixelgl.Window) {
		bounds := win.Bounds()
		imd := imdraw.New(nil)

		// Moving vertical bar
		speed := 200.0 / app.flickerInterval.Seconds() // pixels per second
		t := time.Since(time.Time{}).Seconds()
		x := math.Mod(t*speed, bounds.W()+100) - 50

		imd.Color = colornames.White
		imd.Push(pixel.V(x, 0), pixel.V(x+50, bounds.H()))
		imd.Rectangle(0)
		imd.Draw(win)

		app.needsRedraw = true // Always redraw for motion
	}
}

func (app *App) subpixelTest() func(app *App, win *pixelgl.Window) {
	return func(app *App, win *pixelgl.Window) {
		bounds := win.Bounds()
		width := int(bounds.W())
		height := int(bounds.H())
		pic := pixel.MakePictureData(bounds)

		// Create RGB subpixel pattern
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				switch x % 3 {
				case 0:
					pic.Pix[y*width+x] = app.adjustBrightness(color.RGBA{255, 0, 0, 255})
				case 1:
					pic.Pix[y*width+x] = app.adjustBrightness(color.RGBA{0, 255, 0, 255})
				case 2:
					pic.Pix[y*width+x] = app.adjustBrightness(color.RGBA{0, 0, 255, 255})
				}
			}
		}

		sprite := pixel.NewSprite(pic, bounds)
		sprite.Draw(win, pixel.IM.Moved(bounds.Center()))
	}
}

// Utility functions
func (app *App) adjustBrightness(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	f := app.brightness
	return color.RGBA{
		R: uint8(clamp(float64(r>>8)*f, 0, 255)),
		G: uint8(clamp(float64(g>>8)*f, 0, 255)),
		B: uint8(clamp(float64(b>>8)*f, 0, 255)),
		A: uint8(a >> 8),
	}
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

func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
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
			if line != "" {
				lines = append(lines, line)
			}
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
