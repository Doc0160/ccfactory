package factory

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type FactoryConfig struct {
	Port         string
	MinCycleTime time.Duration
}

type Factory struct {
	config FactoryConfig
	//
	itemStorage []Storage
	// item storages
	// fluid storage
	// energy storage
	//processes
	// details cache ; nbt hash -> item details

	// server stuff
	nextId           int
	newIdMutex       sync.Mutex
	responseChannels map[int]chan RemoteCallResult
	connexions       map[string]*websocket.Conn
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (c FactoryConfig) Build(fn func(*Factory)) {
	if c.Port == "" {
		c.Port = "1847"
	}
	if c.MinCycleTime < time.Second {
		c.MinCycleTime = time.Second
	}

	f := &Factory{
		config:           c,
		newIdMutex:       sync.Mutex{},
		responseChannels: map[int]chan RemoteCallResult{},
		connexions:       map[string]*websocket.Conn{},
		nextId:           1,
		//
		itemStorage: []Storage{},
	}

	fn(f)

	fs := http.FileServer(http.Dir("../client"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		log.Info("Possible new connexion")

		l := Login{}
		err = conn.ReadJSON(&l)
		if err != nil {
			return
		}
		f.connexions[l.Addr] = conn
		log.Info("New connexion", "addr", l.Addr)

		for {
			r := RemoteCallResult{}
			err := conn.ReadJSON(&r)
			if err != nil {
				log.Error(err)
				break
			}
			log.Debug("", "r", r)

			if ch, ok := f.responseChannels[r.ID]; ok {
				ch <- r
			}
		}
	})

	// cycle
	go func() {
		cycles := 0
		for {
			f.UpdateStorage()
			//f.RunProcesses()

			// take bus tasks and do it till 0
			// todo: do
			cycles++
			time.Sleep(f.config.MinCycleTime)
		}
		/*for _, storage := f.itemStorage {
			storage.Update()
		}*/
	}()

	// websocket
	log.Info("Listening on http://localhost:" + c.Port)
	http.ListenAndServe(":"+c.Port, nil)
}

func (f *Factory) AddItemStorage(storage StorageConfig) {
	f.itemStorage = append(f.itemStorage, storage.Build(f))
}
func (f *Factory) ClientConnected(client string) bool {
	_, ok := f.connexions[client]
	return ok
}

func (f *Factory) UpdateStorage() {
	// item
	// foreach storage .Update()
	for _, storage := range f.itemStorage {
		storage.Update()
	}

	// fluid

	// energy
}

// handle bus
// loop
// - update storages
// - run processes
