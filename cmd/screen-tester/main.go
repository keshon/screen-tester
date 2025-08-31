package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/keshon/screen-tester/internal/core"
	_ "github.com/keshon/screen-tester/internal/tests" // auto-register tests
	"github.com/keshon/screen-tester/internal/ui"
	"github.com/keshon/screen-tester/internal/version"
)

func run() {
	monitor := pixelgl.PrimaryMonitor()
	width, height := monitor.Size()

	cfg := pixelgl.WindowConfig{
		Title:       version.AppFullName + " " + version.GoVersion + " " + version.BuildDate,
		Bounds:      pixel.R(0, 0, width, height),
		Monitor:     monitor,
		Undecorated: true,
		Maximized:   true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	ctx := &core.WindowContext{
		Win:          win,
		ScreenWidth:  int(width),
		ScreenHeight: int(height),
		ShowInfo:     true,
		Brightness:   1.0,
	}

	controls := &ui.Controls{}
	tests := core.All()

	for !win.Closed() {
		currentTest := tests[controls.Current]
		controls.HandleInput(ctx, tests)
		currentTest.Run(ctx)
		if ctx.ShowInfo {
			ui.DrawInfo(ctx, currentTest, currentTest.Options(), ctx.Brightness)
		}
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
