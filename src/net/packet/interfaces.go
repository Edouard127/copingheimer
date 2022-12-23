package packet

import (
	"edouard127/copingheimer/src/server/http"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type Listener struct{ net.TCPListener }

func ListenProvider(addr net.TCPAddr) (*Listener, error) {
	l, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		return nil, err
	}
	return &Listener{*l}, nil
}

func (l Listener) Accept() (*Conn, error) {
	conn, err := l.AcceptTCP()
	return NewConnection(conn), err
}

type ConnAction struct {
	Wait bool
}

type Conn struct {
	Socket    *net.TCPConn
	id        int32
	lastState int32 // Last IP received
	states    *States
}

func NewConnection(conn *net.TCPConn) *Conn {
	return &Conn{
		Socket: conn,
		states: NewStates(),
	}
}

func (c *Conn) WritePacket(p Packet) (err error) {
	if c.states.Has(Offline) {
		return fmt.Errorf("cannot write packet to offline connection")
	}
	b := make([]byte, len(p.Data)+8)
	binary.BigEndian.PutUint32(b[0:4], uint32(p.ID))
	binary.BigEndian.PutUint32(b[4:8], uint32(len(p.Data)))
	copy(b[8:], p.Data)
	_, err = c.Socket.Write(b)
	if err != nil {
		return fmt.Errorf("error writing packet %d: %w", p.ID, err)
	}
	return
}

func (c *Conn) ReadPacket(p *Packet) error {
	var (
		Length   int32
		PacketID int32
	)

	if err := c.readBytes(&PacketID, &Length); err != nil {
		return err
	}
	p.ID = PacketID
	p.Data = make([]byte, Length)
	if _, err := io.ReadFull(c.Socket, p.Data); err != nil {
		return err
	}
	return nil
}

func (c *Conn) readBytes(is ...*int32) error {
	var (
		err error
	)
	for i := 0; i < len(is); i++ {
		b := make([]byte, 4)
		if _, err = c.Socket.Read(b); err != nil {
			return err
		}
		*is[i] = int32(binary.BigEndian.Uint32(b))
	}
	return nil
}

type Server struct {
	EventHandlers
	Events
	Dashboard *http.Dashboard
	clients   []*Conn
}

func NewServer(mongo string) *Server {
	return &Server{
		EventHandlers: EventHandlers{},
		Events:        Events{handlers: make(map[int32]*handlerHeap, 0)},
		Dashboard:     http.NewDashboard(mongo),
		clients:       make([]*Conn, 0),
	}
}

func (c *Server) HandlePacket(conn *Conn, p Packet) (err error) {
	if listeners := c.Events.handlers[p.ID]; listeners != nil {
		for _, handler := range *listeners {
			if err = handler.F(c, conn, p); err != nil {
				return fmt.Errorf("error handling packet %d: %w", p.ID, err)
			}
		}
	}
	return
}

func (c *Server) Add(conn *Conn) {
	c.clients = append(c.clients, conn)
}

func (c *Server) Remove(conn *Conn) {
	conn.Socket.Close()
	for i, client := range c.clients {
		if client == conn {
			if err := client.Socket.Close(); err != nil {
				panic(err)
			}
			c.clients = append(c.clients[:i], c.clients[i+1:]...)
		}
	}
}

func (s *Server) FindClient(ip net.IP) *Conn {
	for _, client := range s.clients {
		if client.Socket.RemoteAddr().(*net.TCPAddr).IP.Equal(ip) {
			return client
		}
	}
	return nil
}
