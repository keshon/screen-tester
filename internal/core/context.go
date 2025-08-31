package core

import (
	"time"

	"github.com/faiface/pixel/pixelgl"
)

type WindowContext struct {
	Win             *pixelgl.Window
	Brightness      float64
	FlickerInterval time.Duration
	ShowInfo        bool
	ScreenWidth     int
	ScreenHeight    int
}
