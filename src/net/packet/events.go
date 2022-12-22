package packet

/**
Taken from github.com\tnze\go-mc@v1.19.2\bot\event.go
*/

type Events struct {
	handlers map[int32]*handlerHeap // for specific packet id only
}

func (e *Events) AddListener(listeners ...PacketHandler) {
	for _, l := range listeners {
		var s *handlerHeap
		var ok bool
		if s, ok = e.handlers[l.ID]; !ok {
			s = &handlerHeap{l}
			e.handlers[l.ID] = s
		} else {
			s.Push(l)
		}
	}
}

type PacketHandler struct {
	ID int32
	F  func(*Server, *Conn, Packet) error
}

type handlerHeap []PacketHandler

func (h handlerHeap) Len() int            { return len(h) }
func (h handlerHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *handlerHeap) Push(x interface{}) { *h = append(*h, x.(PacketHandler)) }
func (h *handlerHeap) Pop() interface{} {
	old := *h
	n := len(old)
	*h = old[0 : n-1]
	return old[n-1]
}
