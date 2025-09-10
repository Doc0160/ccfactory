package process

import "ccfactory/server/itemdata"

type Stock struct {
	Item  itemdata.Filter
	Count int
}

type RestockConfig struct {
	Name  string
	Stock []Stock
}
