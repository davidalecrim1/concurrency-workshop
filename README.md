# Concurrency Workshop

This is a practical guide of how to apply the Canonical Go Concurrency Pattern (a.k.a. Worker Pool with Results) in Golang with an actual API from Via CEP in Brazil.

## Canonical Go Concurrency Pattern

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

## Results:

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
