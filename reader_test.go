package blockio_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/mdouchement/blockio"
	"github.com/stretchr/testify/assert"
)

func TestReader8_Read(t *testing.T) {
	buf := bytes.NewBuffer([]byte{4, 'd', 'a', 't', 'a', 5, 'd', 'a', 't', 'u', 'm'})
	r := blockio.NewReader8(buf)

	//
	// Too small block

	n, err := r.Read(make([]byte, blockio.MaxBlock8-1))
	assert.ErrorAs(t, err, &blockio.ErrBlockSizeTooSmall)
	assert.Equal(t, 0, n)

	//
	// Read `data`

	block := make([]byte, blockio.MaxBlock8)
	n, err = r.Read(block)
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, []byte("data"), block[:n])

	//
	// Read remaining `datum`

	block = block[:blockio.MaxBlock8]
	n, err = r.Read(block)
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte("datum"), block[:n])

	//
	// EOF

	block = block[:blockio.MaxBlock8]
	n, err = r.Read(block)
	assert.ErrorAs(t, err, &io.EOF)
	assert.Equal(t, 0, n)
}

func TestReader16_Read(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0, 4, 'd', 'a', 't', 'a', 0, 5, 'd', 'a', 't', 'u', 'm'})
	r := blockio.NewReader16(buf)

	//
	// Too small block

	n, err := r.Read(make([]byte, blockio.MaxBlock16-1))
	assert.ErrorAs(t, err, &blockio.ErrBlockSizeTooSmall)
	assert.Equal(t, 0, n)

	//
	// Read `data`

	block := make([]byte, blockio.MaxBlock16)
	n, err = r.Read(block)
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, []byte("data"), block[:n])

	//
	// Read remaining `datum`

	block = block[:blockio.MaxBlock16]
	n, err = r.Read(block)
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte("datum"), block[:n])

	//
	// EOF

	block = block[:blockio.MaxBlock16]
	n, err = r.Read(block)
	assert.ErrorAs(t, err, &io.EOF)
	assert.Equal(t, 0, n)
}

func TestReader32_Read(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0, 0, 0, 4, 'd', 'a', 't', 'a', 0, 0, 0, 5, 'd', 'a', 't', 'u', 'm'})
	r := blockio.NewReader32(buf)

	//
	// Too small block

	n, err := r.Read(make([]byte, blockio.MaxBlock32-1))
	assert.ErrorAs(t, err, &blockio.ErrBlockSizeTooSmall)
	assert.Equal(t, 0, n)

	//
	// Read `data`

	block := make([]byte, blockio.MaxBlock32)
	n, err = r.Read(block)
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, []byte("data"), block[:n])

	//
	// Read remaining `datum`

	block = block[:blockio.MaxBlock32]
	n, err = r.Read(block)
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte("datum"), block[:n])

	//
	// EOF

	block = block[:blockio.MaxBlock32]
	n, err = r.Read(block)
	assert.ErrorAs(t, err, &io.EOF)
	assert.Equal(t, 0, n)
}
