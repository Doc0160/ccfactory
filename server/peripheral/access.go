package peripheral

type BusAccess struct {
	Client  string
	BusAddr string
	InvAddr string
}

func (a BusAccess) WithSlot(slot int) BusAccessWithSlot {
	return BusAccessWithSlot{
		Client:  a.Client,
		BusAddr: a.BusAddr,
		InvAddr: a.InvAddr,
		InvSlot: slot,
	}
}

type BusAccessWithSlot struct {
	Client  string
	BusAddr string
	InvAddr string
	InvSlot int
}

func NewBusTask(from *BusAccessWithSlot, to *BusAccessWithSlot) *BusTask {
	return &BusTask{
		FromClient:  from.Client,
		FromBusAddr: from.BusAddr,
		FromInvAddr: from.InvAddr,
		FromInvSlot: from.InvSlot,

		ToClient:  to.Client,
		ToBusAddr: to.BusAddr,
		ToInvAddr: to.InvAddr,
		ToInvSlot: to.InvSlot,
	}
}

type BusTask struct {
	FromClient  string
	FromBusAddr string
	FromInvAddr string
	FromInvSlot int

	ToClient  string
	ToBusAddr string
	ToInvAddr string
	ToInvSlot int
}
