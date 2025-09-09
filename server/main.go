package main

import (
	"ccfactory/server/factory"
	"ccfactory/server/factory/storage"
)

func main() {
	(&factory.FactoryConfig{}).Build(func(f *factory.Factory) {

		f.AddStorage(&storage.ChestConfig{
			Client:  "C0",
			InvAddr: "minecraft:barrel_0",
			BusAddr: "dimstorage:dimensional_chest_0",
		})

		f.AddStorage(&storage.ChestConfig{
			Client:  "C0",
			InvAddr: "minecraft:barrel_1",
			BusAddr: "dimstorage:dimensional_chest_0",
		})

	})
}
