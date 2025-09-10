package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	Logins    map[string]*Client
	respChans map[int]chan Response
	nextId    int
}

func NewServer(port int) *Server {
	if port <= 0 {
		port = 1847
	}

	server := &Server{
		Logins:    map[string]*Client{},
		respChans: map[int]chan Response{},
	}

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

			//log.Debug("Response ", "error ", r.Error, "result", string(r.Result))
			if ch, ok := server.respChans[response.Id]; ok {
				ch <- response
			}
		}
	})

	log.Info(fmt.Sprintf("Listening on http://localhost:%d", port))
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	return server
}

func (s *Server) NumClients() int {
	return len(s.Logins)
}

func (s *Server) Login(name string, client *Client) {
	client.Login = name
	log.Info("Login", "name", name)
	if old, ok := s.Logins[name]; ok && old != nil {
		log.Warn("Logged in from another address, closing old : " + client.Login)
		old.Login = ""
		old.Websocket.Close()
	}
	s.Logins[name] = client
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
	if ws, ok := s.Logins[client_name]; !ok || ws == nil {
		return nil, errors.New("client not found :" + client_name)

	}

	id := s.nextId
	s.nextId++
	request.Id = id

	respCh := make(chan Response)
	s.respChans[id] = respCh

	s.Logins[client_name].Websocket.WriteJSON(request)

	resp := <-respCh
	delete(s.respChans, id)
	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}

	if len(resp.Result) == 0 {
		return nil, nil
	}

	return resp.Result, nil
}

type Client struct {
	Login     string
	Server    *Server
	Websocket *websocket.Conn
}

func (c *Client) Close() {
	log.Warn("Client closed", "name", c.Login)
	c.Server.Logins[c.Login] = nil
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
