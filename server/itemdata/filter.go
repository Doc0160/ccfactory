package itemdata

type Filter interface {
	Apply(item *Item, detail *Detail) bool
	String() string
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
