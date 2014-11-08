package main

type Throttler struct {
	c chan struct{}
}

func NewThrottler(size int) *Throttler {
	return &Throttler{c: make(chan struct{}, size)}
}

func (t *Throttler) Acquire() {
	t.c <- struct{}{}
}

func (t *Throttler) Release() {
	<-t.c
}
