package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
	},
	Timeout: 10 * time.Second,
}

func main() {
	startAtCep := 7000000
	endAtCep := 7999999
	workers := 1000

	var wg sync.WaitGroup
	jobs := make(chan Cep, workers)
	results := make(chan CepResult, workers)

	startTime := time.Now()
	fmt.Println("Amount of Workers: ", workers)
	fmt.Println("Start time: ", startTime.Format(time.RFC3339))

	for range workers {
		wg.Add(1)
		go processCep(&wg, jobs, results)
	}

	go func() {
		for i := startAtCep; i < endAtCep; i++ {
			jobs <- Cep(i)
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	ceps := make([]CepResult, 0, endAtCep-startAtCep)
	for res := range results {
		ceps = append(ceps, res)
	}

	fmt.Println("Amount of CEPs in Brazil Processed: ", len(ceps))

	endTime := time.Now()
	fmt.Println("End time: ", endTime.Format(time.RFC3339))
	fmt.Println("Diff time (in seconds): ", endTime.Sub(startTime).Seconds())

	amountOfValidCeps := 0
	amountOfInvalidCeps := 0

	for _, address := range ceps {
		if address.Error == nil {
			amountOfValidCeps++
		} else {
			amountOfInvalidCeps++
		}
	}

	fmt.Println("Amount of VALID CEPs Processed: ", amountOfValidCeps)
	fmt.Println("Amount of INVALID CEPs Processed: ", amountOfInvalidCeps)

	count := 0
	fmt.Println("Sample of VALID CEPs: ")
	for _, cep := range ceps {
		if cep.Error == nil {
			fmt.Printf("CEP: %s, City: %s, UF: %s\n", cep.Address.Cep, cep.Address.Localidade, cep.Address.UF)
			count++
		}

		if count == 5 {
			break
		}
	}

	count = 0
	fmt.Println("Sample of INVALID CEPs: ")
	for _, cep := range ceps {
		if cep.Error != nil {
			fmt.Printf("CEP: %s\n", cep.Address.Cep)
			count++
		}

		if count == 5 {
			break
		}
	}
}

func processCep(wg *sync.WaitGroup, jobs chan Cep, results chan CepResult) {
	defer wg.Done()

	for cep := range jobs {
		res := requestCep(cep)
		results <- res
	}
}

type Cep int

type CepResult struct {
	Address Address
	Error   error
}

type Address struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia"`
	DDD         string `json:"ddd"`
	SIAFI       string `json:"siafi"`
	Error       string `json:"erro"`
}

func requestCep(cep Cep) CepResult {
	url := fmt.Sprintf("https://viacep.com.br/ws/%08d/json", cep)
	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return CepResult{
			Address: Address{
				Cep: fmt.Sprintf("%08d", cep),
			},
			Error: fmt.Errorf("received an error in the GET request: %v", err),
		}
	}

	defer resp.Body.Close()

	var address Address
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return CepResult{
			Address: Address{
				Cep: fmt.Sprintf("%08d", cep),
			},
			Error: fmt.Errorf("received an error while parsing the body: %v", err),
		}
	}

	if address.Error == "true" {
		return CepResult{
			Address: Address{
				Cep: fmt.Sprintf("%08d", cep),
			},
			Error: fmt.Errorf("The API returned an error in the body: field 'error' = %v", address.Error),
		}
	}

	return CepResult{
		Address: address,
	}
}
