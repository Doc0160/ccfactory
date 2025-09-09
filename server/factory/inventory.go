package factory

type Inventory interface {
	Access
}

type Access struct {
	GetClient
	GetAddr
}

type GetClient struct {
	client *string
}

func (c *GetClient) GetClient() *string {
	return c.client
}

type GetAddr struct {
	addr *string
}

func (c *GetAddr) GetAddr() *string {
	return c.addr
}
