package worker

import (
	"sync"
)

type Iterator[T any] struct {
	index   uint
	restart bool // restarts the iterator if the end is reached
	entries []T
	next    func(*Iterator[T]) T // only used if entries is empty
	group   *sync.Mutex
}

func NewIterator[T any]() *Iterator[T] {
	return &Iterator[T]{
		index:   0,
		restart: false,
		entries: make([]T, 0),
		group:   &sync.Mutex{},
	}
}

func (p *Iterator[T]) AddEntry(proxy ...T) *Iterator[T] {
	p.group.Lock()
	defer p.group.Unlock()
	p.entries = append(p.entries, proxy...)
	return p
}

func (p *Iterator[T]) RemoveEntry(index int) *Iterator[T] {
	p.group.Lock()
	defer p.group.Unlock()
	p.entries = append(p.entries[:index], p.entries[index+1:]...)
	return p
}

func (p *Iterator[T]) SetNext(next func(*Iterator[T]) T) *Iterator[T] {
	p.next = next
	return p
}

func (p *Iterator[T]) SetRestart(restart bool) *Iterator[T] {
	p.restart = restart
	return p
}

func (p *Iterator[T]) Next() T {
	p.group.Lock()
	defer p.group.Unlock()

	var next T

	if int(p.index) >= len(p.entries)-1 {
		if p.restart {
			p.index = 0
			next = p.entries[p.index]
			p.index++
		}
	}

	next = p.next(p)
	p.index++

	return next
}

func (p *Iterator[T]) HasNext() bool {
	p.group.Lock()
	defer p.group.Unlock()

	return int(p.index) < len(p.entries)
}

func (p *Iterator[T]) Len() int32 {
	p.group.Lock()
	defer p.group.Unlock()

	return int32(len(p.entries))
}

func (p *Iterator[T]) Index() uint {
	return p.index
}
