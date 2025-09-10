package storage

import (
	"ccfactory/server/itemdata"
	"ccfactory/server/peripheral"
	"ccfactory/server/server"
)

type ChestConfig struct {
	Access peripheral.BusAccess
}

func (config *ChestConfig) IntoChest(server *server.Server, itemdata *itemdata.Data) *Chest {
	return &Chest{
		ChestConfig:    config,
		GlobalItemData: itemdata,
		inventory: &peripheral.Inventory{
			BusAccess: config.Access,
			Server:    server,
		},
	}
}

type Chest struct {
	*ChestConfig
	inventory      *peripheral.Inventory
	GlobalItemData *itemdata.Data
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
		item := list[i]

		itemDetail, err := s.inventory.GetItemDetail(i)
		if err != nil {
			log.Error(err)
			continue
		}

		detail := itemdata.Detail{
			Label:   itemDetail.Label,
			MaxSize: itemDetail.MaxCount,
			Other: itemdata.DetailOthers{
				Tags:         itemDetail.Tags,
				Damage:       itemDetail.Damage,
				MaxDamage:    itemDetail.MaxDamage,
				Durability:   itemDetail.Durability,
				Enchantments: []itemdata.Enchantment{},
			},
		}
		for _, e := range itemDetail.Enchantments {
			detail.Other.Enchantments = append(detail.Other.Enchantments, itemdata.Enchantment{
				Name:  e.Name,
				Label: e.Label,
				Level: e.Level,
			})
		}

		s.GlobalItemData.RegisterStoredItem(&itemdata.Item{
			Name:    item.Name,
			NbtHash: item.Nbt,
		}, &detail).Provide(&itemdata.Provider{
			Priority: item.Count,
			Provided: item.Count,
			Access:   s.Access.WithSlot(i),
		})
	}

	return nil
}

/*
func (s *Chest) List() ([]ItemDetail, error) {
	size, err := s.inventory.Size()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//og.Debug("", "size", size)

	list, err := s.inventory.List()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug("", "list", list)

	info := []ItemDetail{}
	for i := 0; i < size; i++ {
		if list[i] == nil {
			continue
		}
		item := list[i]

		itemDetail, err := s.inventory.GetItemDetail(i)
		if err != nil {
			log.Error(err)
			continue
		}

		log.Debug("", "detail", itemDetail)
		detail := itemdata.Detail{
			Label:   itemDetail.Label,
			MaxSize: itemDetail.MaxCount,
			Other: itemdata.DetailOthers{
				Tags:         itemDetail.Tags,
				Damage:       itemDetail.Damage,
				MaxDamage:    itemDetail.MaxDamage,
				Durability:   itemDetail.Durability,
				Enchantments: []itemdata.Enchantment{},
			},
		}
		for _, e := range itemDetail.Enchantments {
			detail.Other.Enchantments = append(detail.Other.Enchantments, itemdata.Enchantment{
				Name:  e.Name,
				Label: e.Label,
				Level: e.Level,
			})
		}

		info = append(info, ItemDetail{
			Item: itemdata.Item{
				Name:    item.Name,
				NbtHash: item.Nbt,
			},
			Detail: detail,
			Provider: itemdata.Provider{
				Priority: item.Count,
				Provided: item.Count,
				Access:   s.Access.WithSlot(i),
			},
		})
	}

	return info, nil
}*/
