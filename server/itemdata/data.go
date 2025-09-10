package itemdata

import (
	"ccfactory/server/debug"
)

type Data struct {
	labelMap map[string][]*Item
	nameMap  map[string][]*Item
	items    map[string]*ItemInfo
}

func New() *Data {
	return &Data{
		labelMap: map[string][]*Item{},
		nameMap:  map[string][]*Item{},
		items:    map[string]*ItemInfo{},
	}
}

func (d *Data) Clear() {
	debug.Dump("items", d.items)
	debug.Dump("itemsLabelMap", d.labelMap)
	debug.Dump("itemsNameMap", d.nameMap)

	d.items = map[string]*ItemInfo{}
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

	if _, ok := d.items[key]; !ok {
		d.items[key] = &ItemInfo{
			Item:      item,
			Detail:    detail,
			Providers: []*Provider{},
		}
		d.items[key].Init()
	}

	return d.items[key]
}

func (d *Data) SearchItem(filter Filter) *ItemInfo {
	var bestInfo *ItemInfo

	switch f := filter.(type) {
	case LabelFilter:
		if items, ok := d.labelMap[f.Label]; ok {
			for _, key := range items {
				info := d.items[key.String()]
				if bestInfo != nil && info.Stored <= bestInfo.Stored {
					continue
				}
				bestInfo = info
			}
		}
	case NameFilter:
		if items, ok := d.nameMap[f.Name]; ok {
			for _, key := range items {
				info := d.items[key.String()]
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
					info := d.items[key.String()]
					if bestInfo != nil && info.Stored <= bestInfo.Stored {
						continue
					}
					bestInfo = info
				}
			}
		}
	case CustomFiler:
		for _, info := range d.items {
			if f.Apply(info.Item, info.Detail) {
				if bestInfo != nil && info.Stored <= bestInfo.Stored {
					continue
				}
				bestInfo = info
			}
		}
	}

	return bestInfo
}
