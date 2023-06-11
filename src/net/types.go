package net

import "io"

type (
	// UnsignedInt is a unsigned int
	UnsignedInt uint
)

func (ui UnsignedInt) WriteTo(w io.Writer) (int64, error) {
	n := uint32(ui)
	nn, err := w.Write([]byte{byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n)})
	return int64(nn), err
}

func (ui *UnsignedInt) ReadFrom(r io.Reader) (n int64, err error) {
	var b [4]byte
	nn, err := r.Read(b[:])
	if err != nil {
		return int64(nn), err
	}
	n = int64(nn)
	*ui = UnsignedInt(uint(b[0])<<24 | uint(b[1])<<16 | uint(b[2])<<8 | uint(b[3]))
	return
}
