package image

import "sync"

type Image struct {
	sync.Mutex `json:"-"`
	Name       string   `json:"name"`
	Size       uint64   `json:"size"`
	Blocks     []*Block `json:"blocks"`
}

func EmptyImage() *Image {
	return &Image{}
}

func NewImage(name string, size uint64) *Image {
	image := new(Image)
	image.Name = name
	image.Size = size

	blockCount := size/BlockSize + 1

	image.Blocks = make([]*Block, blockCount)

	for i := range image.Blocks {
		image.Blocks[i] = NewBlock(BlockSize, uint64(i*BlockSize))
	}

	return image
}

//func (i *Image) Delete() {
//	b := i.blocks
//	for b != nil {
//		_, _ = req.Delete(b.Url, nil)
//	}
//}

func (i *Image) BlockAt(off uint64) (block *Block, pos uint) {
	index := off / BlockSize
	block = i.Blocks[index]
	pos = uint(off % BlockSize)
	return
}

//
//func (i *Image) WriteAt(p []byte, off int64) (n int, err error) {
//
//	poff := uint64(0)
//	block, _ := i.findBlock(uint64(off))
//	imageOff := off
//
//	for poff < uint64(len(p)) {
//
//		blockOff := imageOff % BlockSize
//		wc := uint64(BlockSize - (imageOff % BlockSize)) //本块可以写入的字节数
//
//		// 检查偏移是否超过p本身的长度
//		if poff+wc > uint64(len(p)) {
//			wc = uint64(len(p)) - poff
//		}
//
//		block.WriteAt(p[poff:poff+wc], uint64(blockOff))
//
//		block = block.next
//		poff += wc
//		imageOff += int64(wc)
//	}
//
//	return 0, nil
//}
//
//func (i *Image) ReadAt(p []byte, off int64) (n int, err error) {
//
//	poff := 0
//	block, _ := i.BlockAt(uint64(off))
//	imageOff := off
//
//	for poff < len(p) {
//
//		blockOff := uint64(imageOff % BlockSize)
//		wc := int(BlockSize - (imageOff % BlockSize)) //本块可以写入的字节数
//
//		// 检查偏移是否超过p本身的长度
//		if poff+wc > len(p) {
//			wc = len(p) - poff
//		}
//
//		block.ReadAt(p[poff:poff+wc], blockOff)
//
//		block = block.next
//		poff += wc
//		imageOff += int64(wc)
//	}
//
//	return 0, nil
//}
//
//func (i *Image) Close() error {
//	i.Flush()
//	return nil
//}
//
//func (i *Image) Flush() {
//	b := i.blocks
//	for b != nil {
//		b.PutData()
//		b.Cache = nil
//	}
//}
