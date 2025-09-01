# Monitor test tool for Linux, macOS and Windows

![Screen Tester Banner](./assets/banner.png)

A simple fullscreen monitor testing tool written in Go using the [Pixel](https://github.com/faiface/pixel) library.  
It provides several test screens for checking colors, gradients, dead pixels, motion clarity, and subpixel layouts.

---

## Available tests
* **`Red`**  
  Solid red screen

* **`Green`**  
  Solid green screen

* **`Blue`**  
  Solid blue screen

* **`White`**  
  Solid white screen

* **`Black`**  
  Solid black screen

* **`Horizontal Gradient`**  
  Black to white gradient (Shift+Up/Down to invert)

* **`Vertical Gradient`**  
  Black to white gradient (Shift+Up/Down to invert)

* **`Small Checkerboard`**  
  Black & white checkerboard with adjustable square size (Shift+Up/Down)

* **`Pixel Grid`**  
  Grid overlay with adjustable cells size (Shift+Up/Down)

* **`Motion Balls`**  
  Bouncing balls with background cycling (Up/Down: background, Shift+Up/Down: speed)

* **`Dead Pixel Recovery`**  
  Flashes colors to exercise dead pixels (Shift+Up/Down to adjust speed)


---

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
