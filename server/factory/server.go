package factory

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

type Login struct {
	Client string `json:"client"`
}

type Request struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
	Args []any  `json:"args"`
}

type Response struct {
	Id     int             `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  string          `json:"error"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (f *Factory) StartServer() {
	if f.config.Port == "" {
		f.config.Port = "1847"
	}
	f.conns = map[string]*websocket.Conn{}
	f.respChans = map[int]chan Response{}

	fs := http.FileServer(http.Dir("../client"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		l := Login{}
		err = conn.ReadJSON(&l)
		if err != nil {
			log.Error(err)
			return
		}
		f.conns[l.Client] = conn
		log.Debug("Login ", "client", l.Client)

		for {
			r := Response{}
			err = conn.ReadJSON(&r)
			if err != nil {
				log.Error(err)
				return
			}
			//log.Debug("Response ", "error ", r.Error, "result", string(r.Result))
			if ch, ok := f.respChans[r.Id]; ok {
				ch <- r
			}
		}

	})

	log.Info("Listening on http://localhost:" + f.config.Port)
	http.ListenAndServe(":"+f.config.Port, nil)
}
