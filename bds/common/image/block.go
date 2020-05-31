package image

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

const BlockSize = 4 * 1024 * 1024

type Block struct {
	ID         uuid.UUID `json:"id"`
	Size       uint      `json:"size"`
	Offset     uint64    `json:"offset"`
	Created    bool      `json:"created"`
	LastUpdate time.Time `json:"last_update"`

	sync.Mutex `json:"-"`
	Cache      []byte `json:"-"`
}

func (b *Block) Cached(cache []byte) {
	b.Lock()
	defer b.Unlock()
	if cache != nil {
		b.Cache = cache
		return
	}

	b.Cache = make([]byte, b.Size)
}

func NewBlock(size uint, offset uint64) *Block {
	b := new(Block)
	b.Init(size, offset)
	return b
}

func (b *Block) Init(size uint, offset uint64) {
	b.ID = uuid.New()

	b.Size = size
	b.Offset = offset

	b.Created = false
	b.LastUpdate = time.Now()
}

func (b *Block) WriteAt(data []byte, off uint) (n uint64) {
	//defer log.Tracef("Write block: %s over", b.ID.String())

	b.Lock()
	defer b.Unlock()

	if b.Cache == nil {
		return
	}

	log.Tracef("Write block: %s", b.ID.String())

	copy(b.Cache[off:], data)
	b.LastUpdate = time.Now()

	if len(b.Cache[off:]) > len(data) {
		return uint64(len(data))
	} else {
		return uint64(len(b.Cache[off:]))
	}
}

func (b *Block) ReadAt(data []byte, off uint) (n uint64) {
	//defer log.Tracef("Read block: %s over", b.ID.String())

	b.Lock()
	defer b.Unlock()

	if b.Cache == nil {
		return
	}

	log.Tracef("Read block: %s", b.ID.String())

	copy(data, b.Cache[off:])

	if len(b.Cache[off:]) > len(data) {
		return uint64(len(data))
	} else {
		return uint64(len(b.Cache[off:]))
	}
}

func (b *Block) UpdateFrom(block *Block) {
	b.Lock()
	defer b.Unlock()

	b.ID = block.ID

	b.Size = block.Size
	b.Offset = block.Offset

	b.Created = block.Created
	b.LastUpdate = block.LastUpdate
}
