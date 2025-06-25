package main

import (
	"fmt"
	"math/rand"
	"time"
)

func repeatFunc[T any, K any](done <-chan K, fn func() T) <-chan T {
	stream := make(chan T)

	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				return

			case stream <- fn():
			}
		}
	}()

	return stream
}

func take[T any, K any](done <-chan K, stream <-chan T, n int) <-chan T {
	taken := make(chan T)

	go func() {
		defer close(taken)
		for range n {
			select {
			case <-done:
				return
			case receivedVal := <-stream:
				taken <- receivedVal
			}
		}
	}()

	return taken
}

// this is expected to be slow on purpose
// for us to use another pattern to solve slow stages.
func primeFinder(done <-chan bool, randIntStream <-chan int) <-chan int {
	isPrime := func(randomInt int) bool {
		for i := randomInt - 1; i > 1; i-- {
			if randomInt%i == 0 {
				return false
			}
		}
		return true
	}

	primes := make(chan int)
	go func() {
		defer close(primes)
		for {
			select {
			case <-done:
				return
			case randomInt := <-randIntStream:
				if isPrime(randomInt) {
					primes <- randomInt
				}
			}
		}
	}()

	return primes
}

func main() {
	done := make(chan bool)
	defer close(done)

	randomNumberFn := func() int {
		return rand.Intn(100_000_000)
	}

	start := time.Now()

	randIntStream := repeatFunc(done, randomNumberFn)
	primeStream := primeFinder(done, randIntStream)

	i := 0
	for val := range take(done, primeStream, 20) {
		fmt.Printf("%d - %d\n", i, val)
		i++
	}

	fmt.Printf("Total Execution time: (in seconds): %f\n", time.Since(start).Seconds())
}
