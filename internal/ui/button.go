package ui

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"

	"github.com/keshon/screen-tester/internal/core"
)

type Button struct {
	Text   string
	Bounds pixel.Rect
}

func DrawButton(ctx *core.WindowContext, btn Button, hovered bool) {
	imd := imdraw.New(nil)
	col := colornames.Gray
	if hovered {
		col = colornames.Red
	}
	imd.Color = col
	imd.Push(btn.Bounds.Min, btn.Bounds.Max)
	imd.Rectangle(0)
	imd.Draw(ctx.Win)

	// текст
	txt := text.New(pixel.V(btn.Bounds.Min.X+10, btn.Bounds.Min.Y+10), Atlas)
	txt.Color = colornames.White
	fmt.Fprint(txt, btn.Text)
	txt.Draw(ctx.Win, pixel.IM)
}
