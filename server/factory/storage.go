package factory

type StorageConfig interface {
	Build(*Factory) Storage
}

type Storage interface {
	Update()
	Deposit(stack *DetailStack, busSlot int)
}
