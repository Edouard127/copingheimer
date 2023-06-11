package server

import (
	"container/list"
	"context"
	"edouard127/copingheimer/src/core"
	pk "github.com/Tnze/go-mc/net/packet"
	"sync/atomic"
	"time"
)

// keepAliveInterval represents the interval when the server sends keep alive
const keepAliveInterval = time.Second * 10

// keepAliveWaitInterval represents how long does the worker have to answer
const keepAliveWaitInterval = time.Second * 3

type ClientManager struct {
	join chan AliveClient
	quit chan AliveClient

	pingList  *list.List
	waitList  *list.List
	listIndex map[AliveClient]*list.Element
	listTimer *time.Timer
	waitTimer *time.Timer

	effort         chan BestEffortClient
	workerWaitList *list.List
	workerIndex    map[int32]BestEffortClient
	maxIP          uint

	id            atomic.Int64
	clientCounter atomic.Int64
}

type AliveClient interface {
	SendDisconnect()
	SendKeepAlive(id int64)
}

type WorkerClient interface {
	GetID() int32
	GetInstances() int32
	SendOffsetIP(offset int32)
	SendSignalID(signal core.Signal)
}

type BestEffortClient interface {
	SendPacket(p pk.Packet)
	AliveClient
	WorkerClient
}

func NewClientManager(ip uint) *ClientManager {
	return &ClientManager{
		join:           make(chan AliveClient),
		quit:           make(chan AliveClient),
		pingList:       list.New(),
		waitList:       list.New(),
		listIndex:      make(map[AliveClient]*list.Element),
		listTimer:      time.NewTimer(keepAliveInterval),
		waitTimer:      time.NewTimer(keepAliveWaitInterval),
		effort:         make(chan BestEffortClient),
		workerWaitList: list.New(),
		workerIndex:    make(map[int32]BestEffortClient, 0),
		maxIP:          ip,
		id:             atomic.Int64{},
		clientCounter:  atomic.Int64{},
	}
}

func (m *ClientManager) ClientJoin(client AliveClient) { m.join <- client }
func (m *ClientManager) ClientLeft(client AliveClient) { m.quit <- client }

func (m *ClientManager) WorkerJoin(client BestEffortClient) {
	m.workerIndex[client.GetID()] = client
	m.clientCounter.Add(1)
	m.maxIP = core.IPInt(core.DefaultIP) + uint(m.GetOffsetMultiplier())
}
func (m *ClientManager) WorkerWait(client BestEffortClient) { m.effort <- client }

// Run implement Component for ClientManager
func (m *ClientManager) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case c := <-m.join:
			m.pushClient(c)
		case c := <-m.quit:
			m.removeClient(c)
		case now := <-m.listTimer.C:
			m.pingClient(now)
		case <-m.waitTimer.C:
			m.disconnectClient()
		case c := <-m.effort:
			m.waitClient(c)
		}
	}
}

func (m *ClientManager) pushClient(c AliveClient) {
	m.listIndex[c] = m.pingList.PushBack(
		keepAliveItem{client: c, t: time.Now()},
	)
}

func (m *ClientManager) removeClient(c AliveClient) {
	elem := m.listIndex[c]
	delete(m.listIndex, c)
	if elem.Prev() == nil {
		// At present, it is difficult to distinguish
		// which linked list the client is in,
		// so both timers will be reset
		defer keepAliveSetTimer(m.pingList, m.listTimer, keepAliveInterval)
		defer keepAliveSetTimer(m.waitList, m.waitTimer, keepAliveWaitInterval)
	}
	m.pingList.Remove(elem)
	m.waitList.Remove(elem)
}

func (m *ClientManager) pingClient(now time.Time) {
	if elem := m.pingList.Front(); elem != nil {
		c := m.pingList.Remove(elem).(keepAliveItem).client

		// Send keep alive id
		c.SendKeepAlive(m.id.Load())
		m.id.Add(1)

		// Clientbound ClientManager packet is sent, move the client to waiting list.
		m.listIndex[c] = m.waitList.PushBack(
			keepAliveItem{client: c, t: now},
		)
	}
	// Wait for next earliest client
	keepAliveSetTimer(m.pingList, m.listTimer, keepAliveInterval)
}

func (m *ClientManager) tickClient(c BestEffortClient) {
	elem, _ := m.listIndex[c]

	if elem.Prev() == nil {
		if !m.waitTimer.Stop() {
			<-m.waitTimer.C
		}
		defer keepAliveSetTimer(m.waitList, m.waitTimer, keepAliveWaitInterval)
	}

	// move the client to ping list
	m.listIndex[c] = m.pingList.PushBack(
		keepAliveItem{client: c, t: time.Now()},
	)
}

func (m *ClientManager) disconnectClient() {
	if elem := m.waitList.Front(); elem != nil {
		c := m.waitList.Remove(elem).(keepAliveItem).client
		m.waitList.Remove(elem)
		m.ClientLeft(c)
		c.SendDisconnect()
	}
	keepAliveSetTimer(m.waitList, m.waitTimer, keepAliveWaitInterval)
}

func keepAliveSetTimer(l *list.List, timer *time.Timer, interval time.Duration) {
	if first := l.Front(); first != nil {
		item := first.Value.(keepAliveItem)
		interval -= time.Since(item.t)
		if interval < 0 {
			interval = 0
		}
	}
	timer.Reset(interval)
	return
}

type keepAliveItem struct {
	client AliveClient
	t      time.Time
}

func (m *ClientManager) waitClient(c BestEffortClient) {
	delete(m.workerIndex, c.GetID())
	m.workerWaitList.PushBack(c)

	if len(m.workerIndex) == 0 {
		// All clients are waiting, resume the workers
		for elem := m.workerWaitList.Front(); elem != nil; elem = elem.Next() {
			c := elem.Value.(BestEffortClient)
			m.workerIndex[c.GetID()] = c
			c.SendOffsetIP(c.GetInstances())
			c.SendSignalID(core.Resume)
		}
		m.maxIP += uint(m.GetOffsetMultiplier())
		m.workerWaitList.Init()
	}
}

func (m *ClientManager) GetOffsetMultiplier() int32 {
	var multiplier int32 = 0
	for _, c := range m.workerIndex {
		multiplier += c.GetInstances()
	}
	return multiplier
}
