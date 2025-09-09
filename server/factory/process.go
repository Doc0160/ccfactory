package factory

type ProcessConfig interface {
	Build(*Factory) Process
}

type Process interface {
	Run()
}
