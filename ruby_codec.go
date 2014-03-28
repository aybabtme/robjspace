package robjspace

import (
	"encoding/json"
	"io"
)

type Decoder struct {
	dec *json.Decoder
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{json.NewDecoder(r)}
}

func (d *Decoder) Decode(rObj *RubyObject) (err error) {
	var schema *objectSchema
	err = d.dec.Decode(schema)
	if err != nil {
		return err
	}
	return rObj.loadSchema(schema)
}

type Encoder struct {
	enc *json.Encoder
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{json.NewEncoder(w)}
}

func (e *Encoder) Encode(rObj *RubyObject) (err error) {
	return e.enc.Encode(rObj.saveSchema())
}
