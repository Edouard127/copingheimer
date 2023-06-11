package packet

import (
	"io"
	"net"

	pk "github.com/Tnze/go-mc/net/packet"
)

const MaxPacketSize = 1024 // 1KB

type Listener struct{ net.TCPListener }

func ListenProvider(addr net.TCPAddr) *Listener {
	l, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		panic(err)
	}
	return &Listener{*l}
}

func (l Listener) Accept() (*Conn, error) {
	conn, err := l.AcceptTCP()
	return NewConnection(conn), err
}

type Conn struct {
	Socket net.Conn
	io.Reader
	io.Writer
	ID int32
}

func NewConnection(conn net.Conn) *Conn {
	return &Conn{
		Socket: conn,
		Reader: conn,
		Writer: conn,
	}
}

func (c *Conn) IP() net.IP {
	ip := c.Socket.RemoteAddr().String()
	p, _, _ := net.SplitHostPort(ip)
	return net.ParseIP(p)
}

func (c *Conn) ReadPacket(p *pk.Packet) error {
	return p.UnPack(c.Reader, -1)
}

func (c *Conn) WritePacket(p pk.Packet) error {
	return p.Pack(c.Writer, -1) // uncompressed
}
