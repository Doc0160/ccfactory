package process

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
	p.factory.Log("test", 6)
	//log.Debug(p.factory.SearchItem(name("minecraft:torch")))
}
