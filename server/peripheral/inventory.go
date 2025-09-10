package peripheral

import (
	"ccfactory/server/server"
)

type Into interface {
	IntoInventory() *Inventory
}

type Inventory struct {
	Server    *server.Server
	BusAccess BusAccess
}

func (p *Inventory) Size() (int, error) {
	response, err := p.Server.Call(
		p.BusAccess.Client,
		&server.Request{
			Type: "peripheral",
			Args: []any{
				p.BusAccess.InvAddr,
				"size",
			},
		})
	if err != nil {
		return -1, err
	}

	size, err := server.Into[int](response)
	if err != nil {
		return -1, err
	}

	return size, nil
}

type Item struct {
	Name  string `json:"name"`
	Nbt   string `json:"nbt,omitempty"`
	Count int    `json:"count"`
}

func (p *Inventory) List() ([]*Item, error) {
	response, err := p.Server.Call(
		p.BusAccess.Client,
		&server.Request{
			Type: "peripheral",
			Args: []any{
				p.BusAccess.InvAddr,
				"list",
			},
		})
	if err != nil {
		return nil, err
	}

	list, err := server.Into[[]*Item](response)
	if err != nil {
		return nil, err
	}

	return list, nil
}

type Detail struct {
	Name         string          `json:"name"`
	Label        string          `json:"displayName"`
	Count        int             `json:"count"`
	MaxCount     int             `json:"maxCount"`
	Tags         map[string]bool `json:"tags,omitempty"`
	Damage       int             `json:"damage,omitempty"`
	MaxDamage    int             `json:"maxDamage,omitempty"`
	Durability   float64         `json:"durability,omitempty"`
	Enchantments []Enchantment   `json:"enchantments,omitempty"`
}

type Enchantment struct {
	Name  string `json:"name"`
	Label string `json:"displayName"`
	Level int    `json:"level"`
}

func (p *Inventory) GetItemDetail(slot int) (*Detail, error) {
	response, err := p.Server.Call(
		p.BusAccess.Client,
		&server.Request{
			Type: "peripheral",
			Args: []any{
				p.BusAccess.InvAddr,
				"getItemDetail",
				slot + 1,
			},
		})
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, nil
	}

	detail, err := server.Into[*Detail](response)
	if err != nil {
		return nil, err
	}

	return detail, nil
}

func (p *Inventory) PushItems(toName, fromSlot, limit, toSlot int) (int, error) {

	return -1, nil
}

func (p *Inventory) PullItems(fromName, fromSlot, limit, toSlot int) (int, error) {

	return -1, nil
}
