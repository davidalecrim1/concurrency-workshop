# Concurrency Workshop

This is a practical guide of how to apply some concorrency patterns in Golang.

## Worker Pool Pattern:

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

### Hands on folder `worker_pool` of this pattern:

This used an external API (Via CEP from Brazil) to apply a worker pool.

Using workers=1:

```bash
Amount of Workers:  1
Start time:  2025-06-15T18:09:46-03:00
Amount of CEPs in Brazil Processed:  999999
End time:  2025-06-19T14:54:17-03:00
Diff time (in seconds):  199456.730876541
Amount of VALID CEPs Processed:  12824
Amount of INVALID CEPs Processed:  987175
Sample of VALID CEPs:
CEP: 07010-000, City: Guarulhos, UF: SP
CEP: 07010-001, City: Guarulhos, UF: SP
CEP: 07010-003, City: Guarulhos, UF: SP
CEP: 07010-010, City: Guarulhos, UF: SP
CEP: 07010-015, City: Guarulhos, UF: SP
Sample of INVALID CEPs:
CEP: 07000000
CEP: 07000001
CEP: 07000002
CEP: 07000003
CEP: 07000004
```

Using workers=100:

```bash
Amount of Workers:  100
Start time:  2025-06-16T09:30:28-03:00
Amount of CEPs in Brazil Processed:  999999
End time:  2025-06-16T18:55:22-03:00
Diff time (in seconds):  31771.85346225
Amount of VALID CEPs Processed:  15077
Amount of INVALID CEPs Processed:  984922
Sample of VALID CEPs:
CEP: 07010-003, City: Guarulhos, UF: SP
CEP: 07010-000, City: Guarulhos, UF: SP
CEP: 07010-001, City: Guarulhos, UF: SP
CEP: 07010-010, City: Guarulhos, UF: SP
CEP: 07010-015, City: Guarulhos, UF: SP
Sample of INVALID CEPs:
CEP: 07000024
CEP: 07000029
CEP: 07000037
CEP: 07000033
CEP: 07000077
```

Using workers=1000

```bash
Amount of Workers:  1000
Start time:  2025-06-18T17:51:05-03:00
Amount of CEPs in Brazil Processed:  999999
End time:  2025-06-18T19:49:39-03:00
Diff time (in seconds):  7110.883386541
Amount of VALID CEPs Processed:  7132
Amount of INVALID CEPs Processed:  992867
Sample of VALID CEPs:
CEP: 07010-000, City: Guarulhos, UF: SP
CEP: 07010-001, City: Guarulhos, UF: SP
CEP: 07010-003, City: Guarulhos, UF: SP
CEP: 07010-010, City: Guarulhos, UF: SP
CEP: 07010-015, City: Guarulhos, UF: SP
Sample of INVALID CEPs:
CEP: 07000600
CEP: 07000729
CEP: 07000840
CEP: 07000717
CEP: 07000777
```

Using workers=10000

```bash
Amount of Workers:  10000
Start time:  2025-06-18T14:50:08-03:00
Amount of CEPs in Brazil Processed:  999999
End time:  2025-06-18T15:04:31-03:00
Diff time (in seconds):  863.83120925
Amount of VALID CEPs Processed:  36
Amount of INVALID CEPs Processed:  999963
Sample of VALID CEPs:
CEP: 07064-021, City: Guarulhos, UF: SP
CEP: 07064-020, City: Guarulhos, UF: SP
CEP: 07064-022, City: Guarulhos, UF: SP
CEP: 07152-842, City: Guarulhos, UF: SP
CEP: 07152-834, City: Guarulhos, UF: SP
Sample of INVALID CEPs:
CEP: 07004160
CEP: 07009353
CEP: 07002995
CEP: 07005288
CEP: 07000370
```

Beyond the 10_000, it overloaded my network and increase the amount of invalid CEPs cause by network connection.

## Other Patterns

- Pipeline
- Fan Out && Fan In
- Generator

We will create an example using all of this pattern in the folder `pipeline`.

The pipeline pattern is a concurrency pattern that connects stages via channels, where each stage transforms data and sends it to the next stage. It’s great for decoupling, parallelism, and stream processing.

Sample code:

```go
package main

import (
    "fmt"
)

func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func multiplyBy2(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * 2
        }
        close(out)
    }()
    return out
}

func add1(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n + 1
        }
        close(out)
    }()
    return out
}

func main() {
    in := gen(1, 2, 3)
    out := add1(multiplyBy2(in))

    for result := range out {
        fmt.Println(result) // 3, 5, 7
    }
}
```

### Hands on folder `pipeline` of the `generator`, `pipeline` and `fan out` pattern:

Without Fan In / Fan Out pattern (folder `001`):

```bash
davidalecrim@Davids-Macbook-Pro-M4-Pro 001 % go run main.go
0 - 61593341
1 - 96418121
2 - 60666797
3 - 36491321
4 - 28664239
5 - 7691863
6 - 6207367
7 - 84081713
8 - 40813369
9 - 37995019
10 - 14608591
11 - 35037193
12 - 53967547
13 - 67103549
14 - 38412467
15 - 2698321
16 - 71137097
17 - 11231081
18 - 93647527
19 - 37801867
Total Execution time: (in seconds): 4.221439
```

With Fan In and Fan Out Pattern (folder `002`):

```bash
davidalecrim@Davids-Macbook-Pro-M4-Pro 002 % go run main.go
Amount of CPUs:  12
0 - 14005063
1 - 16417229
2 - 13178293
3 - 36242369
4 - 57813551
5 - 97841369
6 - 65167831
7 - 44850517
8 - 22551409
9 - 66343579
10 - 68644973
11 - 48501823
12 - 96359057
13 - 15607373
14 - 23021861
15 - 94668961
16 - 13219571
17 - 41215879
18 - 51082751
19 - 20677051
Total Execution time: (in seconds): 0.628412
```

A buffer didn't changed much the results.
