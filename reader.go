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
	MaxBlock24 = 0xFFFFFF
	MaxBlock32 = 0xFFFFFFFF
)

var (
	// ErrBlockSizeTooSmall is returned when the block size too small for a reader.
	ErrBlockSizeTooSmall = errors.New("block size too small for the reader")
	// ErrSizeTooLarge is returned when the given size exceed the block size used by a reader.
	ErrSizeTooLarge = errors.New("size too large for the block reader")
)

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

// NewReader24 returns a new reader that is able to read blocks of size MaxBlock24.
func NewReader24(r io.Reader) io.Reader {
	r24 := &reader{
		src:  r,
		size: MaxBlock24,
	}
	buf := make([]byte, 4)

	r24.rsize = func() (int, error) {
		_, err := r24.src.Read(buf[1:]) // 3 bytes because we work on 24bit. We let to zero the fourth byte at index 0 for binary.BigEndian.Uint32's behavior.
		if err != nil {
			return 0, err
		}

		return int(binary.BigEndian.Uint32(buf)), nil
	}

	return r24
}

// NewReader24Custom returns a new reader that is able to read blocks up to size MaxBlock24.
func NewReader24Custom(r io.Reader, size int) (io.Reader, error) {
	if size > MaxBlock24 {
		return nil, ErrSizeTooLarge
	}

	r24c := &reader{
		src:  r,
		size: size,
	}
	buf := make([]byte, 4)

	r24c.rsize = func() (int, error) {
		_, err := r24c.src.Read(buf[1:]) // 3 bytes because we work on 24bit. We let to zero the fourth byte at index 0 for binary.BigEndian.Uint32's behavior.
		if err != nil {
			return 0, err
		}

		return int(binary.BigEndian.Uint32(buf)), nil
	}

	return r24c, nil
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

// NewReader32Custom returns a new reader that is able to read blocks up to size MaxBlock32.
func NewReader32Custom(r io.Reader, size int) (io.Reader, error) {
	if size > MaxBlock32 {
		return nil, ErrSizeTooLarge
	}

	r32c := &reader{
		src:  r,
		size: size,
	}
	buf := make([]byte, 4)

	r32c.rsize = func() (int, error) {
		_, err := r32c.src.Read(buf)
		if err != nil {
			return 0, err
		}

		return int(binary.BigEndian.Uint32(buf)), nil
	}

	return r32c, nil
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
