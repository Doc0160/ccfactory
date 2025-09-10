package main

import (
	"ccfactory/server/factory"
	"ccfactory/server/itemdata"
	"ccfactory/server/peripheral"
	"ccfactory/server/server"
	"ccfactory/server/storage"
)

func main() {
	(&factory.FactoryConfig{
		Server:      server.NewServer(1847),
		LogClients:  []string{"C0"},
		DetailCache: itemdata.NewDetailCache("detail_cache.gob"),
	}).Build(func(f *factory.Factory) {

		f.AddItemStorage(&storage.ChestConfig{
			Access: peripheral.BusAccess{
				Client:  "C0",
				InvAddr: "minecraft:barrel_0",
				BusAddr: "dimstorage:dimensional_chest_0",
			},
		})

		f.AddItemStorage(&storage.ChestConfig{
			Access: peripheral.BusAccess{
				Client:  "C2",
				InvAddr: "minecraft:barrel_4",
				BusAddr: "dimstorage:dimensional_chest_2",
			},
		})
		f.AddItemStorage(&storage.ChestConfig{
			Access: peripheral.BusAccess{
				Client:  "C2",
				InvAddr: "minecraft:barrel_5",
				BusAddr: "dimstorage:dimensional_chest_2",
			},
		})
		f.AddItemStorage(&storage.ChestConfig{
			Access: peripheral.BusAccess{
				Client:  "C2",
				InvAddr: "minecraft:barrel_6",
				BusAddr: "dimstorage:dimensional_chest_2",
			},
		})
		f.AddItemStorage(&storage.ChestConfig{
			Access: peripheral.BusAccess{
				Client:  "C2",
				InvAddr: "minecraft:barrel_7",
				BusAddr: "dimstorage:dimensional_chest_2",
			},
		})

		/*f.AddStorage(&storage.ChestConfig{
			BusAccess: BusAccess{
				Client:  "C0",
				InvAddr: "minecraft:barrel_0",
				BusAddr: "dimstorage:dimensional_chest_0",
			},
		})*/
		/*
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
			})*/

	})
}
