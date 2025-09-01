package input

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/keshon/screen-tester/internal/core"
	"github.com/keshon/screen-tester/internal/ui"
)

func HandleMenuInput(ctx *core.WindowContext, menu *ui.Menu) (selected *ui.Button) {
	win := ctx.Win
	mousePos := win.MousePosition()

	for i, btn := range menu.Buttons {
		if mousePos.X >= btn.Bounds.Min.X && mousePos.X <= btn.Bounds.Max.X &&
			mousePos.Y >= btn.Bounds.Min.Y && mousePos.Y <= btn.Bounds.Max.Y {
			menu.Hovered = i
		}
	}

	if win.JustPressed(pixelgl.MouseButtonLeft) {
		selected = &menu.Buttons[menu.Hovered]
	}

	if win.JustPressed(pixelgl.KeyDown) {
		menu.Hovered++
		if menu.Hovered >= len(menu.Buttons) {
			menu.Hovered = 0
		}
	}
	if win.JustPressed(pixelgl.KeyUp) {
		menu.Hovered--
		if menu.Hovered < 0 {
			menu.Hovered = len(menu.Buttons) - 1
		}
	}
	if win.JustPressed(pixelgl.KeyEnter) {
		selected = &menu.Buttons[menu.Hovered]
	}
	return
}
