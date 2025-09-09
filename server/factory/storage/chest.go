package storage

import (
	"ccfactory/server/factory"
)

type ChestConfig struct {
	Client  string
	InvAddr string
	BusAddr string
}

type Chest struct {
	config  *ChestConfig
	factory *factory.Factory

	size   int
	stacks []factory.DetailStack

	//items map[*Item]*Item
}

func (c *ChestConfig) Build(f *factory.Factory) factory.Storage {
	return &Chest{
		config:  c,
		factory: f,

		size:   0,
		stacks: []factory.DetailStack{},
	}
}

func (c *Chest) Size() (int, error) {
	sizeResult, err := c.factory.CallPeripheral(c.config.Client, c.config.InvAddr, "size")
	if err != nil {
		return -1, err
	}

	size, err := Into[int](sizeResult)
	return size, nil
}

type ListItem struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Nbt   string `json:"nbt"`
}

func (c *Chest) List() ([]*ListItem, error) {
	listResult, err := c.factory.CallPeripheral(c.config.Client, c.config.InvAddr, "list")
	if err != nil {
		return nil, err
	}

	list, err := factory.Into[[]*ListItem](listResult)
	if err != nil {
		return nil, err
	}

	return list, nil
}

type ItemDetail struct {
	Name        string          `json:"name"`
	Nbt         string          `json:"nbt,omitempty"`
	DisplayName string          `json:"displayName"`
	Tags        map[string]bool `json:"tags"`
	Count       int             `json:"count"`
	MaxCount    int             `json:"maxCount"`
}

func (c *Chest) GetItemDetail(slot int) (*ItemDetail, error) {
	detailResult, err := c.factory.CallPeripheral(c.config.Client, c.config.InvAddr, "getItemDetail", slot+1)
	if err != nil {
		return nil, err
	}
	if string(detailResult) == "" {
		return nil, nil
	}

	detail, err := Into[*ItemDetail](detailResult)
	if err != nil {
		return nil, err
	}

	return detail, nil
}

func (c *Chest) Update() {
	//skip not connected
	if !c.factory.ClientConnected(c.config.Client) {
		log.Warn("Client not connected", "client", c.config.Client)
		return
	}

	var err error
	c.size, err = c.Size()
	if err != nil {
		log.Error(err)
	}
	c.stacks = make([]factory.DetailStack, c.size)

	list, _ := c.List()
	for i := 0; i < c.size; i++ {
		if i > len(list) {
			continue
		}
		item := list[i]
		if item == nil {
			continue
		}

		detail, err := c.GetItemDetail(i)
		if err != nil {
			log.Error(err)
		}
		if detail == nil {
			continue
		}

		c.stacks[i] = factory.DetailStack{
			Item: &factory.Item{
				Name:    item.Name,
				NbtHash: item.Nbt,
			},
			Detail: &factory.Detail{
				Label:   detail.DisplayName,
				MaxSize: detail.MaxCount,
			},
			Size: item.Count,
		}

		c.factory.
			RegisterStoredItem(c.stacks[i].Item, c.stacks[i].Detail).
			Provide(&factory.Provider{
				Provided: item.Count,
				Priority: -item.Count,
				Extractor: &ChestExtractor{
					chest:   c,
					invSlot: i,
				},
			})

	}

	//log.Debug("", "addr", c.config.InvAddr, "size", size)

	//c.GetItemDetail(0)
}

type ChestExtractor struct {
	chest   *Chest
	invSlot int
}

func (ce ChestExtractor) Extract(size int, bus_slot int) error {
	config := ce.chest.config
	chest := ce.chest
	r, err := chest.factory.CallPeripheral(config.Client, config.BusAddr,
		"pullItems",
		config.InvAddr,
		ce.invSlot+1,
		size,
		bus_slot+1)
	if err != nil {
		return err
	}

	invStack := &chest.stacks[ce.invSlot]
	i, err := Into[int](r)
	if err != nil {
		return err
	}

	invStack.Size -= i
	if invStack.Size <= 0 {
		invStack = nil
	}

	return nil
}
