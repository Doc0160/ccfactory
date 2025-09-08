package storage

import (
	"ccfactory/server/factory"
)

type ChestConfig struct {
	Client  string
	InvAddr string
	BusAddr string
}

func (s *ChestConfig) Build(f *factory.Factory) factory.Storage {
	return &Chest{
		config:  s,
		factory: f,
	}
}

type Chest struct {
	config  *ChestConfig
	factory *factory.Factory
}

func (s *Chest) Update() {
	if !s.factory.ClientConnected(s.config.Client) {
		log.Warn("Client missing", "name", s.config.Client)
		return
	}

	/*size := s.factory.Call(s.config.Client, factory.RemoteCall{
		Action: factory.ActionPeripheralCall,
		Data: factory.PeripheralCallData{
			Name:   s.config.InvAddr,
			Method: "size",
		},
	}).Result.([]any)[0].(float64)
	fmt.Println(size)

	list := s.factory.Call(s.config.Client, factory.RemoteCall{
		Action: factory.ActionPeripheralCall,
		Data: factory.PeripheralCallData{
			Name:   s.config.InvAddr,
			Method: "list",
		},
	})
	fmt.Println(list.Result)*/
}
