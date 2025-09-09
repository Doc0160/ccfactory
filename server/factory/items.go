package factory

type Item struct {
	Name    string
	NbtHash string
}

func (i *Item) String() string {
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
	Providers BinaryHeap[Provider]
}

func (i *ItemInfo) Provide(provider *Provider) {
	provided := provider.Provided
	if provided > 0 {
		i.Stored += provided
		i.Providers.Push(*provider)
		provider.Extractor.Extract(1, 0)
	}
}

type Provider struct {
	Priority  int
	Provided  int
	Extractor Extractor
}

func (p Provider) Less(other Less) bool {
	return p.Priority < other.(Provider).Priority
}

type Extractor interface {
	Extract(size int, bus_slot int) error
}

/*
type Item struct {
	Name    string
	NbtHash string
}

type ItemInfo struct {
	Detail  *Detail
	NStored int32
	NBackup int32
}

type Detail struct {
	Label   string
	MaxSize int32
}
*/
