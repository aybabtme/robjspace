package rubyobj_test

import (
	"bytes"
	"github.com/aybabtme/rubyobj"
	"io"
	"io/ioutil"
	"testing"
)

func BenchmarkDecode_TinyDump(b *testing.B)  { decode(b, "testdata/tiny.json") }
func BenchmarkDecode_SmallDump(b *testing.B) { decode(b, "testdata/small.json") }

func decode(b *testing.B, filename string) {
	r := jsonReader(b, filename)

	var err error

	b.ResetTimer()
	for n := 0; n < b.N; n++ {

		dec := rubyobj.NewDecoder(r)
		for err != io.EOF {
			rObj := rubyobj.RubyObject{}
			err = dec.Decode(&rObj)
			if err != nil && err != io.EOF {
				b.Fatal(err)
			}
		}

		r.Reset()
	}
}

func BenchmarkEncode_TinyDump(b *testing.B)  { encode(b, "testdata/tiny.json") }
func BenchmarkEncode_SmallDump(b *testing.B) { encode(b, "testdata/small.json") }

func encode(b *testing.B, filename string) {

	objects := decodeAll(b, jsonReader(b, filename))

	w := bytes.NewBuffer(make([]byte, 0, 1<<23))

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		enc := rubyobj.NewEncoder(w)
		for _, obj := range objects {

			if err := enc.Encode(&obj); err != nil {
				b.Fatal(err)
			}
		}
		w.Reset()
	}
}

func jsonReader(b *testing.B, filename string) *bytes.Buffer {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		b.Fatal(err)
	}
	return bytes.NewBuffer(data)
}

func decodeAll(b *testing.B, r io.Reader) (objects []rubyobj.RubyObject) {

	rObj := rubyobj.RubyObject{}
	var err error

	dec := rubyobj.NewDecoder(r)
	for err != io.EOF {
		err = dec.Decode(&rObj)
		if err != nil && err != io.EOF {
			b.Fatal(err)
		}
		objects = append(objects, rObj)

	}
	return
}
