package ui

import (
	"github.com/faiface/pixel"
	"github.com/keshon/screen-tester/internal/core"
)

type Menu struct {
	Buttons []Button
	Hovered int
}

func LayoutMenuButtons(menu *Menu, winW, winH float64) {
	btnWidth := 300.0
	btnHeight := 30.0
	padding := 10.0

	totalHeight := float64(len(menu.Buttons))*(btnHeight+padding) - padding
	startY := (winH+totalHeight)/2 - btnHeight

	for i := range menu.Buttons {
		x := (winW - btnWidth) / 2
		y := startY - float64(i)*(btnHeight+padding)
		menu.Buttons[i].Bounds = pixel.R(x, y, x+btnWidth, y+btnHeight)
	}
}

func DrawMenu(ctx *core.WindowContext, menu *Menu) {
	for i, btn := range menu.Buttons {
		hovered := i == menu.Hovered
		DrawButton(ctx, btn, hovered)
	}
}
