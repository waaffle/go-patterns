package main

import (
	"fmt"
)

type Semaphore interface {
	Lock(weight int) error   // блокирующий
	TryLock(weight int) bool // неблокирующий
	Release(weight int)
	Cap() int
	FreeSpace() int
}

type Sem struct {
	ch chan struct{}
}

func NewSemaphore(cap int) Semaphore {
	return &Sem{
		ch: make(chan struct{}, cap),
	}
}

func (s *Sem) Lock(weight int) error {
	for range weight {
		s.ch <- struct{}{}
	}
	return nil
}

func (s *Sem) TryLock(weight int) bool {
	if s.FreeSpace() < weight {
		return false
	}
	for i := range weight {
		select {
		case s.ch <- struct{}{}:
		default:
            s.Release(i)
            return false
		}
	}
	return true
}

func (s *Sem) Release(weight int) {
	for range weight {
		select {
		case <-s.ch:
		default:
		}
	}
}

func (s *Sem) Cap() int {
	return cap(s.ch)
}

func (s *Sem) FreeSpace() int {
	return cap(s.ch) - len(s.ch)
}

func main() {
	sem := NewSemaphore(5)
	sem.Lock(4)
	fmt.Println(sem.FreeSpace())
	sem.Release(3)
	fmt.Println(sem.FreeSpace())
}
