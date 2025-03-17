package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type BrasilApiResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	ApiName      string `json:"-"`
}

type ViaCepResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"estado"`
	City         string `json:"localidade"`
	Neighborhood string `json:"bairro"`
	Street       string `json:"logradouro"`
	ApiName      string `json:"-"`
}

func main() {
	channelBrasilApi := make(chan BrasilApiResponse)
	channelViaCepApi := make(chan ViaCepResponse)

	go GetCepBrasilApi(channelBrasilApi)
	go GetCepViaCepApi(channelViaCepApi)

	select {
	case resViaCEP := <-channelBrasilApi:
		fmt.Printf("BrasilApi: %+v\n", resViaCEP)

	case resApiCEP := <-channelViaCepApi:
		fmt.Printf("ViaCEP: %+v\n", resApiCEP)

	case <-time.After(1 * time.Second):
		fmt.Printf("Timeout")
	}
}

func GetCepBrasilApi(c chan BrasilApiResponse) {
	var brasilApi BrasilApiResponse

	res, err := RequestAPI("https://brasilapi.com.br/api/cep/v1/49040020")
	if err != nil {
		panic("Error on brasilapi request")
	}

	err = json.Unmarshal(res, &brasilApi)
	if err != nil {
		panic("Error on unmarshal brasilapi response" + err.Error())
	}
	brasilApi.ApiName = "BrasilApi"

	c <- brasilApi
}

func GetCepViaCepApi(c chan ViaCepResponse) {
	var viaCEP ViaCepResponse

	res, err := RequestAPI("http://viacep.com.br/ws/49040020/json")
	if err != nil {
		panic("Error on viacep request")
	}

	err = json.Unmarshal(res, &viaCEP)
	if err != nil {
		panic("Error on unmarshal viacep response")
	}
	viaCEP.ApiName = "ViaCEP"

	c <- viaCEP
}

func RequestAPI(url string) ([]byte, error) {
	req, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
