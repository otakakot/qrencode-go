package qrencode

import (
	"bytes"
)

// The test benchmark shows that encoding with boolBitVector/boolBitGrid is
// twice as fast as byteBitVector/byteBitGrid and uin32BitVector/uint32BitGrid.

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

/*
type BitVector struct {
	byteBitVector
}

type BitGrid struct {
	byteBitGrid
}

func (v *BitVector) AppendBits(b BitVector) {
	v.byteBitVector.AppendBits(b.byteBitVector)
}

func NewBitGrid(width, height int) *BitGrid {
	return &BitGrid{newByteBitGrid(width, height)}
}
*/

/*
type BitVector struct {
	uint32BitVector
}

type BitGrid struct {
	uint32BitGrid
}

func (v *BitVector) AppendBits(b BitVector) {
	v.uint32BitVector.AppendBits(b.uint32BitVector)
}

func NewBitGrid(width, height int) *BitGrid {
	return &BitGrid{newUint32BitGrid(width, height)}
}
*/

func (v *BitVector) String() string {
	b := bytes.Buffer{}
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
	b := bytes.Buffer{}
	for y, w, h := 0, g.Width(), g.Height(); y < h; y++ {
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
func (g *BitGrid) GetRGB565Size(blockSize, margin int) (width, height int) {
	width = blockSize * (2*margin + g.Width())
	height = blockSize * (2*margin + g.Height())
	return
}

// Return an image of the grid, with black blocks for true items and
// white blocks for false items, with the given block size and margin.
func (g *BitGrid) ToRGB565WithMargin(blockSize, margin int) []uint16 {
	width := uint16(blockSize * (2*margin + g.Width()))
	height := uint16(blockSize * (2*margin + g.Height()))
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

	for y := range uint16(g.Height()) {
		for x := range uint16(g.Width()) {
			if g.Get(int(x), int(y)) {
				x0 := uint16(blockSize) * (x + uint16(margin))
				y0 := uint16(blockSize) * (y + uint16(margin))
				for dy := uint16(0); dy < uint16(blockSize); dy++ {
					for dx := uint16(0); dx < uint16(blockSize); dx++ {
						idx := (y0+dy)*width + (x0 + dx)
						if idx >= 0 && idx < uint16(size) {
							pixels[idx] = black
						}
					}
				}
			}
		}
	}

	return pixels
}
