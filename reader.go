package blockio

import (
	"encoding/binary"
	"errors"
	"io"
)

// Max size of blocks in bytes.
const (
	MaxBlock8  = 0xFF
	MaxBlock16 = 0xFFFF
	MaxBlock32 = 0xFFFFFFFF
)

// ErrBlockSizeTooSmall is returned when the block size too small for a reader.
var ErrBlockSizeTooSmall = errors.New("block size too small for the reader")

type reader struct {
	src   io.Reader
	size  int
	rsize func() (int, error)
}

// NewReader8 returns a new reader that is able to read blocks of size MaxBlock8.
func NewReader8(r io.Reader) io.Reader {
	r8 := &reader{
		src:  r,
		size: MaxBlock8,
	}
	buf := make([]byte, 1)

	r8.rsize = func() (int, error) {
		_, err := r8.src.Read(buf)
		if err != nil {
			return 0, err
		}

		return int(buf[0]), nil
	}

	return r8
}

// NewReader16 returns a new reader that is able to read blocks of size MaxBlock16.
func NewReader16(r io.Reader) io.Reader {
	r16 := &reader{
		src:  r,
		size: MaxBlock16,
	}
	buf := make([]byte, 2)

	r16.rsize = func() (int, error) {
		_, err := r16.src.Read(buf)
		if err != nil {
			return 0, err
		}

		return int(binary.BigEndian.Uint16(buf)), nil
	}

	return r16
}

// NewReader32 returns a new reader that is able to read blocks of size MaxBlock32.
func NewReader32(r io.Reader) io.Reader {
	r32 := &reader{
		src:  r,
		size: MaxBlock32,
	}
	buf := make([]byte, 4)

	r32.rsize = func() (int, error) {
		_, err := r32.src.Read(buf)
		if err != nil {
			return 0, err
		}

		return int(binary.BigEndian.Uint32(buf)), nil
	}

	return r32
}

func (r *reader) Read(p []byte) (n int, err error) {
	if cap(p) < r.size {
		return 0, ErrBlockSizeTooSmall
	}

	n, err = r.rsize()
	if err != nil {
		return 0, err
	}

	return r.src.Read(p[:n])
}
