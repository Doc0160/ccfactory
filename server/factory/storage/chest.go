package storage

import (
	"ccfactory/server/factory"
)

type ChestConfig struct {
	BusAccess
}

type Chest struct {
	*ChestConfig
	factory *Factory

	size   int
	stacks []*factory.DetailStack
}

var _ Storage = (*Chest)(nil)

func (c *ChestConfig) Build(f *Factory) Storage {
	return &Chest{
		ChestConfig: c,
		factory:     f,

		size:   0,
		stacks: []*factory.DetailStack{},
	}
}

func (c *Chest) Size() (int, error) {
	sizeResult, err := c.factory.CallPeripheral(c.Client, c.InvAddr, "size")
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
	listResult, err := c.factory.CallPeripheral(c.Client, c.InvAddr, "list")
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
	detailResult, err := c.factory.CallPeripheral(c.Client, c.InvAddr, "getItemDetail", slot+1)
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

func (c *Chest) Deposit(stack *factory.DetailStack, busSlot int) {
	_, err := c.factory.CallPeripheral(c.Client,
		c.BusAddr,
		"pushItems",
		c.InvAddr,
		busSlot+1,
		stack.Size,
		0+1)
	if err != nil {
		log.Error(err)
	}
}

func (c *Chest) Update() {
	//skip not connected
	if !c.factory.IsClientConnected(c.Client) {
		c.factory.Log(c.Client+" not connected", 14)
		return
	}

	var err error
	c.size, err = c.Size()
	if err != nil {
		log.Error(err)
	}
	c.stacks = make([]*factory.DetailStack, c.size)

	list, _ := c.List()
	for invSlot := 0; invSlot < c.size; invSlot++ {
		if invSlot >= len(list) {
			continue
		}
		item := list[invSlot]
		if item == nil {
			continue
		}

		detail, err := c.GetItemDetail(invSlot)
		if err != nil {
			log.Error(err)
		}
		if detail == nil {
			continue
		}

		c.stacks[invSlot] = &factory.DetailStack{
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
			RegisterStoredItem(c.stacks[invSlot].Item, c.stacks[invSlot].Detail).
			Provide(factory.NewProvider(
				item.Count,
				item.Count,
				func(size int, bus_slot int) error {
					r, err := c.factory.CallPeripheral(c.Client, c.BusAddr,
						"pullItems",
						c.InvAddr,
						invSlot+1,
						size,
						bus_slot+1)
					if err != nil {
						return err
					}

					i, err := Into[int](r)
					if err != nil {
						return err
					}

					c.stacks[invSlot].Size -= i
					if c.stacks[invSlot].Size <= 0 {
						c.stacks[invSlot] = nil
					}

					return nil
				}))

	}
}
