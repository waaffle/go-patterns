package main

import (
	"fmt"
)

func process[T any](closeCh <-chan T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)
		for {
			select {
			case <-closeCh:
				return
			default:
				fmt.Println("default")
				//...processing
			}
		}
	}()

	return out
}

func main() {
	closeCh := make(chan struct{})

	ch := process(closeCh)
	close(closeCh)
	<-ch

	fmt.Println("terminated")
}
