package qrencode

import (
	"strings"
)

type BitVector struct {
	boolBitVector
}

type BitGrid struct {
	boolBitGrid
}

func (v *BitVector) AppendBits(b BitVector) {
	v.boolBitVector.AppendBits(b.boolBitVector)
}

func NewBitGrid(width, height int) *BitGrid {
	return &BitGrid{newBoolBitGrid(width, height)}
}

func (v *BitVector) String() string {
	var b strings.Builder
	b.Grow(v.Length())
	for i, l := 0, v.Length(); i < l; i++ {
		if v.Get(i) {
			b.WriteString("1")
		} else {
			b.WriteString("0")
		}
	}
	return b.String()
}

func (g *BitGrid) String() string {
	w, h := g.Width(), g.Height()
	var b strings.Builder
	b.Grow((w + 1) * h)
	for y := range h {
		for x := range w {
			if g.Empty(x, y) {
				b.WriteString(" ")
			} else if g.Get(x, y) {
				b.WriteString("#")
			} else {
				b.WriteString("_")
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

// ToRGB565WithSize returns RGB565 pixel data that fits within the specified width and height.
// It automatically calculates the appropriate block size to fit the QR code.
// The actual size may be smaller than specified to maintain square pixels.
func (g *BitGrid) ToRGB565WithSize(maxWidth, maxHeight int) []uint16 {
	margin := 4

	gridWithMargin := g.Width() + 2*margin

	blockSizeWidth := maxWidth / gridWithMargin

	blockSizeHeight := maxHeight / gridWithMargin

	blockSize := blockSizeWidth

	blockSize = min(blockSizeWidth, blockSizeHeight)

	blockSize = max(blockSize, 1)

	return g.ToRGB565WithMargin(blockSize, margin)
}

// Return an image of the grid, with black blocks for true items and
// white blocks for false items, with the given block size and a
// default margin.
func (g *BitGrid) ToRGB565(blockSize int) []uint16 {
	return g.ToRGB565WithMargin(blockSize, 4)
}

// GetRGB565Size returns the actual width and height of the RGB565 image
// that will be generated with the given block size and margin.
func (g *BitGrid) GetRGB565Size(blockSize int, margin int) (width, height int) {
	width = blockSize * (2*margin + g.Width())
	height = blockSize * (2*margin + g.Height())
	return
}

// Return an image of the grid, with black blocks for true items and
// white blocks for false items, with the given block size and margin.
func (g *BitGrid) ToRGB565WithMargin(blockSize int, margin int) []uint16 {
	gridWidth := g.Width()
	gridHeight := g.Height()
	width := uint16(blockSize * (2*margin + gridWidth))
	height := uint16(blockSize * (2*margin + gridHeight))
	size := int(width * height)

	if size <= 0 || size > 1024*1024 {
		return nil
	}

	pixels := make([]uint16, size)

	white := uint16(0xFFFF) // RGB565: 11111 111111 11111
	black := uint16(0x0000) // RGB565: 00000 000000 00000

	for i := range pixels {
		pixels[i] = white
	}

	blockSizeU16 := uint16(blockSize)
	marginU16 := uint16(margin)

	for y := range gridHeight {
		y0 := blockSizeU16 * (uint16(y) + marginU16)
		for x := range gridWidth {
			if g.Get(x, y) {
				x0 := blockSizeU16 * (uint16(x) + marginU16)
				for dy := uint16(0); dy < blockSizeU16; dy++ {
					rowStart := (y0 + dy) * width
					for dx := uint16(0); dx < blockSizeU16; dx++ {
						pixels[rowStart+x0+dx] = black
					}
				}
			}
		}
	}

	return pixels
}
