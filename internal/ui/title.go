package ui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

func DrawTitle(win *pixelgl.Window, title string, screenWidth, screenHeight float64) {
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	txt := text.New(pixel.V(0, 0), atlas)

	txt.Color = colornames.White
	txt.Clear()
	txt.WriteString(title)

	margin := 20.0
	pos := pixel.V(screenWidth-txt.Bounds().W()-margin, margin)
	mat := pixel.IM.Moved(pos)

	txt.Draw(win, mat)
}
