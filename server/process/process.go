package process

import (
	"ccfactory/server/itemdata"
	"ccfactory/server/server"
)

type ProcessConfig interface {
	IntoProcess(*server.Server, *itemdata.Data, *itemdata.DetailCache) Process
}

type Process interface {
	GetName() string
	Run() error
}
