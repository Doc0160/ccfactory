package storage

import (
	"ccfactory/server/access"
	"ccfactory/server/debug"
	"ccfactory/server/itemdata"
	"ccfactory/server/misc"
	"ccfactory/server/peripheral"
	"ccfactory/server/server"
)

type ChestConfig struct {
	Access access.BusAccess
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
			Access: config.Access,
			Server: server,
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
	defer debug.Timer("chest.Update")()

	var size int
	var sizeError error
	var list []*peripheral.Item
	var listError error

	misc.DoParrallel(func() {
		size, sizeError = s.inventory.Size()
	}, func() {
		list, listError = s.inventory.List()
	})
	if sizeError != nil {
		log.Error(sizeError)
		return sizeError
	}
	if listError != nil {
		log.Error(listError)
		return listError
	}

	for i := 0; i < size; i++ {
		if i >= len(list) {
			break
		}
		if list[i] == nil {
			continue
		}
		pitem := list[i]
		item := itemdata.FromItem(pitem)

		detail := s.detailCache.GetOrSet(item.String(), func() *itemdata.Detail {
			defer debug.Timer("detail")()
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
			//todo : Access:   s.Access.WithSlot(i),
		})
	}

	return nil
}
