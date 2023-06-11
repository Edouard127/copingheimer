package packet

import (
	"fmt"
	pk "github.com/Tnze/go-mc/net/packet"
)

type Events[T any] struct {
	handlers map[int32]*handlerHeap[T] // for specific packet id only
}

func NewEvents[T any]() *Events[T] {
	return &Events[T]{handlers: make(map[int32]*handlerHeap[T], 0)}
}

func (e *Events[T]) HandlePacket(t *T, p pk.Packet) (err error) {
	if listeners := e.handlers[p.ID]; listeners != nil {
		for _, handler := range *listeners {
			if err = handler.F(t, p); err != nil {
				return fmt.Errorf("error handling packet %d: %w", p.ID, err)
			}
		}
	}
	return
}

func (e *Events[T]) AddListener(listeners ...Handler[T]) {
	for _, l := range listeners {
		var s *handlerHeap[T]
		var ok bool
		if s, ok = e.handlers[l.ID]; !ok {
			s = &handlerHeap[T]{l}
			e.handlers[l.ID] = s
		} else {
			s.Push(l)
		}
	}
}

type Handler[T any] struct {
	ID int32
	F  func(*T, pk.Packet) error
}

type handlerHeap[T any] []Handler[T]

func (h handlerHeap[T]) Len() int            { return len(h) }
func (h handlerHeap[T]) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *handlerHeap[T]) Push(x interface{}) { *h = append(*h, x.(Handler[T])) }
func (h *handlerHeap[T]) Pop() interface{} {
	old := *h
	n := len(old)
	*h = old[0 : n-1]
	return old[n-1]
}
