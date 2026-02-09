package main

import "fmt"

func decorator[T any](ch <-chan T, action func(T) T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)
		for v := range ch {
			out <- action(v)
		}
	}()

	return out
}

func main() {

	ch1 := make(chan int)

	go func() {
		defer close(ch1)
		for i := range 10 {
			ch1 <- i
		}
	}()

	decorate := func(v int) int {
		return v * v
	}

	ch2 := decorator(ch1, decorate)

	for v := range ch2 {
		fmt.Println(v)
	}

}
