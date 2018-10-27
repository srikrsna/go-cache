package protobuf

import (
	"bytes"
	"errors"
	"io"
	"sync"

	"github.com/golang/protobuf/proto"
)

// ErrNotProtobufMessage is returned when a non protobuf message type is passed for marshalling or unmarshaling
var ErrNotProtobufMessage = errors.New("non protobuf message provided for marshaling/unmarshaling")

// Protobuf implements the MarhsalUnmarshaler interface.
type Protobuf struct{}

// Marshal to marshals a data structure into the given io.Writer
func (Protobuf) Marshal(w io.Writer, v interface{}) error {

	p, ok := v.(proto.Message)
	if !ok {
		return ErrNotProtobufMessage
	}

	buf := wpool.Get().(*proto.Buffer)
	buf.Reset()

	if err := buf.Marshal(p); err != nil {
		return err
	}

	_, err := w.Write(buf.Bytes())

	wpool.Put(buf)

	return err
}

// Unmarshal unmarshals the data structure present in r in its encoded form into v. v should be a pointer type.
func (Protobuf) Unmarshal(r io.Reader, v interface{}) error {
	p, ok := v.(proto.Message)
	if !ok {
		return ErrNotProtobufMessage
	}

	buf := rpool.Get().(*bytes.Buffer)
	buf.Reset()

	if _, err := io.Copy(buf, r); err != nil {
		return err
	}
	err := proto.Unmarshal(buf.Bytes(), p)

	rpool.Put(buf)

	return err
}

var wpool = sync.Pool{
	New: func() interface{} {
		return new(proto.Buffer)
	},
}

var rpool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}
