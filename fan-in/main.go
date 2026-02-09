package main

import (
	"fmt"
	"sync"
)

func mergeChannels[T any](chans ...<-chan T) <-chan T {
	out := make(chan T)
	wg := &sync.WaitGroup{}

	for _, ch := range chans {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range ch {
				out <- v
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go func() {
		defer func() {
			close(ch1)
			close(ch2)
			close(ch3)
		}()

		for i := range 100 {
			ch1 <- i
			ch2 <- i - 1000
			ch3 <- i + 1000
		}
	}()

	mergeCh := mergeChannels(ch1, ch2, ch3)

	for v := range mergeCh {
		fmt.Println(v)
	}
}
