package core

import (
	"image/color"
	"strings"

	"github.com/faiface/pixel/pixelgl"
)

func Clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func AdjustBrightness(c color.RGBA, brightness float64) color.RGBA {
	return color.RGBA{
		R: uint8(Clamp(float64(c.R)*brightness, 0, 255)),
		G: uint8(Clamp(float64(c.G)*brightness, 0, 255)),
		B: uint8(Clamp(float64(c.B)*brightness, 0, 255)),
		A: c.A,
	}
}

func AdjustBrightnessWithKeys(ctx *WindowContext, step float64) {
	if ctx.Win.JustPressed(pixelgl.KeyUp) {
		ctx.Brightness += step
		if ctx.Brightness > 1 {
			ctx.Brightness = 1
		}
	}
	if ctx.Win.JustPressed(pixelgl.KeyDown) {
		ctx.Brightness -= step
		if ctx.Brightness < 0 {
			ctx.Brightness = 0
		}
	}
}

func WrapText(text string, limit int) []string {
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
