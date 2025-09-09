package factory

func (f *Factory) Call(addr string, action Action, args ...any) RemoteCallResult {
	f.newIdMutex.Lock()
	id := f.nextId
	f.nextId++
	f.newIdMutex.Unlock()

	call := RemoteCall{
		ID:     id,
		Action: action,
		Args:   args,
	}

	respCh := make(chan RemoteCallResult)
	f.responseChannels[id] = respCh

	f.connexions[addr].WriteJSON(call)

	resp := <-respCh
	delete(f.responseChannels, id)
	return resp
}

type Login struct {
	Addr string `json:"addr"`
}

type RemoteCall struct {
	ID int `json:"id"`
	// See Action*
	Action Action `json:"action"`
	Args   []any  `json:"args"`
}

type RemoteCallResult struct {
	ID     int `json:"id"`
	Result any `json:"result,omitempty"`
	Error  any `json:"error,omitempty"`
}

type Action string

const (
	ActionLog            Action = "log"
	ActionPeripheralCall        = "peripheral"
)
