package rubyobj_test

import (
	"bytes"
	"github.com/aybabtme/rubyobj"
	"io"
	"io/ioutil"
	"runtime"
	"sync"
	"testing"
)

func BenchmarkDecode_TinyDump(b *testing.B)   { decode(b, "testdata/tiny.json") }
func BenchmarkDecode_SmallDump(b *testing.B)  { decode(b, "testdata/small.json") }
func BenchmarkDecode_MediumDump(b *testing.B) { decode(b, "testdata/medium.json") }
func BenchmarkDecode_BigDump(b *testing.B)    { decode(b, "testdata/big.json") }
func BenchmarkDecode_HugeDump(b *testing.B)   { decode(b, "testdata/huge.json") }

func BenchmarkParallelDecode_TinyDump(b *testing.B)   { parallelDecode(b, "testdata/tiny.json") }
func BenchmarkParallelDecode_SmallDump(b *testing.B)  { parallelDecode(b, "testdata/small.json") }
func BenchmarkParallelDecode_MediumDump(b *testing.B) { parallelDecode(b, "testdata/medium.json") }
func BenchmarkParallelDecode_BigDump(b *testing.B)    { parallelDecode(b, "testdata/big.json") }
func BenchmarkParallelDecode_HugeDump(b *testing.B)   { parallelDecode(b, "testdata/huge.json") }

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

func parallelDecode(b *testing.B, filename string) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(runtime.NumCPU()))

	r := jsonReader(b, filename)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {

		objC, errC := rubyobj.ParallelDecode(r, uint(1))

		wg := sync.WaitGroup{}
		wg.Add(2)
		go readObj(&wg, objC, b)
		go readErr(&wg, errC, b)
		wg.Wait()

		r.Reset()
	}
}

func readObj(wg *sync.WaitGroup, objC <-chan rubyobj.RubyObject, b *testing.B) {
	defer wg.Done()
	for obj := range objC {
		_ = obj
	}
}

func readErr(wg *sync.WaitGroup, errC <-chan error, b *testing.B) {
	defer wg.Done()
	for err := range errC {
		b.Logf("error: %v", err)
	}
}

func BenchmarkEncode_TinyDump(b *testing.B)   { encode(b, "testdata/tiny.json") }
func BenchmarkEncode_SmallDump(b *testing.B)  { encode(b, "testdata/small.json") }
func BenchmarkEncode_MediumDump(b *testing.B) { encode(b, "testdata/medium.json") }
func BenchmarkEncode_BigDump(b *testing.B)    { encode(b, "testdata/big.json") }
func BenchmarkEncode_HugeDump(b *testing.B)   { encode(b, "testdata/huge.json") }

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
