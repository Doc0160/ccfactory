package process

import itemdata "ccfactory/server/itemData"

type Stock struct {
	Item  itemdata.Filter
	Count int
}

type RestockConfig struct {
	Name  string
	Stock []Stock
}
