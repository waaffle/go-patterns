package main

import (
	"fmt"
	"strconv"
	"sync"
)

func parse(in <-chan string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		for v := range in {
			out <- fmt.Sprintf("parsed: %v", v)
		}
	}()

	return out
}

func send(in <-chan string, parralelFactor int) <-chan string {
	out := make(chan string)
	wg := &sync.WaitGroup{}

	for i := range parralelFactor {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range in {
				out <- fmt.Sprintf("Отправлено горутиной %v значение %v", i+1, v)
			}
		}()
	}

	go func(){
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	in := make(chan string)

	go func() {
		defer close(in)
		for i := range 10 {
			in <- strconv.Itoa(i)
		}
	}()

	for v := range send(parse(in), 2) {
		fmt.Println(v)
	}
}
