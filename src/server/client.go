package server

import (
	"edouard127/copingheimer/src/core"
	"edouard127/copingheimer/src/net"
	"edouard127/copingheimer/src/net/packet"
	pk "github.com/Tnze/go-mc/net/packet"
	"go.uber.org/zap"
)

// ServerClient represents a client connected to the server.
type ServerClient struct {
	log       *zap.Logger
	conn      *packet.Conn
	ip        int32
	instances int32
	id        int32
}

func NewServerClient(log *zap.Logger, conn *packet.Conn) *ServerClient {
	return &ServerClient{log: log, conn: conn}
}

func (c *ServerClient) SendPacket(p pk.Packet) {
	c.conn.WritePacket(p)
}

func (c *ServerClient) SendDisconnect() {
	c.conn.Socket.Close()
}

func (c *ServerClient) SendKeepAlive(id int64) {
	c.SendPacket(pk.Marshal(net.CPacketKeepAlive, pk.Long(id)))
}

func (c *ServerClient) SendOffsetIP(offset int32) {
	c.SendPacket(pk.Marshal(net.CPacketOffsetIP, pk.Int(offset)))
}

func (c *ServerClient) SendSignalID(signal core.Signal) {
	c.SendPacket(pk.Marshal(net.CPacketSignal, pk.Int(signal)))
}

func (c *ServerClient) GetID() int32 {
	return c.id
}

func (c *ServerClient) GetInstances() int32 {
	return c.instances
}
