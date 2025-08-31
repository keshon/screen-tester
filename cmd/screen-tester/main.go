package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"app/internal/core"
	_ "app/internal/tests" // auto-register tests
	"app/internal/ui"
)

func run() {
	monitor := pixelgl.PrimaryMonitor()
	width, height := monitor.Size()

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
