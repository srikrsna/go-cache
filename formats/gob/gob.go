package gob

import (
	"encoding/gob"
	"io"
)

// Gob implements the MarhsalUnmarshaler interface. It used go's excellent binary encoding format 'encoding/gob'.
// Gob requires types to be registered before they can marshalled/unmarshalled
type Gob struct{}

// Marshal to marshals a data structure into the given io.Writer
func (Gob) Marshal(w io.Writer, v interface{}) error {
	return gob.NewEncoder(w).Encode(v)
}

// Unmarshal unmarshals the data structure present in r in its encoded form into v. v should be a pointer type.
func (Gob) Unmarshal(r io.Reader, v interface{}) error {
	return gob.NewDecoder(r).Decode(v)
}
