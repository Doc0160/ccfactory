package main

import (
	"ccfactory/server/factory"
	"ccfactory/server/factory/process"
	"ccfactory/server/factory/storage"
)

func main() {
	(&factory.FactoryConfig{
		LogClients: []string{"C0"},
	}).Build(func(f *factory.Factory) {

		f.AddStorage(&storage.ChestConfig{
			BusAccess: factory.BusAccess{
				Client:  "C0",
				InvAddr: "minecraft:barrel_0",
				BusAddr: "dimstorage:dimensional_chest_0",
			},
		})

		f.AddStorage(&storage.ChestConfig{
			BusAccess: factory.BusAccess{
				Client:  "C0",
				InvAddr: "minecraft:barrel_3",
				BusAddr: "dimstorage:dimensional_chest_0",
			},
		})

		f.AddProcess(&process.StockConfig{
			BusAccess: factory.BusAccess{
				Client:  "C0",
				InvAddr: "minecraft:barrel_1",
				BusAddr: "dimstorage:dimensional_chest_0",
			},
		})

	})
}
