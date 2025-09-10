package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

type Server struct {
	Logins    sync.Map // map[string]*Client
	respChans sync.Map // map[int]chan Response
	nextId    int64
}

func NewServer(port int) *Server {
	if port <= 0 {
		port = 1847
	}

	server := &Server{}

	fs := http.FileServer(http.Dir("../client"))
	http.Handle("/", fs)

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("Websocket", "error", err)
			return
		}
		client := &Client{
			Login:     "",
			Server:    server,
			Websocket: conn,
		}
		defer client.Close()

		var login Login
		err = conn.ReadJSON(&login)
		if err != nil {
			log.Error("Websocket", "error", err)
			return
		}
		server.Login(login.Client, client)

		for {
			var response Response
			err = conn.ReadJSON(&response)
			if err != nil {
				log.Error(err)
				return
			}

			if chVal, ok := server.respChans.Load(response.Id); ok {
				ch := chVal.(chan Response)
				ch <- response
			}
		}
	})

	log.Info(fmt.Sprintf("Listening on http://localhost:%d", port))
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	return server
}

func (s *Server) NumClients() int {
	count := 0
	s.Logins.Range(func(_, _ any) bool {
		count++
		return true
	})
	return count
}

func (s *Server) Login(name string, client *Client) {
	client.Login = name
	log.Info("Login", "name", name)

	if old, ok := s.Logins.Load(name); ok && old != nil {
		oldClient := old.(*Client)
		log.Warn("Logged in from another address, closing old: " + oldClient.Login)
		oldClient.Login = ""
		oldClient.Close()
	}

	s.Logins.Store(name, client)
}

type RawMessage = json.RawMessage

func Into[T any](r RawMessage) (T, error) {
	var v T
	err := json.Unmarshal(r, &v)
	return v, err
}

func (s *Server) Call(client_name string, request *Request) (RawMessage, error) {
	if client_name == "" {
		return nil, errors.New("client not found :" + client_name)
	}

	wsVal, ok := s.Logins.Load(client_name)
	if !ok || wsVal == nil {
		return nil, errors.New("client not found: " + client_name)
	}
	ws := wsVal.(*Client)

	id := int(atomic.AddInt64(&s.nextId, 1) - 1)
	respCh := make(chan Response, 1)
	s.respChans.Store(id, respCh)

	request.Id = id
	ws.WriteJSON(request)

	resp := <-respCh

	s.respChans.Delete(id)

	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}
	if len(resp.Result) == 0 || string(resp.Result) == "{}" {
		return nil, nil
	}

	return resp.Result, nil
}

type Client struct {
	Login     string
	Server    *Server
	Websocket *websocket.Conn
	mutex     sync.Mutex
}

func (c *Client) WriteJSON(v interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Websocket.WriteJSON(v)
}

func (c *Client) Close() {
	log.Warn("Client closed", "name", c.Login)
	c.Server.Logins.Delete(c.Login)
	c.Websocket.Close()
}

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
