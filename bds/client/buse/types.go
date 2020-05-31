/*
source https://github.com/samalba/buse-go

MIT License

Copyright (c) 2017 Sam Alba

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package buse

import (
	"os"
)

// Rewrote type definitions for #defines and structs to workaround cgo
// as defined in <linux/nbd.h>

const (
	NBD_SET_SOCK        = (0xab<<8 | 0)
	NBD_SET_BLKSIZE     = (0xab<<8 | 1)
	NBD_SET_SIZE        = (0xab<<8 | 2)
	NBD_DO_IT           = (0xab<<8 | 3)
	NBD_CLEAR_SOCK      = (0xab<<8 | 4)
	NBD_CLEAR_QUE       = (0xab<<8 | 5)
	NBD_PRINT_DEBUG     = (0xab<<8 | 6)
	NBD_SET_SIZE_BLOCKS = (0xab<<8 | 7)
	NBD_DISCONNECT      = (0xab<<8 | 8)
	NBD_SET_TIMEOUT     = (0xab<<8 | 9)
	NBD_SET_FLAGS       = (0xab<<8 | 10)
)

const (
	NBD_CMD_READ  = 0
	NBD_CMD_WRITE = 1
	NBD_CMD_DISC  = 2
	NBD_CMD_FLUSH = 3
	NBD_CMD_TRIM  = 4
)

const (
	NBD_FLAG_HAS_FLAGS  = (1 << 0)
	NBD_FLAG_READ_ONLY  = (1 << 1)
	NBD_FLAG_SEND_FLUSH = (1 << 2)
	NBD_FLAG_SEND_TRIM  = (1 << 5)
)

const (
	NBD_REQUEST_MAGIC = 0x25609513
	NBD_REPLY_MAGIC   = 0x67446698
)

type nbdRequest struct {
	Magic  uint32
	Type   uint32
	Handle uint64
	From   uint64
	Length uint32
}

type nbdReply struct {
	Magic  uint32
	Error  uint32
	Handle uint64
}

type BuseInterface interface {
	ReadAt(p []byte, off uint64) error
	WriteAt(p []byte, off uint64) error
	Disconnect()
	Flush() error
	Trim(off uint64, length uint) error
}

type BuseDevice struct {
	size       uint64
	device     string
	driver     BuseInterface
	deviceFp   *os.File
	socketPair [2]int
	op         [5]func(driver BuseInterface, fp *os.File, chunk []byte, request *nbdRequest, reply *nbdReply) error
	disconnect chan int
}
