package factory

import (
	"ccfactory/server/itemdata"
	"ccfactory/server/server"
	"ccfactory/server/storage"
	"time"
)

type FactoryConfig struct {
	Server     *server.Server
	LogClients []string
}

func (config *FactoryConfig) Build(builder func(*Factory)) {
	factory := &Factory{
		FactoryConfig: config,
		ItemStorages:  []*storage.Chest{},

		ItemData: itemdata.New(),

		cycle: -1,
	}

	builder(factory)

	for factory.Server.NumClients() == 0 {
		time.Sleep(1 * time.Second)
	}
	factory.Cycle()
}
