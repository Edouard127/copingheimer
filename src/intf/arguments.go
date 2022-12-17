package intf

import (
	"bytes"
	"encoding/binary"
)

type Arguments struct {
	Config    string
	Help      bool
	IP        string
	Instances int
	Timeout   int
}

func (a *Arguments) ReadFrom(b []byte) error {
	reader := bytes.NewReader(b)
	if err := binary.Read(reader, binary.LittleEndian, a); err != nil {
		return err
	}
	return nil
}
