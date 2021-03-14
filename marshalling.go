package blockio

import (
	"io"
)

////////////////////////////
//                        //
// Decoder                //
//                        //
////////////////////////////

type (
	// Decode parses bytes to an interface.
	Decode func(data []byte, v interface{}) error

	// A Decoder reads blocks and decodes in an object.
	Decoder struct {
		r      io.Reader
		decode Decode
		buf    []byte
	}
)

// NewBlockDecoder decodes values from r using the given h.
// Provided buf must be large enough to handle blocks.
func NewBlockDecoder(r io.Reader, h Decode, buf []byte) *Decoder {
	return &Decoder{
		r:      r,
		decode: h,
		buf:    buf,
	}
}

// NewBlock8Decoder decodes values from r using the given h from Block8.
func NewBlock8Decoder(r io.Reader, h Decode) *Decoder {
	return NewBlockDecoder(NewReader8(r), h, make([]byte, MaxBlock8))
}

// NewBlock16Decoder decodes values from r using the given h from Block16.
func NewBlock16Decoder(r io.Reader, h Decode) *Decoder {
	return NewBlockDecoder(NewReader16(r), h, make([]byte, MaxBlock16))
}

// NewBlock32Decoder decodes values from r using the given h from Block32.
func NewBlock32Decoder(r io.Reader, h Decode) *Decoder {
	return NewBlockDecoder(NewReader32(r), h, make([]byte, MaxBlock32))
}

// Read reads from its block reader and deserialized data in v.
func (d *Decoder) Read(v interface{}) error {
	n, err := d.r.Read(d.buf[:cap(d.buf)])
	if err != nil {
		return err
	}

	return d.decode(d.buf[:n], v)
}

////////////////////////////
//                        //
// Encoder                //
//                        //
////////////////////////////

type (
	// Encode generates serialized bytes from an intrerface.
	Encode func(v interface{}) ([]byte, error)

	// An Encoder encodes objects and writes them as blocks.
	Encoder struct {
		w      io.Writer
		encode Encode
	}
)

// NewBlockEncoder encodes values to w using the given h.
// Provided w must be a block writer.
func NewBlockEncoder(w io.Writer, h Encode) *Encoder {
	return &Encoder{
		w:      w,
		encode: h,
	}
}

// NewBlock8Encoder encodes values to w using the given h in Block8.
func NewBlock8Encoder(w io.Writer, h Encode) *Encoder {
	return NewBlockEncoder(NewWriter8(w), h)
}

// NewBlock16Encoder encodes values to w using the given h in Block16.
func NewBlock16Encoder(w io.Writer, h Encode) *Encoder {
	return NewBlockEncoder(NewWriter16(w), h)
}

// NewBlock32Encoder encodes values to w using the given h in Block32.
func NewBlock32Encoder(w io.Writer, h Encode) *Encoder {
	return NewBlockEncoder(NewWriter32(w), h)
}

// Write writes marshalized bytes to its writer of the given v.
func (e *Encoder) Write(v interface{}) error {
	payload, err := e.encode(v)
	if err != nil {
		return err
	}

	_, err = e.w.Write(payload)
	return err
}
