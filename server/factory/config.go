package factory

import (
	"ccfactory/server/access"
	"ccfactory/server/itemdata"
	"ccfactory/server/server"
	"ccfactory/server/storage"
	"time"
)

type FactoryConfig struct {
	Server      *server.Server
	DetailCache *itemdata.DetailCache
	LogClients  []string

	BusAccess access.BasicAccess
}

func (config *FactoryConfig) Build(builder func(*Factory)) {
	factory := &Factory{
		FactoryConfig: config,

		ItemStorages: []*storage.Chest{},

		ItemData: itemdata.New(),

		cycle: 0,
	}
	bus := NewBus(config.Server, config.BusAccess, factory)
	factory.Bus = bus

	builder(factory)

	for factory.Server.NumClients() < 2 {
		time.Sleep(1 * time.Second)
	}
	factory.Cycle()
	/*
		misc.DoParrallel(func() {
			factory.Bus.Transfer(peripheral.BusTask{
				FromClient:  "C2",
				FromBusAddr: "dimstorage:dimensional_chest_2",
				FromInvAddr: "minecraft:barrel_4",
				FromInvSlot: 0,

				Count: 1,

				ToClient:  "C0",
				ToBusAddr: "dimstorage:dimensional_chest_0",
				ToInvAddr: "minecraft:barrel_1",
				ToInvSlot: 0,
			})
		}, func() {
			factory.Bus.Transfer(peripheral.BusTask{
				FromClient:  "C2",
				FromBusAddr: "dimstorage:dimensional_chest_2",
				FromInvAddr: "minecraft:barrel_4",
				FromInvSlot: 0,

				Count: 1,

				ToClient:  "C0",
				ToBusAddr: "dimstorage:dimensional_chest_0",
				ToInvAddr: "minecraft:barrel_1",
				ToInvSlot: 0,
			})
		}, func() {
			factory.Bus.Transfer(peripheral.BusTask{
				FromClient:  "C2",
				FromBusAddr: "dimstorage:dimensional_chest_2",
				FromInvAddr: "minecraft:barrel_4",
				FromInvSlot: 0,

				Count: 1,

				ToClient:  "C0",
				ToBusAddr: "dimstorage:dimensional_chest_0",
				ToInvAddr: "minecraft:barrel_1",
				ToInvSlot: 0,
			})
		}, func() {
			factory.Bus.Transfer(peripheral.BusTask{
				FromClient:  "C2",
				FromBusAddr: "dimstorage:dimensional_chest_2",
				FromInvAddr: "minecraft:barrel_4",
				FromInvSlot: 0,

				Count: 1,

				ToClient:  "C0",
				ToBusAddr: "dimstorage:dimensional_chest_0",
				ToInvAddr: "minecraft:barrel_1",
				ToInvSlot: 0,
			})
		}, func() {
			factory.Bus.Transfer(peripheral.BusTask{
				FromClient:  "C2",
				FromBusAddr: "dimstorage:dimensional_chest_2",
				FromInvAddr: "minecraft:barrel_4",
				FromInvSlot: 0,

				Count: 1,

				ToClient:  "C0",
				ToBusAddr: "dimstorage:dimensional_chest_0",
				ToInvAddr: "minecraft:barrel_1",
				ToInvSlot: 0,
			})
		}, func() {
			factory.Bus.Transfer(peripheral.BusTask{
				FromClient:  "C2",
				FromBusAddr: "dimstorage:dimensional_chest_2",
				FromInvAddr: "minecraft:barrel_4",
				FromInvSlot: 0,

				Count: 1,

				ToClient:  "C0",
				ToBusAddr: "dimstorage:dimensional_chest_0",
				ToInvAddr: "minecraft:barrel_1",
				ToInvSlot: 0,
			})
		})
	*/
	//time.Sleep(1 * time.Second)
	//factory.Cycle()
}
