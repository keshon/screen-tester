package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/keshon/screen-tester/internal/core"
	"github.com/keshon/screen-tester/internal/input"
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

	tests := core.AllTests()

	menuButtons := make([]ui.Button, 0, len(tests)+1)
	for _, t := range tests {
		menuButtons = append(menuButtons, ui.Button{
			Text: t.Name(),
		})
	}
	menuButtons = append(menuButtons, ui.Button{
		Text: "Exit",
	})

	menu := &ui.Menu{
		Buttons: menuButtons,
		Hovered: 0,
	}

	showMenu := true
	currentTest := tests[0]
	testControls := &input.TestInput{}

	cursor := imdraw.New(nil)

	for !win.Closed() {
		ctx.Win.Clear(colornames.Black)

		if showMenu {
			ui.DrawPixelTitle(ctx.Win, "SCREEN TESTER", ctx.Win.Bounds().W(), ctx.Win.Bounds().H(), time.Now())
			ui.DrawTitle(ctx.Win,
				fmt.Sprintf("%s - %s\nMade by %s (%s)",
					version.AppFullName,
					version.AppDescription,
					version.AppAuthor,
					version.AppRepo),
				ctx.Win.Bounds().W(),
				ctx.Win.Bounds().H(),
			)

			ui.LayoutMenuButtons(menu, ctx.Win.Bounds().W(), ctx.Win.Bounds().H())
			ui.DrawMenu(ctx, menu)

			if sel := input.HandleMenuInput(ctx, menu); sel != nil {
				if sel.Text == "Exit" {
					break
				}
				if menu.Hovered >= 0 && menu.Hovered < len(tests) {
					testControls.Current = menu.Hovered
					currentTest = tests[testControls.Current]
					showMenu = false
				}
			}

			cursor.Clear()
			cursor.Color = colornames.White
			pos := win.MousePosition()
			cursor.Push(pos)
			cursor.Circle(6, 2)
			cursor.Draw(win)

		} else {
			currentTest = tests[testControls.Current]
			testControls.HandleTestInput(ctx, tests)

			if ctx.Win.JustPressed(pixelgl.KeyEscape) {
				showMenu = true
				continue
			}

			currentTest.Run(ctx)

			if ctx.ShowInfo {
				ui.DrawInfo(ctx, currentTest, currentTest.Options(), ctx.Brightness)
			}
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
