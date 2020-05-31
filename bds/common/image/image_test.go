package image

import (
	"testing"
)

func TestNewImage(t *testing.T) {
	img := NewImage("test", 1024*1024*512)

	offset := uint64(0)
	for i, b := range img.Blocks {
		if b.Offset != offset {
			t.Fatal(img.Blocks[i-1].Offset,
				img.Blocks[i-1].Size,
				b.Offset,
				b.Size)
		}
		offset += uint64(b.Size)
	}
}
