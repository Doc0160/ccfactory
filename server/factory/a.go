package factory

type RemoteCallResult struct {
	ID     int   `json:"id"`
	Result []any `json:"result"`
	Error  any   `json:"error"`
}
type RemoteCall struct {
	// See Action*
	ID     int    `json:"id"`
	Action Action `json:"action"`
	Args   []any  `json:"args"`
}

type Action string

const (
	ActionLog            Action = "log"
	ActionPeripheralCall        = "peripheral"
)

type Login struct {
	Addr string `json:"addr"`
}
