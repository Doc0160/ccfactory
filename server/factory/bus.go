package factory

import (
	"ccfactory/server/access"
	"ccfactory/server/debug"
	"ccfactory/server/peripheral"
	"ccfactory/server/server"
	"sync"
)

type Bus struct {
	server    *server.Server
	access    access.BasicAccess
	inventory peripheral.Inventory
	factory   *Factory

	mutex       sync.RWMutex
	allocations map[int]bool

	size    int
	waiters []chan int
}

func NewBus(server *server.Server, access access.BasicAccess, factory *Factory) *Bus {
	return &Bus{
		server:      server,
		access:      access,
		factory:     factory,
		allocations: map[int]bool{},
		inventory: peripheral.Inventory{
			Server: server,
			Access: access,
		},
	}
}

func (b *Bus) Update() error {
	defer debug.Timer("Update")()

	slots, _ := b.inventory.List()
	for i, _ := range slots {
		if slots[i] == nil {
			continue
		}
		if b.allocations[i] {
			continue
		}
		log.Debug(slots[i])
		for _, storage := range b.factory.ItemStorages {
			_ = storage
		}
	}

	return nil
}

func (b *Bus) Transfer(task peripheral.BusTask) error {
	defer debug.Timer("Transfer")()

	b.mutex.Lock()
	if b.size == 0 {
		size, err := b.inventory.Size()
		if err != nil {
			return err
		}
		b.size = size
	}

	freeSlot := -1
	for i := b.size - 1; i >= 0; i-- {
		if !b.allocations[i] {
			freeSlot = i
			b.allocations[i] = true
			break
		}
	}

	if freeSlot == -1 {
		log.Error("not found")
		ch := make(chan int, 1)
		b.waiters = append(b.waiters, ch)
		b.mutex.Unlock()

		freeSlot = <-ch
	} else {
		b.mutex.Unlock()
	}

	_, err := b.server.Call(task.FromClient, &server.Request{
		Type: "peripheral",
		Args: []any{
			task.FromBusAddr,
			"pullItems",
			task.FromInvAddr,
			task.FromInvSlot + 1,
			task.Count,
			freeSlot + 1,
		},
	})
	if err != nil {
		log.Error(err)
	}

	_, err = b.server.Call(task.ToClient, &server.Request{
		Type: "peripheral",
		Args: []any{
			task.ToBusAddr,
			"pushItems",
			task.ToInvAddr,
			freeSlot + 1,
			task.Count,
			task.ToInvSlot + 1,
		},
	})
	if err != nil {
		log.Error(err)
	}

	//free slot
	b.mutex.Lock()
	if len(b.waiters) > 0 {
		ch := b.waiters[0]
		b.waiters = b.waiters[1:]
		b.allocations[freeSlot] = true // immediately re-allocate to waiter
		ch <- freeSlot
	} else {
		b.allocations[freeSlot] = false
	}
	b.mutex.Unlock()
	return nil
}
