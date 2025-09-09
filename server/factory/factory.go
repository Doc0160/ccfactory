package factory

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gorilla/websocket"
)

type FactoryConfig struct {
	Port string
}

func (c *FactoryConfig) newFactory() *Factory {
	return &Factory{
		config: c,

		nameMap:  map[string][]*Item{},
		labelMap: map[string][]*Item{},
	}
}

type Factory struct {
	config *FactoryConfig
	nextId int

	conns     map[string]*websocket.Conn
	respChans map[int]chan Response

	item_storage []Storage

	items    map[*Item]*ItemInfo
	nameMap  map[string][]*Item
	labelMap map[string][]*Item
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

		f.EndOfCycle()
		time.Sleep(20 * time.Second)
	}
}

func (f *Factory) EndOfCycle() {
	log.Debug("", "oak_log", f.nameMap["minecraft:oak_log"], "Oak Log", f.labelMap["Oak Log"])

	log.Info(f.items)

	f.labelMap = map[string][]*Item{}
	f.nameMap = map[string][]*Item{}
	f.items = map[*Item]*ItemInfo{}
}

func (f *Factory) LogMessage(conn string, str string, color int) *Response {
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

type RawMessage = json.RawMessage

/*func (r CallPeripheralResult) IntoInt() (int, error) {
	var i int
	err := json.Unmarshal(r, &i)
	return i, err
}*/

/*func (r CallPeripheralResult) Into(v any) error {
	return json.Unmarshal(r, v)
}*/

func Into[T any](r RawMessage) (T, error) {
	var v T
	err := json.Unmarshal(r, &v)
	return v, err
}
