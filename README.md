# BlockIO

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/mdouchement/blockio)
[![Go Report Card](https://goreportcard.com/badge/github.com/mdouchement/blockio)](https://goreportcard.com/report/github.com/mdouchement/blockio)
[![License](https://img.shields.io/github/license/mdouchement/blockio.svg)](http://opensource.org/licenses/MIT)


BlockIO is a simple package to write and read a file as binary blocks. It's the same idea as putting several JSON objects line per line in a plaintext file.

**This package is not threadsafe.**

## Format

```
[block_size][block_data][block_size][block_data][block_size][block_data]...
```

Possible `block_size` values:
- Block8 coded on 1 byte for maximum 255-byte block length
- Block16 coded on 2 bytes for maximum 65535-byte block length
- Block24 coded on 3 bytes for maximum 16777215-byte block length
- Block32 coded on 4 bytes for maximum 4294967295-byte block length (take care your memory usage ;)

## Usage

```go
package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fxamacker/cbor/v2"
	"github.com/klauspost/compress/zstd"
	"github.com/mdouchement/blockio"
)

type payload struct {
	FieldA int
	FieldB string
}

func main() {
	var buf bytes.Buffer

	zstdE, err := zstd.NewWriter(nil)
	if err != nil {
		panic(err)
	}

	encoder := blockio.NewBlock16Encoder(&buf, func(v interface{}) ([]byte, error) {
		payload, err := cbor.Marshal(v)
		if err != nil {
			return nil, err
		}

		return zstdE.EncodeAll(payload, nil), nil
	})

	v1 := payload{
		FieldA: 42,
		FieldB: strings.Repeat("test", 42),
	}

	if err := encoder.Write(v1); err != nil {
		panic(err)
	}

	fmt.Println(buf.Bytes())
	fmt.Println(buf.String())

	//
	//

	buf = *bytes.NewBuffer(buf.Bytes())
	var v2 payload

	zstdD, err := zstd.NewReader(nil)
	if err != nil {
		panic(err)
	}

	decoder := blockio.NewBlock16Decoder(&buf, func(data []byte, v interface{}) (err error) {
		data, err = zstdD.DecodeAll(data, nil)
		if err != nil {
			return err
		}

		return cbor.Unmarshal(data, v)
	})

	if err := decoder.Read(&v2); err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", v2)
}
```

## License

**MIT**


## Contributing

All PRs are welcome.

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
5. Push to the branch (git push origin my-new-feature)
6. Create new Pull Request