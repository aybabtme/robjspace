package rubyobj

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aybabtme/fatherhood"
	"io"
	// "log"
	"sync"
)

// Trivial codec

// Decoder decodes RubyObjects from an io.Reader.  It wraps a json.Decoder and is
// pretty slow, but simple to use for those acustomed to the json.Decoder from
// the stdlib.
type Decoder struct {
	dec *json.Decoder
}

// NewDecoder returns a trivial decoder wrapping a json.Decoder of the stdlib.
// It is pretty slow but simple to use.
//
// For performance, prefer ParallelDecode.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{json.NewDecoder(r)}
}

var schema objectSchema

// Decode decodes a ruby object from the underlying io.Reader.
func (d *Decoder) Decode(rObj *RubyObject) (err error) {
	schema.clear()
	err = d.dec.Decode(&schema)
	if err != nil {
		return err
	}
	return rObj.loadSchema(&schema)
}

// Encoder encodes RubyObjects to an io.Writer.  It wraps a json.Encoder and is
// pretty slow, but simple to use for those acustomed to the json.Encoder from
// the stdlib.
type Encoder struct {
	enc *json.Encoder
}

// NewEncoder returns a trivial encoder wrapping a json.Encoder of the stdlib.
// It is pretty slow but simple to use.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{json.NewEncoder(w)}
}

// Encode encodes the object onto the underlying io.Writer.
func (e *Encoder) Encode(rObj *RubyObject) (err error) {
	return e.enc.Encode(rObj.saveSchema())
}

// Parallel

// ParallelDecode will use many goroutines to decode io.Reader.  io.Reader MUST
// present JSON objects seperated by \n characters.
//
// Decoding will use para + 1 goroutines:
//      1 x goroutines to read all the lines in the io.Reader
//   para x goroutines to decode the lines
//
// The decoding will continue until it reaches EOF in the io.Reader, and return
// all the errors it encountered on the error channel, including Read errors,
// unmarshalling errors and loading errors.
func ParallelDecode(r io.Reader, para uint) (<-chan RubyObject, <-chan error) {

	bufLen := para << 2
	decodedC := make(chan RubyObject, bufLen)
	errc := make(chan error, bufLen)

	go func() {
		defer close(decodedC)
		defer close(errc)

		lineC := scanLines(r, errc, bufLen)

		wg := sync.WaitGroup{}
		// log.Printf("[para] starting workers")
		for i := uint(0); i < para; i++ {
			// log.Printf("[para] -> worker %d", i)
			wg.Add(1)
			go decodeLines(&wg, lineC, decodedC, errc)
		}
		// log.Printf("[para] waiting for workers")
		wg.Wait()
	}()

	return decodedC, errc

}

func scanLines(r io.Reader, errc chan<- error, para uint) <-chan []byte {
	linesC := make(chan []byte, para)

	go func() {
		defer close(linesC)

		br := bufio.NewReader(r)

		var line []byte
		var err error

		// log.Printf("[scan] scanning...")
		// defer log.Printf("[scan] done")
		for {
			line, err = br.ReadBytes('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				// log.Printf("[scan] -> error")
				errc <- err
			}
			linesC <- line
			// log.Printf("[scan] -> scanned")

		}

	}()

	return linesC
}

func decodeLines(wg *sync.WaitGroup, lineC <-chan []byte, decoded chan<- RubyObject, errc chan<- error) {
	defer wg.Done()

	rObj := RubyObject{}
	schema := objectSchema{}
	var err error

	r := bytes.NewBuffer(nil)
	dec := fatherhood.NewDecoder(r)

	for line := range lineC {

		_, _ = r.Write(line)
		err = dec.EachMember(&schema, decodeObjSchema)
		if err != nil {
			errc <- err
			continue
		}

		err = rObj.loadSchema(&schema)
		if err != nil {
			errc <- err
			continue
		}
		decoded <- rObj

	}

}

func decodeObjSchema(dec *fatherhood.Decoder, s interface{}, member string) error {
	schema := s.(*objectSchema)
	switch member {
	case "address":
		return dec.ReadString(&schema.Address)
	case "class":
		return dec.ReadString(&schema.Class)
	case "node_type":
		return dec.ReadString(&schema.NodeType)
	case "references":
		schema.References = make([]string, 0)
		return dec.EachValue(&schema.References, decodeReference)
	case "type":
		return dec.ReadString(&schema.Type)
	case "value":
		return dec.ReadString(&schema.Value)
	case "line":
		return dec.ReadUint64(&schema.Line)
	case "method":
		return dec.ReadString(&schema.Method)
	case "file":
		return dec.ReadString(&schema.File)
	case "fd":
		return dec.ReadInt(&schema.Fd)
	case "bytesize":
		return dec.ReadUint64(&schema.Bytesize)
	case "capacity":
		return dec.ReadUint64(&schema.Capacity)
	case "length":
		return dec.ReadUint64(&schema.Length)
	case "size":
		return dec.ReadUint64(&schema.Size)
	case "encoding":
		return dec.ReadString(&schema.Encoding)
	case "default":
		return dec.ReadString(&schema.Default)
	case "name":
		return dec.ReadString(&schema.Name)
	case "struct":
		return dec.ReadString(&schema.Struct)
	case "ivars":
		return dec.ReadUint64(&schema.Ivars)
	case "generation":
		return dec.ReadUint64(&schema.Generation)
	case "memsize":
		return dec.ReadUint64(&schema.Memsize)
	case "frozen":
		return dec.ReadBool(&schema.Frozen)
	case "embedded":
		return dec.ReadBool(&schema.Embedded)
	case "broken":
		return dec.ReadBool(&schema.Broken)
	case "fstring":
		return dec.ReadBool(&schema.Fstring)
	case "shared":
		return dec.ReadBool(&schema.Shared)
	case "flags":
		schema.Flags = flagSchema{}
		return dec.EachMember(&schema.Flags, decodeFlagSchema)
	}
	// unsupported member
	return dec.Discard()
}

func decodeFlagSchema(dec *fatherhood.Decoder, f interface{}, member string) error {
	flag := f.(*flagSchema)
	switch member {
	case "wb_protected":
		return dec.ReadBool(&flag.WbProtected)
	case "old":
		return dec.ReadBool(&flag.Old)
	case "marked":
		return dec.ReadBool(&flag.Marked)
	}
	// unsupported member
	return dec.Discard()
}

func decodeReference(dec *fatherhood.Decoder, a interface{}, t fatherhood.JSONType) error {
	arr := a.(*[]string)
	switch t {
	case fatherhood.String:
		var v struct {
			val string
		}
		err := dec.ReadString(&v.val)
		if err != nil {
			return err
		}
		*arr = append(*arr, v.val)
		return nil
	}
	return fmt.Errorf("unexpected type in 'reference' array, %#v", t)
}
