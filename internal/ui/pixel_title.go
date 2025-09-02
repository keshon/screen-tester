package ui

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var fontMap = map[rune][][]int{
	'S': {
		{1, 1, 1},
		{1, 0, 0},
		{1, 1, 1},
		{0, 0, 1},
		{1, 1, 1},
	},
	'C': {
		{1, 1, 1},
		{1, 0, 0},
		{1, 0, 0},
		{1, 0, 0},
		{1, 1, 1},
	},
	'R': {
		{1, 1, 0},
		{1, 0, 1},
		{1, 1, 0},
		{1, 0, 1},
		{1, 0, 1},
	},
	'E': {
		{1, 1, 1},
		{1, 0, 0},
		{1, 1, 0},
		{1, 0, 0},
		{1, 1, 1},
	},
	'N': {
		{1, 0, 1},
		{1, 1, 1},
		{1, 1, 1},
		{1, 1, 1},
		{1, 0, 1},
	},
	'T': {
		{1, 1, 1},
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 0},
	},
	' ': {
		{0},
	},
}

func DrawPixelTitle(win *pixelgl.Window, title string, screenWidth, screenHeight float64, t time.Time) {
	imd := imdraw.New(nil)
	imd.Color = colornames.White

	dotSize := 8.0
	spacing := 5.0
	scale := dotSize + spacing

	totalWidth := 0.0
	for _, ch := range title {
		matrix, ok := fontMap[ch]
		if !ok {
			totalWidth += scale * 4
			continue
		}
		totalWidth += float64(len(matrix[0]))*scale + scale
	}

	x := (screenWidth - totalWidth) / 2
	y := screenHeight - 100.0

	for _, ch := range title {
		matrix, ok := fontMap[ch]
		if !ok {
			x += scale * 4
			continue
		}
		for row := 0; row < len(matrix); row++ {
			for col := 0; col < len(matrix[row]); col++ {
				if matrix[row][col] == 1 {
					cx := x + float64(col)*scale
					cy := y - float64(row)*scale
					imd.Push(
						pixel.V(cx-dotSize/2, cy-dotSize/2),
						pixel.V(cx+dotSize/2, cy+dotSize/2),
					)
					imd.Rectangle(0)
				}
			}
		}
		x += float64(len(matrix[0]))*scale + scale
	}
	imd.Draw(win)
}
