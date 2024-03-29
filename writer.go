package blockio

import (
	"encoding/binary"
	"errors"
	"io"
)

// ErrBlockSize is returned when the block size to be written exceeds the writer capabilities.
var ErrBlockSize = errors.New("block size exceed writer capabilities")

type writer struct {
	dst   io.Writer
	buf   []byte
	wsize func(l int) (int, error)
}

// NewWriter8 returns a new writer that is able to write blocks of size up to MaxBlock8.
func NewWriter8(w io.Writer) io.Writer {
	w8 := &writer{
		dst: w,
		buf: make([]byte, MaxBlock8+1),
	}
	w8.wsize = func(l int) (int, error) {
		if l > MaxBlock8 {
			return 0, ErrBlockSize
		}

		w8.buf[0] = byte(uint8(l))
		return 1, nil
	}

	return w8
}

// NewWriter16 returns a new writer that is able to write blocks of size up to MaxBlock16.
func NewWriter16(w io.Writer) io.Writer {
	w16 := &writer{
		dst: w,
		buf: make([]byte, MaxBlock16+2),
	}
	w16.wsize = func(l int) (int, error) {
		if l > MaxBlock16 {
			return 0, ErrBlockSize
		}

		binary.BigEndian.PutUint16(w16.buf[:2], uint16(l))
		return 2, nil
	}

	return w16
}

// NewWriter24 returns a new writer that is able to write blocks of size up to MaxBlock24.
func NewWriter24(w io.Writer) io.Writer {
	w24 := &writer{
		dst: w,
		buf: make([]byte, MaxBlock24+3),
	}
	w24.wsize = func(l int) (int, error) {
		if l > MaxBlock24 {
			return 0, ErrBlockSize
		}

		binary.BigEndian.PutUint32(w24.buf[:4], uint32(l))
		w24.buf[0], w24.buf[1], w24.buf[2] = w24.buf[1], w24.buf[2], w24.buf[3] // Translate over 3 bytes because we work on 24bit.
		return 3, nil
	}

	return w24
}

// NewWriter32 returns a new writer that is able to write blocks of size up to MaxBlock32.
func NewWriter32(w io.Writer) io.Writer {
	w32 := &writer{
		dst: w,
		buf: make([]byte, MaxBlock32+4),
	}
	w32.wsize = func(l int) (int, error) {
		if l > MaxBlock32 {
			return 0, ErrBlockSize
		}

		binary.BigEndian.PutUint32(w32.buf[:4], uint32(l))
		return 4, nil
	}

	return w32
}

func (w *writer) Write(block []byte) (n int, err error) {
	n, err = w.wsize(len(block))
	if err != nil {
		return 0, err
	}

	n += copy(w.buf[n:], block)
	return w.dst.Write(w.buf[:n])
}
