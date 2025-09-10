package storage

import (
	"ccfactory/server/itemdata"
	"ccfactory/server/peripheral"
	"ccfactory/server/server"
)

type ChestConfig struct {
	Access peripheral.BusAccess
}

func (config *ChestConfig) IntoChest(
	server *server.Server,
	itemdata *itemdata.Data,
	detailCache *itemdata.DetailCache) *Chest {
	return &Chest{
		ChestConfig:    config,
		globalItemData: itemdata,
		detailCache:    detailCache,
		inventory: &peripheral.Inventory{
			BusAccess: config.Access,
			Server:    server,
		},
	}
}

type Chest struct {
	*ChestConfig
	inventory      *peripheral.Inventory
	globalItemData *itemdata.Data
	detailCache    *itemdata.DetailCache
}

type ItemDetail struct {
	Item     itemdata.Item
	Detail   itemdata.Detail
	Provider itemdata.Provider
}

func (s *Chest) Update() error {
	size, err := s.inventory.Size()
	if err != nil {
		log.Error(err)
		return err
	}

	list, err := s.inventory.List()
	if err != nil {
		log.Error(err)
		return err
	}

	for i := 0; i < size; i++ {
		if i >= len(list) {
			continue
		}
		if list[i] == nil {
			continue
		}
		pitem := list[i]
		item := &itemdata.Item{
			Name:    pitem.Name,
			NbtHash: pitem.Nbt,
		}

		detail := s.detailCache.Get(item.String(), func() *itemdata.Detail {
			pdetail, err := s.inventory.GetItemDetail(i)
			if err != nil {
				log.Error(err)
				return nil
			}

			detail := &itemdata.Detail{
				Label:   pdetail.Label,
				MaxSize: pdetail.MaxCount,
				Other: itemdata.DetailOthers{
					Tags:         pdetail.Tags,
					Damage:       pdetail.Damage,
					MaxDamage:    pdetail.MaxDamage,
					Durability:   pdetail.Durability,
					Enchantments: []itemdata.Enchantment{},
				},
			}

			for _, e := range pdetail.Enchantments {
				detail.Other.Enchantments = append(detail.Other.Enchantments, itemdata.Enchantment{
					Name:  e.Name,
					Label: e.Label,
					Level: e.Level,
				})
			}

			return detail
		})

		s.globalItemData.RegisterStoredItem(item, detail).Provide(&itemdata.Provider{
			Priority: pitem.Count,
			Provided: pitem.Count,
			Access:   s.Access.WithSlot(i),
		})
	}

	return nil
}
