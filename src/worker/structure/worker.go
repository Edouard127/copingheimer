package structure

import (
	"bytes"
	"edouard127/copingheimer/src/intf"
	"edouard127/copingheimer/src/net/packet"
	"edouard127/copingheimer/src/worker/utils"
	"fmt"
	provider "github.com/Tnze/go-mc/bot"
	"net"
	"net/http"
	"time"
)

type Worker struct {
	Socket *packet.Conn
	ip     utils.IteratorFunc

	Arguments *intf.Arguments
	States    *packet.States
	Tasks     int
	Auth      int32
}

func NewWorker(arguments *intf.Arguments) *Worker {
	blackList, err := intf.ReadBlacklist(arguments)
	if err != nil {
		fmt.Println("Error while reading the blacklist:", err)
	}
	return &Worker{
		Socket: nil,
		ip: utils.IPSubnetIterator(net.IPNet{
			IP:   net.ParseIP(arguments.IP),
			Mask: net.IPMask{0, 0, 0, 0},
		}, blackList),
		Arguments: arguments,
		States:    packet.NewStates(),
	}
}

func (w *Worker) Connect(addr net.TCPAddr) {
	conn, err := net.DialTCP("tcp", nil, &addr)
	if err != nil {
		fmt.Println("failed to connect to the node, is the node running?")
		w.States.Set(packet.Offline, true)
		return
	}
	w.Socket = packet.NewConnection(conn)
	fmt.Println("Connected to the server")

	if err := w.Socket.WritePacket(
		packet.Marshal(
			packet.SPacketLogin,
			int32(w.Arguments.Instances),
		),
	); err != nil {
		fmt.Println("failed to send login packet:", err)
		w.States.Set(packet.Offline, true)
		return
	}
}

func (w *Worker) Handle() {
	go w.keepNetwork()
	go w.handleConnection()
	fmt.Println("Let's all cope together !")

	r := /*w.Arguments.Mode == "random"*/ false // TODO: implement random mode
	for {
		if w.States.Has(packet.Wait) {
			fmt.Println("waiting for the server to be ready")
			continue
		}

		if w.Tasks < w.Arguments.Instances {
			ip := w.ip(1, r)
			if ip.Int()%int32(w.Arguments.Instances) == 0 && !w.States.Has(packet.Offline) {
				if err := w.Socket.WritePacket(
					packet.Marshal(
						packet.SPacketIP,
						ip.Int(),
					),
				); err != nil {
					panic(err)
				}
			}
			w.Tasks++
			go w.scan()
		}
	}
}

func (w *Worker) scan() {
	s := w.ip(1, false)
	if data, _, err := provider.PingAndListTimeout(s.String(), time.Duration(w.Arguments.Timeout)*time.Millisecond); err == nil {
		status := &intf.StatusResponse{}
		if err := status.Put(data); err != nil {
			fmt.Println("error while parsing the status:", err)
			return
		}
		if !w.States.Has(packet.Offline) {
			// Post the status to the server
			parsedIP, _, err := net.SplitHostPort(w.Arguments.Node)
			if err != nil {
				fmt.Println("error while parsing the node ip:", err)
				return
			}
			go func() {
				b, err := status.Json(s.String())
				if err != nil {
					fmt.Println("error while marshalling the status:", err)
					return
				}
				req, err := http.NewRequest("POST", "http://"+parsedIP+":80/api/status", bytes.NewBuffer(b))
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					fmt.Println("error while creating the request:", err)
					return
				}
				client := &http.Client{}
				if _, err := client.Do(req); err != nil {
					fmt.Println("error while sending the request:", err)
					return
				}
				req.Body.Close()
			}()
		}
	} else {
		fmt.Println("error while scanning the ip:", err)
	}
	w.Tasks--
}

func (w *Worker) keepNetwork() {
	var last time.Time
	for {
		time.Sleep(5 * time.Second)
		if _, err := net.DialTimeout("tcp", "google.com:80", 5*time.Second); err != nil {
			fmt.Println("connection lost, adjusting instances count")
			w.Arguments.Instances /= 16
			w.ip(int32(-(w.Arguments.Instances * 10)), false)
			last = time.Now()
			if !w.States.Has(packet.Offline) {
				if err := w.Socket.WritePacket(
					packet.Marshal(
						packet.SPacketLogin,
						int32(w.Arguments.Instances),
					),
				); err != nil {
					panic(err)
				}
			}
		} else {
			if last.Add(5*time.Second).Before(time.Now()) && last != (time.Time{}) {
				fmt.Println("connection recovered")
				last = time.Now()
			}
		}
	}
}

func (w *Worker) handleConnection() {
	var p packet.Packet
	for {
		if err := w.Socket.ReadPacket(&p); err == nil {
			if p.ID == packet.CPacketLogin {
				var (
					instances  int32
					serverAuth int32
				)
				if err := p.Scan(&instances, &serverAuth); err != nil {
					fmt.Println("failed to read login packet:", err)
					return
				}
				w.Auth = serverAuth
				w.Arguments.Instances = int(instances)
			}
		}
	}
}
