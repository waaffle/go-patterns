package main

import (
	"fmt"
	"sync"
)

func splitCh[T any](in <-chan T, n int) []<-chan T {
	out := make([]chan T, n)
	for i := range out {
		out[i] = make(chan T)
	}

	go func() {
		for v := range in {
			for i := range out {
				out[i] <- v
			}
		}
		for i := range out {
			close(out[i])
		}
	}()

	// cannot cast []chan T to []<-chan T
	resCh := make([]<-chan T, n) 
	for i := range resCh {
		resCh[i] = out[i]
	}

	return resCh
}

func main() {
	in := make(chan int)
	wg := &sync.WaitGroup{}

	go func() {
		defer close(in)
		for i := range 10 {
			in <- i
		}
	}()

	chans := splitCh(in, 2)
	wg.Add(2)
	go func() {
		defer wg.Done()
		for v := range chans[0] {
			fmt.Println("ch1 ", v)
		}
	}()
	go func() {
		defer wg.Done()
		for v := range chans[1] {
			fmt.Println("ch2 ", v)
		}
	}()

	wg.Wait()

}
