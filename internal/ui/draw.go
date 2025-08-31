package ui

import (
	"fmt"
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"

	"github.com/keshon/screen-tester/internal/core"
	"github.com/keshon/screen-tester/internal/version"
)

var atlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)

func DrawInfo(ctx *core.WindowContext, test core.ScreenTest, opts core.TestOptions, brightness float64) {
	lines := []string{
		fmt.Sprintf("%s", test.Name()),
		fmt.Sprintf("Resolution: %dx%d", ctx.ScreenWidth, ctx.ScreenHeight),
		fmt.Sprintf("Brightness: %.1f", brightness),
	}

	if size, ok := opts.Extra["size"]; ok {
		lines = append(lines, fmt.Sprintf("Size: %d px", size.(int)))
	}

	if direction, ok := opts.Extra["direction"]; ok {
		lines = append(lines, fmt.Sprintf("Direction: %v", direction.(string)))
	}

	if speed, ok := opts.Extra["speed"]; ok {
		lines = append(lines, fmt.Sprintf("Speed: %.0f ms", speed.(time.Duration).Seconds()*1000))
	}

	if test.Description() != "" {
		lines = append(lines, "")
		lines = append(lines, core.WrapText(test.Description(), 60)...)
	}

	lines = append(lines, "")
	lines = append(lines, "Controls:")
	lines = append(lines, "Left / Right: Switch tests")
	lines = append(lines, "F1: Toggle info")
	lines = append(lines, "ESC: Exit")
	lines = append(lines, "")
	lines = append(lines, version.AppFullName)
	lines = append(lines, version.AppDescription)
	lines = append(lines, version.AppRepo)
	lines = append(lines, fmt.Sprintf("Made by %s", version.AppAuthor))

	imd := imdraw.New(nil)
	imd.Color = color.RGBA{0, 0, 0, 255}

	x := 10.0
	y := ctx.Win.Bounds().H() - 10.0
	paddingTop := 30.0
	paddingBottom := 10.0
	paddingSides := 30.0
	lineHeight := 14.0

	var boxWidth float64
	for _, line := range lines {
		txt := text.New(pixel.V(0, 0), atlas)
		fmt.Fprint(txt, line)
		if w := txt.Bounds().W(); w > boxWidth {
			boxWidth = w
		}
	}

	boxHeight := float64(len(lines))*lineHeight + paddingTop + paddingBottom

	imd.Push(pixel.V(x, y-boxHeight), pixel.V(x+boxWidth+paddingSides, y))
	imd.Rectangle(0)
	imd.Draw(ctx.Win)

	for i, line := range lines {
		txt := text.New(pixel.V(x+paddingSides/2, y-paddingTop-lineHeight*float64(i)), atlas)
		txt.Color = colornames.White
		fmt.Fprint(txt, line)
		txt.Draw(ctx.Win, pixel.IM)
	}

}
