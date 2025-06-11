package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	startAtCep := 0o7300000
	endAtCep := 0o7400000
	workers := 1000

	var wg sync.WaitGroup
	jobs := make(chan Cep, endAtCep-startAtCep)
	results := make(chan CepResult, endAtCep-startAtCep)

	for range workers {
		wg.Add(1)
		go processCep(&wg, jobs, results)
	}

	for i := startAtCep; i < endAtCep; i++ {
		jobs <- Cep(i)
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	addressFromBrazil := make([]CepResult, 0, endAtCep-startAtCep)
	for res := range results {
		addressFromBrazil = append(addressFromBrazil, res)
	}

	fmt.Println("Amount of CEPs in Brazil Processed: ", len(addressFromBrazil))

	amountOfValidCeps := 0
	amountOfInvalidCeps := 0

	var sampleValidResult *CepResult
	var sampleInvalidResult *CepResult

	for _, address := range addressFromBrazil {
		if address.Error == nil {
			amountOfValidCeps++

			if sampleValidResult == nil {
				sampleValidResult = &address
			}

		} else {
			amountOfInvalidCeps++

			if sampleInvalidResult == nil {
				sampleInvalidResult = &address
			}
		}
	}

	fmt.Println("Amount of VALID CEPs Processed: ", amountOfValidCeps)
	fmt.Println("Amount of INVALID CEPs Processed: ", amountOfInvalidCeps)

	fmt.Printf("Sample VALID Result: %s\n", sampleValidResult.Address.Cep)
	fmt.Printf("Sample INVALID Result: %s\n", sampleInvalidResult.Address.Cep)
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
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := fmt.Sprintf("https://viacep.com.br/ws/%08d/json", cep)

	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return CepResult{
			Error: fmt.Errorf("received an error in the GET request: %v", err),
		}
	}

	defer resp.Body.Close()

	var address Address
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return CepResult{
			Error: fmt.Errorf("received an error while parsing the body: %v", err),
		}
	}

	if address.Error == "true" {
		return CepResult{
			Error: fmt.Errorf("The API returned an error in the body: field 'error' = %v", address.Error),
		}
	}

	return CepResult{
		Address: address,
	}
}
