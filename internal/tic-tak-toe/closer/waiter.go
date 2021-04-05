package closer

import (
	"fmt"
	"sync"
)

type CloseFunc func() error

type Closer struct {
	m    sync.Mutex
	once sync.Once

	closers []CloseFunc
}

func NewCloser() *Closer {
	return &Closer{
		closers: []CloseFunc{},
	}
}

func (w *Closer) Add(c CloseFunc) {
	w.m.Lock()
	w.closers = append(w.closers, c)
	w.m.Unlock()
}

func (w *Closer) Close() {
	w.m.Lock()
	closers := w.closers
	w.closers = nil
	w.m.Unlock()

	w.once.Do(func() {
		for _, c := range closers {
			if err := c(); err != nil {
				fmt.Println(err)
			}
		}
	})
}
