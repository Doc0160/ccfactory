package itemdata

import (
	"ccfactory/server/peripheral"
	"container/heap"
)

type Item struct {
	Name    string
	NbtHash string
}

func FromItem(i *peripheral.Item) *Item {
	return &Item{
		Name:    i.Name,
		NbtHash: i.Nbt,
	}
}

func (i *Item) String() string {
	if i.NbtHash == "" {
		return i.Name
	}
	return i.Name + ":" + i.NbtHash
}

type Detail struct {
	Label   string
	MaxSize int
	Other   DetailOthers
}
type DetailOthers struct {
	Tags         map[string]bool
	Damage       int
	MaxDamage    int
	Durability   float64
	Enchantments []Enchantment
}
type Enchantment struct {
	Name  string
	Label string
	Level int
}

type ItemInfo struct {
	Item      *Item
	Detail    *Detail
	Stored    int
	Providers Providers
}

func (i *ItemInfo) Init() {
	heap.Init(&i.Providers)
}

func (i *ItemInfo) Provide(p *Provider) {
	i.Stored += p.Provided
	heap.Push(&i.Providers, p)
}

type Provider struct {
	Priority int
	Provided int
	Access   peripheral.BusAccessWithSlot
}

type Providers []*Provider

func (h Providers) Len() int {
	return len(h)
}
func (h Providers) Less(i, j int) bool {
	return h[i].Priority < h[j].Priority
}
func (h Providers) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *Providers) Push(x any) {
	*h = append(*h, x.(*Provider))
}
func (h *Providers) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
func (i *Providers) Fix() {
	heap.Fix(i, 0)
}
func (h Providers) Peek() *Provider {
	if len(h) == 0 {
		return nil
	}
	return h[0]
}

/*

func (pq *PriorityQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}
*/

var _ heap.Interface = (*Providers)(nil)
