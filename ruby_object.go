package rubyobj

import (
	"fmt"
	"strconv"
	"strings"
)

// RubyObject is the deserialized form of an object in an ObjectSpace dump.
type RubyObject struct {
	Type  RubyType
	Value interface{}
	Name  string

	NodeType string

	Address    uint64
	Class      uint64
	References []uint64

	Default    uint64
	Generation uint64

	Bytesize uint64

	Fd       int
	File     string
	Encoding string

	Method string

	Ivars    uint64
	Length   uint64
	Line     uint64
	Memsize  uint64
	Capacity uint64
	Size     uint64

	Struct string
	flags  flagType
}

func (ro RubyObject) Broken() bool {
	return ro.flags&broken != 0
}

func (ro RubyObject) Frozen() bool {
	return ro.flags&frozen != 0
}

func (ro RubyObject) Fstring() bool {
	return ro.flags&fstring != 0
}

func (ro RubyObject) GcMarked() bool {
	return ro.flags&gcMarked != 0
}

func (ro RubyObject) GcOld() bool {
	return ro.flags&gcOld != 0
}

func (ro RubyObject) GcWbProtected() bool {
	return ro.flags&gcWbProtected != 0
}

func (ro RubyObject) Shared() bool {
	return ro.flags&shared != 0
}

func (ro RubyObject) Embedded() bool {
	return ro.flags&embedded != 0
}

func (r *RubyObject) loadSchema(schema *objectSchema) error {

	var err error
	var errs []string

	accumulate := func(err error) {
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	r.Type, err = typeFromName(schema.Type)
	accumulate(err)

	r.Address, err = parseHexUint64(schema.Address)
	accumulate(err)

	r.Class, err = parseHexUint64(schema.Class)
	accumulate(err)

	r.References, err = parseEachUint64(schema.References)
	accumulate(err)

	r.Default, err = parseHexUint64(schema.Default)
	accumulate(err)

	r.NodeType = schema.NodeType

	r.Line = schema.Line
	r.Method = schema.Method
	r.File = schema.File
	r.Fd = schema.Fd
	r.Bytesize = schema.Bytesize
	r.Capacity = schema.Capacity
	r.Length = schema.Length
	r.Size = schema.Size
	r.Encoding = schema.Encoding

	r.Name = schema.Name
	r.Struct = schema.Struct
	r.Ivars = schema.Ivars
	r.Generation = schema.Generation
	r.Memsize = schema.Memsize

	r.flags = flagsFromSchema(schema)

	switch r.Type {
	case Float:
		r.Value, err = strconv.ParseFloat(schema.Value, 64)
		accumulate(err)
	default:
		r.Value = schema.Value
	}

	if len(errs) != 0 {
		return fmt.Errorf("got %d errors decoding Ruby object: %s", len(errs), strings.Join(errs, ", "))
	}
	return nil
}

func (ro *RubyObject) saveSchema() (schema *objectSchema) {
	return &objectSchema{
		Address:    formatUint64(ro.Address),
		Class:      formatUint64(ro.Class),
		NodeType:   ro.NodeType,
		References: formatEachUint64(ro.References),
		Type:       ro.Type.Name(),
		Value:      fmt.Sprintf("%v", ro.Value),
		Line:       ro.Line,
		Method:     ro.Method,
		File:       ro.File,
		Fd:         ro.Fd,
		Bytesize:   ro.Bytesize,
		Capacity:   ro.Capacity,
		Length:     ro.Length,
		Size:       ro.Size,
		Encoding:   ro.Encoding,
		Default:    formatUint64(ro.Default),
		Name:       ro.Name,
		Struct:     ro.Struct,
		Ivars:      ro.Ivars,
		Generation: ro.Generation,
		Memsize:    ro.Memsize,
		Frozen:     ro.Frozen(),
		Embedded:   ro.Embedded(),
		Broken:     ro.Broken(),
		Fstring:    ro.Fstring(),
		Shared:     ro.Shared(),
		Flags: flagSchema{
			WbProtected: ro.GcWbProtected(),
			Old:         ro.GcOld(),
			Marked:      ro.GcMarked(),
		},
	}
}

func parseHexUint64(hexStr string) (uint64, error) {
	if len(hexStr) != 14 {
		return 0, nil
	}
	return strconv.ParseUint(hexStr[2:], 16, 64)
}

func formatUint64(ui uint64) string {
	hex := make([]byte, 0, 14)
	hex = append(hex, byte('0'))
	hex = append(hex, byte('x'))

	val := fmt.Sprintf("%x", ui)
	padLen := cap(hex) - len(val)
	for i := len(hex); i < padLen; i++ {
		hex = append(hex, byte('0'))
	}

	for _, r := range val {
		hex = append(hex, byte(r))
	}
	return string(hex)
}

func parseEachUint64(hexArr []string) ([]uint64, error) {
	var tmp uint64
	var err error
	out := make([]uint64, 0, len(hexArr))
	for i, hexStr := range hexArr {
		tmp, err = parseHexUint64(hexStr)
		if err != nil {
			return nil, fmt.Errorf("%dth value (%s), %v", i, hexStr, err)
		}
		out = append(out, tmp)
	}
	return out, nil
}

func formatEachUint64(uiArr []uint64) []string {
	hexArr := make([]string, 0, len(uiArr))
	for _, ui := range uiArr {
		hexArr = append(hexArr, formatUint64(ui))
	}
	return hexArr
}
