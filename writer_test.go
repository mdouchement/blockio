package blockio_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/mdouchement/blockio"
	"github.com/stretchr/testify/assert"
)

func TestWriter8_Write(t *testing.T) {
	var buf bytes.Buffer
	w := blockio.NewWriter8(&buf)

	//
	// Write `data`

	data := []byte("data")
	n, err := w.Write(data)
	assert.NoError(t, err)

	ld := len(data)
	assert.Equal(t, ld+1, n)

	expected := append([]byte{byte(ld)}, data...)
	assert.Equal(t, expected, buf.Bytes())

	//
	// Append another `data` to buffer

	n, err = w.Write(data)
	assert.NoError(t, err)

	assert.Equal(t, ld+1, n)

	expected = append(expected, expected...)
	assert.Equal(t, expected, buf.Bytes())

	//
	// Just before block size limit

	buf.Reset()

	data = bytes.Repeat([]byte{'b'}, blockio.MaxBlock8)
	n, err = w.Write(data)
	assert.NoError(t, err)

	ld = len(data)
	assert.Equal(t, ld+1, n)

	expected = append([]byte{byte(ld)}, data...)
	assert.Equal(t, expected, buf.Bytes())

	//
	// Block out of limit

	data = bytes.Repeat([]byte{'b'}, blockio.MaxBlock8+1)
	n, err = w.Write(data)
	assert.ErrorAs(t, err, &blockio.ErrBlockSize)
	assert.Equal(t, 0, n)
}

func TestWriter16_Write(t *testing.T) {
	var buf bytes.Buffer
	w := blockio.NewWriter16(&buf)

	bufw := make([]byte, blockio.MaxBlock16+2)

	//
	// Write `data`

	data := []byte("data")
	n, err := w.Write(data)
	assert.NoError(t, err)

	ld := len(data)
	assert.Equal(t, ld+2, n)

	binary.BigEndian.PutUint16(bufw[:2], uint16(ld))
	copy(bufw[2:], data)
	expected := bufw[:ld+2]
	assert.Equal(t, expected, buf.Bytes())

	//
	// Append another `data` to buffer

	n, err = w.Write(data)
	assert.NoError(t, err)

	assert.Equal(t, ld+2, n)

	expected = append(expected, expected...)
	assert.Equal(t, expected, buf.Bytes())

	//
	// Just before block size limit

	buf.Reset()

	data = bytes.Repeat([]byte{'b'}, blockio.MaxBlock16)
	n, err = w.Write(data)
	assert.NoError(t, err)

	ld = len(data)
	assert.Equal(t, ld+2, n)

	binary.BigEndian.PutUint16(bufw[:2], uint16(ld))
	copy(bufw[2:], data)
	assert.Equal(t, bufw[:ld+2], buf.Bytes())

	//
	// Block out of limit

	data = bytes.Repeat([]byte{'b'}, blockio.MaxBlock16+1)
	n, err = w.Write(data)
	assert.ErrorAs(t, err, &blockio.ErrBlockSize)
	assert.Equal(t, 0, n)
}

func TestWriter32_Write(t *testing.T) {
	var buf bytes.Buffer
	w := blockio.NewWriter32(&buf)

	bufw := make([]byte, blockio.MaxBlock32+4)

	//
	// Write `data`

	data := []byte("data")
	n, err := w.Write(data)
	assert.NoError(t, err)

	ld := len(data)
	assert.Equal(t, ld+4, n)

	binary.BigEndian.PutUint32(bufw[:4], uint32(ld))
	copy(bufw[4:], data)
	expected := bufw[:ld+4]
	assert.Equal(t, expected, buf.Bytes())

	//
	// Append another `data` to buffer

	n, err = w.Write(data)
	assert.NoError(t, err)

	assert.Equal(t, ld+4, n)

	expected = append(expected, expected...)
	assert.Equal(t, expected, buf.Bytes())

	// // NEED A LOT OF MEMORY!
	// // Just before block size limit

	// buf.Reset()

	// data = bytes.Repeat([]byte{'b'}, blockio.MaxBlock32)
	// n, err = w.Write(data)
	// assert.NoError(t, err)

	// ld = len(data)
	// assert.Equal(t, ld+4, n)

	// binary.BigEndian.PutUint16(bufw[:4], uint16(ld))
	// copy(bufw[4:], data)
	// assert.Equal(t, bufw[:ld+4], buf.Bytes())

	// //
	// // Block out of limit

	// data = bytes.Repeat([]byte{'b'}, blockio.MaxBlock32+1)
	// n, err = w.Write(data)
	// assert.ErrorAs(t, err, &blockio.ErrBlockSize)
	// assert.Equal(t, 0, n)
}
