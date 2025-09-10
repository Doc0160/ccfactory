package itemdata

import (
	"ccfactory/server/peripheral"
	"container/heap"
)

type Filter interface {
	Apply(item *Item, detail *Detail) bool
}

type LabelFilter struct {
	Label string
}

var _ Filter = (*LabelFilter)(nil)

func (f LabelFilter) String() string {
	return "Label=" + f.Label
}
func (f LabelFilter) Apply(item *Item, detail *Detail) bool {
	return detail.Label == f.Label
}

type NameFilter struct {
	Name string
}

var _ Filter = (*NameFilter)(nil)

func (f NameFilter) String() string {
	return "Name=" + f.Name
}
func (f NameFilter) Apply(item *Item, detail *Detail) bool {
	return item.Name == f.Name
}

type BothFilter struct {
	Name  string
	Label string
}

var _ Filter = (*BothFilter)(nil)

func (f BothFilter) String() string {
	return "Label=" + f.Label + "&Name=" + f.Name
}
func (f BothFilter) Apply(item *Item, detail *Detail) bool {
	return item.Name == f.Name && detail.Label == f.Label
}

type CustomFiler struct {
	// Filter descption for logging
	Desc     string
	Function func(*Item, *Detail) bool
}

var _ Filter = (*CustomFiler)(nil)

func (f CustomFiler) String() string {
	return f.Desc
}

func (f CustomFiler) Apply(item *Item, detail *Detail) bool {
	return f.Function(item, detail)
}

type Item struct {
	Name    string
	NbtHash string
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

/*

func (pq *PriorityQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}
*/

var _ heap.Interface = (*Providers)(nil)
