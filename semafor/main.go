package main

import (
	"errors"
	"fmt"
	"go-patterns/semafor/mock"
	"go-patterns/semafor/model"
	"sync"
	"time"
)

type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(gCount int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, gCount),
	}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.ch
}

type ResultUser struct {
	user *model.User
	err  error
}

func DeActivateUser(user *model.User, chClose chan struct{}) ResultUser {
	select {
	case <-time.After(3 * time.Second):
		user.IsActive = false
		return ResultUser{
			user: user,
			err:  nil,
		}
	case <-chClose:
		return ResultUser{
			user: nil,
			err:  errors.New("error"),
		}
	}
}

func main() {
	sem := NewSemaphore(2)
	wg := &sync.WaitGroup{}
	chClose := make(chan struct{})
	outCh := make(chan ResultUser)
	once := &sync.Once{}

	for _, user := range mock.Usrs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem.Acquire()
			defer sem.Release()
			res := DeActivateUser(user, chClose)
			outCh <- res
		}()
	}

	go func() {
		wg.Wait()
		close(outCh)
	}()

	for v := range outCh {
		if v.err != nil {
			once.Do(func() {
				close(chClose)
			})
		} else {
			fmt.Println(v.user.Name)
		}
	}
}
