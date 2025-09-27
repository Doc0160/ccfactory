package peripheral

import "fmt"

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

func NewBusTask(from BusAccessWithSlot, count int, to BusAccessWithSlot) BusTask {
	return BusTask{
		FromClient:  from.Client,
		FromBusAddr: from.BusAddr,
		FromInvAddr: from.InvAddr,
		FromInvSlot: from.InvSlot,

		Count: count,

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

	Count int

	ToClient  string
	ToBusAddr string
	ToInvAddr string
	ToInvSlot int
}

func (t BusTask) String() string {
	str := ""
	str += t.FromClient + "@" + t.FromInvAddr + "[" + fmt.Sprint(t.FromInvSlot) + "]"
	str += " -> " + t.FromClient + "@" + t.FromBusAddr
	if t.FromClient != t.ToClient && t.FromBusAddr != t.ToBusAddr {
		str += "/" + t.ToClient + "@" + t.ToBusAddr
	}
	str += " -> " + t.ToClient + "@" + t.ToInvAddr + "[" + fmt.Sprint(t.ToInvSlot) + "]"
	return str
}
