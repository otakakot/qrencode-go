# qrencode-go

This repository is a fork of [qpliu/qrencode-go](https://github.com/qpliu/qrencode-go) to make it work with TinyGo.

# Example

```go
package main

import (
	"tinygo.org/x/drivers/examples/ili9341/initdisplay"
	"tinygo.org/x/drivers/ili9341"

	"github.com/otakakot/qrencode-go/qrencode"
)

func main() {
	display := initdisplay.InitDisplay()

	width, height := display.Size()
	if width < 320 || height < 240 {
		display.SetRotation(ili9341.Rotation270)
	}

	grid, err := qrencode.Encode("https://github.com/otakakot", qrencode.ECLevelL)
	if err != nil {
		panic(err)
	}

	qrSize := 240
	pixels := grid.ToRGB565WithSize(qrSize, qrSize)

	blockSize := qrSize / (grid.Width() + 8)
	if blockSize < 1 {
		blockSize = 1
	}
	actualWidth, actualHeight := grid.GetRGB565Size(blockSize, 4)

	x := (320 - actualWidth) / 2
	y := (240 - actualHeight) / 2

	if err = display.DrawRGBBitmap(int16(x), int16(y), pixels, int16(actualWidth), int16(actualHeight)); err != nil {
		panic(err)
	}
}
```
