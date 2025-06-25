package main

import (
	"fmt"
	"sync"
)

func process(wg *sync.WaitGroup, i int, res chan int) {
	defer wg.Done()
	res <- i * 2
}

func main() {
	info := []int{}
	res := make(chan int, 10)

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go process(&wg, i, res)
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	for val := range res {
		info = append(info, val)
	}
	// Alternative:
	// outer:
	// 	for {
	// 		select {
	// 		case val, ok := <-res:
	// 			if !ok {
	// 				break outer
	// 			}
	// 			info = append(info, val)
	// 		}
	// 	}

	fmt.Println("expected 100 got", len(info))
}
