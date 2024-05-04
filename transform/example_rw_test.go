package transform_test

import (
	"bytes"
	"fmt"
	"io"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// NewReader creates readers that encode and decode Latin1 from/to
// UTF-8.
func ExampleNewReader() {
	var (
		utf8  = []byte("tschüß")
		rUTF8 = bytes.NewReader(utf8)

		encLatin1 = charmap.ISO8859_1.NewEncoder()
		decLatin1 = charmap.ISO8859_1.NewDecoder()

		buf1 = &bytes.Buffer{}
		buf2 = &bytes.Buffer{}
	)

	r := transform.NewReader(rUTF8, encLatin1)
	io.Copy(buf1, r)

	latin1 := buf1.Bytes() // copy for printing, before next read empties buf1

	r = transform.NewReader(buf1, decLatin1)
	io.Copy(buf2, r)

	fmt.Println("Original bytes:", utf8)
	fmt.Println("Latin1 bytes:  ", latin1)
	fmt.Println("UTF-8 bytes:   ", buf2.Bytes())
	// Output:
	// Original bytes: [116 115 99 104 195 188 195 159]
	// Latin1 bytes:   [116 115 99 104 252 223]
	// UTF-8 bytes:    [116 115 99 104 195 188 195 159]
}

// NewWriter, along with NewReader, creates a Latin1-to-UTF-16BE
// transcoder.
//
// This example uses an intermediate bytes.Buffer just to be able
// to show the results of the transcoding; a more practical
// example might have the writer wrap a file, like os.Stdout.
func ExampleNewWriter() {
	var (
		latin1  = []byte{116, 115, 99, 104, 252, 223} // tschüß as Latin1 bytes
		rLatin1 = bytes.NewReader(latin1)

		decLatin1  = charmap.ISO8859_1.NewDecoder()
		encUTF16BE = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder()

		buf = &bytes.Buffer{}
	)

	r := transform.NewReader(rLatin1, decLatin1)
	w := transform.NewWriter(buf, encUTF16BE)
	io.Copy(w, r)

	fmt.Println("Latin1 bytes:  ", latin1)
	fmt.Println("UTF-16BE bytes:", buf.Bytes())
	// Output:
	// Latin1 bytes:   [116 115 99 104 252 223]
	// UTF-16BE bytes: [0 116 0 115 0 99 0 104 0 252 0 223]
}
