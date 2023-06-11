package server

import (
	"edouard127/copingheimer/src/core"
	"edouard127/copingheimer/src/net/packet"
	pk "github.com/Tnze/go-mc/net/packet"
	"go.uber.org/zap"
	"net"
)

type Server struct {
	*zap.Logger
	*packet.Listener
	database      *core.Database
	Events        *packet.Events[ServerClient]
	ClientManager *ClientManager
}

func NewServer(options *core.ServerOption) *Server {
	logger, _ := zap.NewDevelopment()
	return &Server{
		Logger: logger.With(zap.String("module", "server")),
		Listener: packet.ListenProvider(net.TCPAddr{
			IP:   net.ParseIP(options.Host),
			Port: 29969,
		}),
		database:      core.InitDatabase(options.MongoDB),
		Events:        packet.NewEvents[ServerClient](),
		ClientManager: NewClientManager(options.StartIP),
	}
}

func (s *Server) Handle() {
	s.Info("server started")
	for {
		conn, err := s.Accept()
		if err != nil {
			panic(err)
		}

		client := NewServerClient(s.With(zap.String("client", conn.Socket.RemoteAddr().String())), conn)

		s.ClientManager.ClientJoin(client)
		go s.handlePacket(client)
	}
}

func (s *Server) handlePacket(client *ServerClient) {
	var p pk.Packet
	for {
		if err := client.conn.ReadPacket(&p); err != nil {
			s.Error("error reading packet", zap.Error(err))
		}

		if err := s.Events.HandlePacket(client, p); err != nil {
			s.Error("error handling packet", zap.Error(err))
		}
	}
}
