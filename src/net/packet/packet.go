package packet

import (
	"bytes"
	"encoding/binary"
	"io"
)

// Packet define a net data package
type Packet struct {
	ID   int32
	Data []byte
}

// Scan decode the packet and fill data into fields
func (p Packet) Scan(fields ...*int32) error {
	r := bytes.NewReader(p.Data)
	for _, v := range fields {
		if err := binary.Read(r, binary.BigEndian, v); err != nil {
			if err == io.EOF {
				break
			}
		}
	}
	return nil
}

// Build a packet
func (p Packet) Build() []byte {
	b := make([]byte, len(p.Data)+8)
	binary.BigEndian.PutUint32(b[0:4], uint32(p.ID))
	binary.BigEndian.PutUint32(b[4:8], uint32(len(p.Data)))
	copy(b[8:], p.Data)
	return b
}

type Builder struct {
	buf bytes.Buffer
}

func (p *Builder) WriteField(fields ...int32) {
	for _, v := range fields {
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, uint32(v))
		p.buf.Write(b)
	}
}

func (p *Builder) Packet(id int32) Packet {
	return Packet{ID: id, Data: p.buf.Bytes()}
}

func Marshal(id int32, fields ...int32) (pk Packet) {
	var pb Builder
	for _, v := range fields {
		pb.WriteField(v)
	}
	return pb.Packet(id)
}
