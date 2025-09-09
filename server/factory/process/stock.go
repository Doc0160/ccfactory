package process

import (
	"ccfactory/server/factory"
)

type StockConfig struct {
	BusAccess
}

type Stock struct {
	*StockConfig
	factory *Factory
}

var _ Process = (*Stock)(nil)

func (c *StockConfig) Build(f *Factory) Process {
	return &Stock{
		StockConfig: c,
		factory:     f,
	}
}

func (p *Stock) Run() {
	if !p.factory.IsClientConnected(p.Client) {
		p.factory.Log(p.Client+" not connected", 14)
		return
	}
	p.factory.Log("test", 6)

	p.factory.PullIntoBus(factory.FilterName("minecraft:torch"), 64)

	msg, err := p.factory.CallPeripheral(p.Client,
		p.BusAddr,
		"pushItems",
		p.InvAddr,
		0+1,
		64,
		0+1)
	if err != nil {
		log.Error(err)
	}
	log.Debug(string(msg))
}
