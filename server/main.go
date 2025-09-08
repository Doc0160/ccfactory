package main

import (
	"ccfactory/server/factory"
	"ccfactory/server/factory/storage"
	"time"
)

func main() {
	factory.FactoryConfig{
		Port:         "1847",
		MinCycleTime: 5 * time.Second,
	}.Build(func(factory *factory.Factory) {
		factory.AddItemStorage(
			&storage.ChestConfig{
				Client:  "C1",
				InvAddr: "minecraft:barrel_1",
				BusAddr: "dimstorage:dimensional_chest_1",
			})
		/*factory.AddItemStorage(ChestItemStorageConfig{
			Client:        "C1",
			InventoryAddr: "minecraft:barrel_1",
			BusAddr:       "dimstorage:dimensional_chest_1",
		})*/
	})

}
