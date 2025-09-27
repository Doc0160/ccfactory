package factory

import (
	"ccfactory/server/itemdata"
	"ccfactory/server/misc"
	"ccfactory/server/peripheral"
	"ccfactory/server/process"
	"ccfactory/server/server"
	"ccfactory/server/storage"
	"fmt"
	"time"
)

type Factory struct {
	*FactoryConfig
	Bus *Bus

	ItemStorages []*storage.Chest
	Processes    []process.Process

	ItemData *itemdata.Data

	BusTasks []peripheral.BusTask

	cycle int
	start time.Time
}

func (f *Factory) AddProcess(config process.ProcessConfig) {
	f.Processes = append(f.Processes, config.IntoProcess(f.Server, f.ItemData, f.DetailCache))
}

func (f *Factory) AddItemStorage(config *storage.ChestConfig) {
	f.ItemStorages = append(f.ItemStorages, config.IntoChest(f.Server, f.ItemData, f.DetailCache))
}

func (factory *Factory) Cycle() {
	factory.StartOffCycle()

	p := misc.NewParrallel()
	for _, storage := range factory.ItemStorages {
		p.Add(func() {
			err := storage.Update()
			if err != nil {
				log.Error(err)
			}
		})
	}
	p.Wait()

	for _, process := range factory.Processes {
		err := process.Run()
		if err != nil {
			log.Error(err)
		}
	}

	factory.Bus.Update()

	for _, task := range factory.BusTasks {
		err := factory.Bus.Transfer(task)
		if err != nil {
			log.Error(err)
		}
	}

	//info := factory.ItemData.SearchItem(itemdata.NameFilter{Name: "minecraft:torch"})
	//log.Debug(info)

	factory.EndOfCycle()
}

func (factory *Factory) StartOffCycle() {
	factory.start = time.Now()
	log.Info("Cycle started")
}

func (factory *Factory) EndOfCycle() {
	factory.ItemData.Clear()

	duration := time.Since(factory.start)
	factory.Log(fmt.Sprintf("CCFactory #%d, cycle=%s", factory.cycle, duration), 1)
	factory.cycle++
}

func (f *Factory) Log(text string, color int) {
	if f.LogClients != nil {
		for _, c := range f.LogClients {
			f.Server.Call(c, &server.Request{
				Type: "log",
				Args: []any{struct {
					Text  string `json:"text"`
					Color int    `json:"color"`
				}{text, color}},
			})
		}
	}
	c := "\x1b[30m"
	switch color {
	case 0: // white
		c = "\x1b[97m"
	case 1: // orange
		c = "\x1b[38;5;202m"
	case 2: // magenta
		c = "\x1b[35m"
	case 3: // lightBlue
		c = "\x1b[36m"
	case 4: // yellow
		c = "\x1b[33m"
	case 5: // lime
		c = "\x1b[32m"
	case 6: // pink
		c = "\x1b[95m"
	case 7: // gray
		c = "\x1b[90m"
	case 8: // lightGray
		c = "\x1b[37m"
	case 9: // cyan
		c = "\x1b[36m"
	case 10: // purple
		c = "\x1b[35m"
	case 11: // blue
		c = "\x1b[34m"
	case 12: // brown
		c = "\x1b[33m"
	case 13: // green
		c = "\x1b[32m"
	case 14: // red
		c = "\x1b[31m"
	case 15: // black
		c = "\x1b[30m"
	}
	logFactory.Info(c + text + "\033[0m")
}
