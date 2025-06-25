package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
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

func fanIn[T any, K any](done <-chan K, channels []<-chan T) <-chan T {
	var wg sync.WaitGroup
	fannedInStream := make(chan T)

	transfer := func(c <-chan T) {
		defer wg.Done()

		for val := range c {
			select {
			case <-done:
				return
			case fannedInStream <- val:
				continue
			}
		}
	}

	for _, c := range channels {
		wg.Add(1)
		go transfer(c)
	}

	go func() {
		wg.Wait()
		close(fannedInStream)
	}()

	return fannedInStream
}

func main() {
	done := make(chan bool)
	defer close(done)

	randomNumberFn := func() int {
		return rand.Intn(100_000_000)
	}

	start := time.Now()

	randIntStream := repeatFunc(done, randomNumberFn)

	// fan out
	numCpus := runtime.NumCPU()
	fmt.Println("Amount of CPUs: ", numCpus)

	primeFinderChannels := make([]<-chan int, numCpus)
	for i := range numCpus {
		primeFinderChannels[i] = primeFinder(done, randIntStream)
	}

	// fan in
	fannedInStream := fanIn(done, primeFinderChannels)

	i := 0
	for val := range take(done, fannedInStream, 20) {
		fmt.Printf("%d - %d\n", i, val)
		i++
	}

	fmt.Printf("Total Execution time: (in seconds): %f\n", time.Since(start).Seconds())
}
