# Concurrency Workshop

This is a practical guide of how to apply the Canonical Go Concurrency Pattern (a.k.a. Worker Pool with Results) in Golang with an actual API from Via CEP in Brazil.

Sample code for reference:
```go
// 1. Create a jobs channel
jobs := make(chan Job)

// 2. Create a results channel
results := make(chan Result)

// 3. Start N worker goroutines
for i := 0; i < N; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        for job := range jobs {
            result := process(job)
            results <- result
        }
    }()
}

// 4. Start a goroutine to feed jobs
go func() {
    for _, job := range inputJobs {
        jobs <- job
    }
    close(jobs)
}()

// 5. Start a goroutine to close the results channel after all workers are done
go func() {
    wg.Wait()
    close(results)
}()

// 6. Main goroutine consumes results
for result := range results {
    // handle result
}
```