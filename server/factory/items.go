package factory

import (
	"ccfactory/server/heap"
	"fmt"
)

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

type DetailStack struct {
	Item   *Item
	Detail *Detail
	Size   int
}

type Detail struct {
	Label   string
	MaxSize int
}

type ItemInfo struct {
	Detail    *Detail
	Stored    int
	Backup    int
	Providers *heap.Heap[Provider]
}

func (i *ItemInfo) String() string {
	return fmt.Sprint("Detail=", i.Detail,
		"; Stored=", i.Stored,
		"; Providers=", i.Providers.Len())
}

func (i *ItemInfo) Provide(provider *Provider) {
	provided := provider.Provided
	if provided > 0 {
		i.Stored += provided
		i.Providers.Push(*provider)
	}
}

type Provider struct {
	priority int
	Provided int
	Extract  Extractor
}

func NewProvider(priority int, provided int, extractor Extractor) *Provider {
	return &Provider{
		priority: priority,
		Provided: provided,
		Extract:  extractor,
	}
}

type Extractor func(size int, bus_slot int) error

type Filter func(item *Item, detail *Detail) bool

func FilterName(name string) Filter {
	return func(item *Item, detail *Detail) bool {
		return item.Name == name
	}
}

func FilterLabel(label string) Filter {
	return func(item *Item, detail *Detail) bool {
		return detail.Label == label
	}
}

func FilterBoth(name string, label string) Filter {
	return func(item *Item, detail *Detail) bool {
		return item.Name == name && detail.Label == label
	}
}
