package misc

import "sync"

type Parrallel struct {
	wg sync.WaitGroup
}

func (p *Parrallel) Add(fns ...func()) {
	p.wg.Add(len(fns))
	for _, fn := range fns {
		go func() {
			defer p.wg.Done()
			fn()
		}()
	}
}

func (p *Parrallel) Wait() {
	p.wg.Wait()
}
func NewParrallel() *Parrallel {
	return &Parrallel{}
}
func DoParrallel(fns ...func()) {
	p := Parrallel{}
	p.Add(fns...)
	p.Wait()
}
