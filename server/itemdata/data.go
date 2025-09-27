package itemdata

import (
	"ccfactory/server/debug"
	"sync"
)

type Data struct {
	labelMap *ItemSliceMap //map[string][]*Item
	nameMap  *ItemSliceMap //map[string][]*Item
	items    sync.Map      //map[string]*ItemInfo
}

type ItemSliceMap struct {
	mu   sync.RWMutex
	data map[string][]*Item
}

func NewItemSliceMap() *ItemSliceMap {
	return &ItemSliceMap{
		data: make(map[string][]*Item),
	}
}
func (m *ItemSliceMap) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = map[string][]*Item{}
}
func (s *ItemSliceMap) Get(key string) []*Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[key]
}
func (s *ItemSliceMap) Append(key string, item *Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = append(s.data[key], item)
}

func New() *Data {
	return &Data{
		labelMap: NewItemSliceMap(),
		nameMap:  NewItemSliceMap(),
	}
}

func (d *Data) Clear() {
	d.items.Clear()
	d.labelMap.Clear()
	d.nameMap.Clear()
}

func (d *Data) RegisterStoredItem(item *Item, detail *Detail) *ItemInfo {
	defer debug.Timer("RegisterStoredItem")()

	label := detail.Label
	name := item.Name
	key := item.String()

	found := false
	for _, i := range d.labelMap.Get(label) {
		if i.String() == item.String() {
			found = true
			break
		}
	}
	if !found {
		d.labelMap.Append(label, item)
	}

	found = false
	for _, i := range d.nameMap.Get(name) {
		if i.String() == item.String() {
			found = true
			break
		}
	}
	if !found {
		d.nameMap.Append(name, item)
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
	defer debug.Timer("SearchItem")()

	var bestInfo *ItemInfo
	switch f := filter.(type) {
	case LabelFilter:
		if items := d.labelMap.Get(f.Label); items != nil {
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
		if items := d.nameMap.Get(f.Name); items != nil {
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
		if items := d.nameMap.Get(f.Name); items != nil {
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
