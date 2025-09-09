package factory

import (
	"errors"
	"time"

	"github.com/gorilla/websocket"
)

type FactoryConfig struct {
	Port       string
	LogClients []string
}

func (c *FactoryConfig) newFactory() *Factory {
	return &Factory{
		FactoryConfig: c,

		nameMap:  map[string][]*Item{},
		labelMap: map[string][]*Item{},
	}
}

type Factory struct {
	*FactoryConfig
	nextId int

	conns     map[string]*websocket.Conn
	respChans map[int]chan Response

	item_storage []Storage

	processes []Process

	items    map[*Item]*ItemInfo
	nameMap  map[string][]*Item
	labelMap map[string][]*Item
}

func (f Factory) Log(text string, color int) {
	if f.LogClients != nil {
		for _, c := range f.LogClients {
			f.LogMessage(c, text, color)
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

func (f *Factory) ClientConnected(conn string) bool {
	_, ok := f.conns[conn]
	return ok
}

func (f *Factory) RegisterStoredItem(item *Item, detail *Detail) *ItemInfo {
	if _, ok := f.nameMap[item.Name]; !ok {
		f.nameMap[item.Name] = []*Item{}
	}
	f.nameMap[item.Name] = append(f.nameMap[item.Name], item)

	if _, ok := f.labelMap[detail.Label]; !ok {
		f.labelMap[detail.Label] = []*Item{}
	}
	f.labelMap[detail.Label] = append(f.labelMap[detail.Label], item)

	f.items[item] = &ItemInfo{
		Detail: detail,
		Stored: 0,
		Backup: 0,
	}
	return f.items[item]
}

func (f *Factory) AddStorage(c StorageConfig) {
	f.item_storage = append(f.item_storage, c.Build(f))
}
func (f *Factory) AddProcess(c ProcessConfig) {
	f.processes = append(f.processes, c.Build(f))
}

func (c *FactoryConfig) Build(fn func(*Factory)) {
	f := c.newFactory()

	fn(f)

	go f.StartServer()

	for {
		// update inv
		for _, inv := range f.item_storage {
			inv.Update()
		}

		// run processes
		for _, proc := range f.processes {
			proc.Run()
		}

		// run processes

		f.EndOfCycle()
		time.Sleep(20 * time.Second)
	}
}

func (f *Factory) EndOfCycle() {
	f.Log("Cycle ran", 1)

	log.Debug("", "oak_log", f.nameMap["minecraft:oak_log"], "Oak Log", f.labelMap["Oak Log"])

	log.Info(f.items)

	f.labelMap = map[string][]*Item{}
	f.nameMap = map[string][]*Item{}
	f.items = map[*Item]*ItemInfo{}
}

func (f *Factory) LogMessage(conn string, str string, color int) *Response {
	if conn == "" {
		return nil
	}
	if _, ok := f.conns[conn]; !ok {
		return nil
	}

	log.Debug("LogMessage")
	id := f.nextId
	f.nextId++

	respCh := make(chan Response)
	f.respChans[id] = respCh

	f.conns[conn].WriteJSON(&Request{
		Id:   id,
		Type: "log",
		Args: []any{struct {
			Text  string `json:"text"`
			Color int    `json:"color"`
		}{str, color}},
	})

	resp := <-respCh
	delete(f.respChans, id)
	return &resp
}

func (f *Factory) CallPeripheral(conn string, args ...any) (RawMessage, error) {
	id := f.nextId
	f.nextId++

	respCh := make(chan Response)
	f.respChans[id] = respCh

	f.conns[conn].WriteJSON(&Request{
		Id:   id,
		Type: "peripheral",
		Args: args,
	})

	resp := <-respCh
	delete(f.respChans, id)
	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}
	return resp.Result, nil
}
