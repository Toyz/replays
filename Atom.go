package replays

import (
	"bytes"
	"encoding/binary"
	"log"
	"unsafe"
)

type Atom struct {
	Name   string
	Size   int32
	Buffer []byte
}

func NewAtom(data []byte) Atom {
	in := bytes.NewBuffer(data)

	a := Atom{}
	if err := binary.Read(in, binary.BigEndian, &a.Size); err != nil {
		log.Printf("failed to  decode size: %f", err)
		return a
	}

	if a.Size < 8 {
		return a
	}

	nameSlice := make([]byte, 4)
	if count, err := in.Read(nameSlice); err != nil {
		log.Printf("failed to read name: %f", err)
		return a
	} else {
		if count == len(nameSlice) {
			a.Name = string(nameSlice)
		}
	}

	s := len(a.Name) + int(unsafe.Sizeof(a.Size))
	if a.Size > 8 {
		a.Buffer = make([]byte, a.Size - int32(s))
		if _, err := in.Read(a.Buffer); err != nil {
			log.Printf("failed to read buffer: %f", err)
			return a
		}
	}

	return a
}
