package main

import (
	"fmt"
	"sync"
	"time"
)

func WorkerPool[T any](inputCh <-chan T, nWorkers int) <-chan T {
	outCh := make(chan T)
	wg := &sync.WaitGroup{}

	for i := range nWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				v, ok := <-inputCh
				if !ok {
					return
				}
				time.Sleep(100 * time.Millisecond)
				fmt.Printf("воркер %v, значение %v \n", i, v)
				outCh <- v
			}
		}()
	}

	go func() {
		wg.Wait()
		close(outCh)
	}()

	return outCh
}

func main() {
	in := make(chan int)
	s := []int{}

	go func() {
		defer close(in)
		for i := range 100 {
			in <- i
		}
	}()

	out := WorkerPool(in, 5)

	for v := range out {
		s = append(s, v)
	}
	fmt.Println(s)
}
