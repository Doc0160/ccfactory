package access

type Client = string
type Addr = string

type GetClient interface {
	GetClient() Client
}
type GetAddr interface {
	GetAddr() Addr
}

type BasicAccessInterface interface {
	GetClient
	GetAddr
}

type BasicAccess struct {
	Client Client
	Addr   Addr
}

var _ BasicAccessInterface = (*BasicAccess)(nil)

func (a BasicAccess) GetClient() Client { return a.Client }
func (a BasicAccess) GetAddr() Addr     { return a.Addr }

type BusAccess struct {
	Client  Client
	InvAddr Addr
	BusAddr Addr
}

var _ BasicAccessInterface = (*BusAccess)(nil)

func (a BusAccess) GetClient() Client { return a.Client }
func (a BusAccess) GetAddr() Addr     { return a.InvAddr }
