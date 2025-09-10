package itemdata

import (
	"sync"
)

type Data struct {
	labelMap map[string][]*Item
	nameMap  map[string][]*Item
	items    sync.Map // map[string]*ItemInfo
}

func New() *Data {
	return &Data{
		labelMap: map[string][]*Item{},
		nameMap:  map[string][]*Item{},
	}
}

func (d *Data) Clear() {
	d.items.Clear()
	d.labelMap = map[string][]*Item{}
	d.nameMap = map[string][]*Item{}
}

func (d *Data) RegisterStoredItem(item *Item, detail *Detail) *ItemInfo {
	label := detail.Label
	name := item.Name
	key := item.String()

	found := false
	for _, i := range d.labelMap[label] {
		if i.String() == item.String() {
			found = true
			break
		}
	}
	if !found {
		d.labelMap[label] = append(d.labelMap[label], item)
	}

	found = false
	for _, i := range d.nameMap[name] {
		if i.String() == item.String() {
			found = true
			break
		}
	}
	if !found {
		d.nameMap[name] = append(d.nameMap[name], item)
	}

	itemVal, loaded := d.items.LoadOrStore(key, &ItemInfo{
		Item:      item,
		Detail:    detail,
		Providers: []*Provider{},
	})

	if !loaded {
		itemVal.(*ItemInfo).Init()
	}

	return itemVal.(*ItemInfo)
}

func (d *Data) SearchItem(filter Filter) *ItemInfo {
	var bestInfo *ItemInfo

	switch f := filter.(type) {
	case LabelFilter:
		if items, ok := d.labelMap[f.Label]; ok {
			for _, key := range items {
				info_gen, _ := d.items.Load(key.String())
				info := info_gen.(*ItemInfo)
				if bestInfo != nil && info.Stored <= bestInfo.Stored {
					continue
				}
				bestInfo = info
			}
		}
	case NameFilter:
		if items, ok := d.nameMap[f.Name]; ok {
			for _, key := range items {
				info_gen, _ := d.items.Load(key.String())
				info := info_gen.(*ItemInfo)
				if bestInfo != nil && info.Stored <= bestInfo.Stored {
					continue
				}
				bestInfo = info
			}
		}
	case BothFilter:
		if items, ok := d.nameMap[f.Name]; ok {
			for _, key := range items {
				if key.Name == f.Name {
					info_gen, _ := d.items.Load(key.String())
					info := info_gen.(*ItemInfo)
					if bestInfo != nil && info.Stored <= bestInfo.Stored {
						continue
					}
					bestInfo = info
				}
			}
		}
	case CustomFiler:
		d.items.Range(func(key, info_gen any) bool {
			info := info_gen.(*ItemInfo)
			if f.Apply(info.Item, info.Detail) {
				if bestInfo != nil && info.Stored <= bestInfo.Stored {
					return true
				}
				bestInfo = info
			}
			return true
		})
	}

	return bestInfo
}
