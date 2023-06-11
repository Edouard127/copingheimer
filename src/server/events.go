package server

import (
	"edouard127/copingheimer/src/core"
	net2 "edouard127/copingheimer/src/net"
	pk "github.com/Tnze/go-mc/net/packet"
)

func (s *Server) HandleLogin(c *ServerClient, p pk.Packet) error {
	var instances pk.Int

	if err := p.Scan(&instances); err != nil {
		return err
	}

	c.instances = int32(instances)
	c.id = int32(s.ClientManager.clientCounter.Load())

	s.ClientManager.WorkerJoin(c)

	c.SendOffsetIP(s.ClientManager.GetOffsetMultiplier() - c.GetInstances())

	return nil
}

func (s *Server) HandleServer(c *ServerClient, p pk.Packet) error {
	var status core.StatusResponse

	if err := p.Scan(&status); err != nil {
		return err
	}

	return s.database.Write(status.IP, status)
}

func (s *Server) HandleIP(c *ServerClient, p pk.Packet) error {
	var ip net2.UnsignedInt

	if err := p.Scan(&ip); err != nil {
		return err
	}

	if uint(ip) >= s.ClientManager.maxIP {
		c.SendSignalID(core.Pause)
		s.ClientManager.waitClient(c)
	}

	return nil
}
