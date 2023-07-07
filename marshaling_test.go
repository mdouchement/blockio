package blockio_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/mdouchement/blockio"
	"github.com/stretchr/testify/assert"
)

func TestEncoder_Write(t *testing.T) {
	v := struct {
		Field1 string
		Field2 int
	}{
		Field1: "test",
		Field2: 42,
	}

	var buf bytes.Buffer

	//

	encoder := blockio.NewBlock8Encoder(&buf, json.Marshal)
	err := encoder.Write(v)
	assert.NoError(t, err)
	assert.Equal(t, "\x1d{\"Field1\":\"test\",\"Field2\":42}", buf.String())

	//

	buf.Reset()
	encoder = blockio.NewBlock16Encoder(&buf, json.Marshal)
	err = encoder.Write(v)
	assert.NoError(t, err)
	assert.Equal(t, "\x00\x1d{\"Field1\":\"test\",\"Field2\":42}", buf.String())

	//

	buf.Reset()
	encoder = blockio.NewBlock24Encoder(&buf, json.Marshal)
	err = encoder.Write(v)
	assert.NoError(t, err)
	assert.Equal(t, "\x00\x00\x1d{\"Field1\":\"test\",\"Field2\":42}", buf.String())

	//

	buf.Reset()
	encoder = blockio.NewBlock32Encoder(&buf, json.Marshal)
	err = encoder.Write(v)
	assert.NoError(t, err)
	assert.Equal(t, "\x00\x00\x00\x1d{\"Field1\":\"test\",\"Field2\":42}", buf.String())
}

func TestDecoder_Read(t *testing.T) {
	type data struct {
		Field1 string
		Field2 int
	}

	//

	buf := bytes.NewBufferString("\x1d{\"Field1\":\"test\",\"Field2\":42}")
	decoder := blockio.NewBlock8Decoder(buf, json.Unmarshal)
	v := data{}
	err := decoder.Read(&v)
	assert.NoError(t, err)
	assert.Equal(t, v.Field1, "test")
	assert.Equal(t, v.Field2, 42)

	//

	buf = bytes.NewBufferString("\x00\x1d{\"Field1\":\"test\",\"Field2\":42}")
	decoder = blockio.NewBlock16Decoder(buf, json.Unmarshal)
	v = data{}
	err = decoder.Read(&v)
	assert.NoError(t, err)
	assert.Equal(t, v.Field1, "test")
	assert.Equal(t, v.Field2, 42)

	//

	buf = bytes.NewBufferString("\x00\x00\x1d{\"Field1\":\"test\",\"Field2\":42}")
	decoder = blockio.NewBlock24Decoder(buf, json.Unmarshal)
	v = data{}
	err = decoder.Read(&v)
	assert.NoError(t, err)
	assert.Equal(t, v.Field1, "test")
	assert.Equal(t, v.Field2, 42)

	//

	buf = bytes.NewBufferString("\x00\x00\x1d{\"Field1\":\"test\",\"Field2\":42}")
	decoder, err = blockio.NewBlock24CustomDecoder(buf, 0x1D, json.Unmarshal)
	assert.NoError(t, err)
	v = data{}
	err = decoder.Read(&v)
	assert.NoError(t, err)
	assert.Equal(t, v.Field1, "test")
	assert.Equal(t, v.Field2, 42)

	//

	buf = bytes.NewBufferString("\x00\x00\x00\x1d{\"Field1\":\"test\",\"Field2\":42}")
	decoder = blockio.NewBlock32Decoder(buf, json.Unmarshal)
	v = data{}
	err = decoder.Read(&v)
	assert.NoError(t, err)
	assert.Equal(t, v.Field1, "test")
	assert.Equal(t, v.Field2, 42)

	//

	buf = bytes.NewBufferString("\x00\x00\x00\x1d{\"Field1\":\"test\",\"Field2\":42}")
	decoder, err = blockio.NewBlock32CustomDecoder(buf, 0x1D, json.Unmarshal)
	assert.NoError(t, err)
	v = data{}
	err = decoder.Read(&v)
	assert.NoError(t, err)
	assert.Equal(t, v.Field1, "test")
	assert.Equal(t, v.Field2, 42)
}
