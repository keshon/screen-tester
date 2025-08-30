# Screen Tester - Monitor test tool

A simple fullscreen monitor testing tool written in Go using the [Pixel](https://github.com/faiface/pixel) library.  
It provides several test screens for checking colors, gradients, dead pixels, motion clarity, and subpixel layouts.

## Features

- Solid color fills (Red, Green, Blue, White, Black, Gray)
- Horizontal and vertical gradients
- Checkerboard and pixel grid patterns
- Dead pixel recovery flicker
- Motion test with bouncing balls
- Subpixel rendering test

## Controls

- **Left / Right Arrow**: Switch between tests  
- **F1**: Toggle info overlay  
- **Esc**: Exit the program  

Per-test controls:

- **Brightness tests**: Adjust brightness with **Up / Down**  
- **Speed tests**: Adjust flicker or motion speed with **Up / Down**  
- **Motion / Subpixel tests**: Use **Shift + Up/Down** for background or channel adjustments  

## Requirements

- Go 1.18+  
- Pixel library dependencies:
  - `github.com/faiface/pixel`
  - `github.com/faiface/pixel/imdraw`
  - `golang.org/x/image`

## Build and Run

```bash
git clone https://github.com/keshon/screen-tester.git
cd screen-tester
go run .
```

or in Windows

```bash
build-n-run.bat
```