package worker

import (
	"bytes"
	"edouard127/copingheimer/src/core"
	net2 "edouard127/copingheimer/src/net"
	"edouard127/copingheimer/src/net/packet"
	"encoding/json"
	"fmt"
	provider "github.com/Tnze/go-mc/bot"
	pk "github.com/Tnze/go-mc/net/packet"
	"go.uber.org/zap"
	"net"
	"time"
)

type Worker struct {
	*zap.Logger
	*packet.Conn
	generator *Iterator[net.IP]
	limiter   chan int
	signal    chan core.Signal

	Events  *packet.Events[Worker]
	Options *core.ClientOption
}

func NewWorker(option *core.ClientOption) *Worker {
	logger, _ := zap.NewDevelopment()

	addr, _ := net.ResolveTCPAddr("tcp", option.Node)
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		logger.Panic("cannot connect to server", zap.Error(err))
	}
	w := &Worker{
		Logger:  logger.With(zap.String("node", addr.String())),
		Conn:    packet.NewConnection(conn),
		limiter: make(chan int, option.Instances),
		signal:  make(chan core.Signal),
		Events:  packet.NewEvents[Worker](),
		Options: option,
	}

	w.generator = NewIterator[net.IP]().SetNext(func(i *Iterator[net.IP]) net.IP {
		return func() net.IP {
			return core.GetIP(core.DefaultIP, i.Index())
		}()
	})

	return w
}

func (w *Worker) Handle() {
	fmt.Println("Let's all cope together !")

	w.WritePacket(pk.Marshal(net2.SPacketLogin, pk.Int(w.Options.Instances)))

	go w.handlePacket()

	core.RunSignal(w.scan, w.signal)
}

func (w *Worker) handlePacket() {
	var p pk.Packet
	for {
		if err := w.Conn.ReadPacket(&p); err != nil {
			w.Error("error reading packet", zap.Error(err))
		}

		if err := w.Events.HandlePacket(w, p); err != nil {
			w.Error("error:", zap.Error(err))
		}
	}
}

func (w *Worker) scan() {
	w.limiter <- 1

	go func() {
		defer func() { <-w.limiter }()

		ip := w.generator.Next()
		w.WritePacket(pk.Marshal(net2.SPacketIP, pk.Int(core.IPInt(ip))))
		data, _, err := provider.PingAndListTimeout(ip.String(), time.Duration(w.Options.Timeout)*time.Millisecond)
		if err != nil {
			w.Debug("server did not respond", zap.Error(err))
		} else {
			w.Debug("server response", zap.ByteString("response", data))

			r := bytes.NewReader(data)
			var status core.StatusResponse
			json.NewDecoder(r).Decode(&status)
			status.IP = ip.String()
			w.WritePacket(pk.Marshal(net2.SPacketServer, &status))
		}
	}()
}
