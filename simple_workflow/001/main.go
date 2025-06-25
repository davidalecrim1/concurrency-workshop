package main

import (
	"fmt"
	"sync"
)

// This code suffers from a race condition because we access the data in the same time.
// We Should solve this with a channel or a mutex.
func process(wg *sync.WaitGroup, i int, info *[]int) {
	defer wg.Done()
	*info = append(*info, i*2)
}

func main() {
	info := []int{}

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go process(&wg, i, &info)
	}

	wg.Wait()
	fmt.Println(info)
	fmt.Println("expected 100 got", len(info))
}
