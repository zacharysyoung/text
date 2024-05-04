// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transform_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"

	"golang.org/x/text/encoding/charmap"
	encunicode "golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// NewReader creates a transformer between UTF-8 and Latin1, and
// vice versa.
func ExampleNewReader() {
	const s = "tschüß"
	bufLatin1 := &bytes.Buffer{}
	bufUTF8 := &bytes.Buffer{}

	fmt.Println("Original bytes:", []byte(s))

	r := transform.NewReader(strings.NewReader(s), charmap.ISO8859_1.NewEncoder())
	io.Copy(bufLatin1, r)

	fmt.Println("Latin1 bytes:  ", bufLatin1.Bytes())

	r = transform.NewReader(bufLatin1, charmap.ISO8859_1.NewDecoder())
	io.Copy(bufUTF8, r)

	fmt.Println("UTF-8 bytes:   ", bufUTF8.Bytes())
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
// example might have the writer wrap os.Stdout.
func ExampleNewWriter() {
	bLatin1 := []byte{116, 115, 99, 104, 252, 223} // tschüß

	buf := &bytes.Buffer{} // intermediate buffer to inspect results

	r := transform.NewReader(bytes.NewReader(bLatin1), charmap.ISO8859_1.NewDecoder())
	w := transform.NewWriter(buf, encunicode.UTF16(encunicode.BigEndian, encunicode.IgnoreBOM).NewEncoder())
	io.Copy(w, r)

	fmt.Println("Latin1 bytes:  ", bLatin1)
	fmt.Println("UTF-16BE bytes:", buf.Bytes())
	// Output:
	// Latin1 bytes:   [116 115 99 104 252 223]
	// UTF-16BE bytes: [0 116 0 115 0 99 0 104 0 252 0 223]
}

func ExampleRemoveFunc() {
	input := []byte(`tschüß; до свидания`)

	b := make([]byte, len(input))

	t := transform.RemoveFunc(unicode.IsSpace)
	n, _, _ := t.Transform(b, input, true)
	fmt.Println(string(b[:n]))

	t = transform.RemoveFunc(func(r rune) bool {
		return !unicode.Is(unicode.Latin, r)
	})
	n, _, _ = t.Transform(b, input, true)
	fmt.Println(string(b[:n]))

	n, _, _ = t.Transform(b, norm.NFD.Bytes(input), true)
	fmt.Println(string(b[:n]))

	// Output:
	// tschüß;досвидания
	// tschüß
	// tschuß
}
