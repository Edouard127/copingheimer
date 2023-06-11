package worker

import (
	"edouard127/copingheimer/src/core"
	"edouard127/copingheimer/src/net"
	pk "github.com/Tnze/go-mc/net/packet"
)

func KeepAlive(w *Worker, p pk.Packet) error {
	p.ID = net.SPacketKeepAlive
	return w.WritePacket(p)
}

func OffsetIP(w *Worker, p pk.Packet) error {
	var ip net.UnsignedInt

	if err := p.Scan(&ip); err != nil {
		return err
	}

	w.generator.index += uint(ip)
	return nil
}

func Signal(w *Worker, p pk.Packet) error {
	var signal pk.Int

	if err := p.Scan(&signal); err != nil {
		return err
	}

	w.signal <- core.Signal(signal)

	return nil
}
