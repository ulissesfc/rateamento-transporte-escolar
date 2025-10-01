package osrm

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DurationMatrix struct {
	Code      string      `json:"code"`
	Durations [][]float64 `json:"durations"`
}

type DistanceMatrix struct {
	Code      string      `json:"code"`
	Distances [][]float64 `json:"distances"`
}

func (c *Client) RequestDurationMatrix(Coords []Coord) ([][]float64, error) {

	pointsBuilder := pointsBuilder(Coords)

	requestURL := fmt.Sprintf("%s/table/v1/driving/%s?annotations=duration", c.baseURL, pointsBuilder)

	fmt.Println("Buscando matriz de durações do OSRM...")
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("falha na requisição ao OSRM: %w", err)
	}
	defer resp.Body.Close()

	var osrmResponse DurationMatrix
	err = json.NewDecoder(resp.Body).Decode(&osrmResponse)
	if err != nil {
		return nil, fmt.Errorf("falha ao decodificar a resposta do OSRM: %w", err)
	}

	if osrmResponse.Code != "Ok" {
		return nil, fmt.Errorf("OSRM retornou um erro: %s", osrmResponse.Code)
	}

	fmt.Println("Matriz de durações obtida com sucesso.")
	return osrmResponse.Durations, nil
}

func (c *Client) RequestDistanceMatrix(pointsBuilder string) ([][]float64, error) {
	requestURL := fmt.Sprintf("%s/table/v1/driving/%s?annotations=duration", c.baseURL, pointsBuilder)

	fmt.Println("Buscando matriz de durações do OSRM...")
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("falha na requisição ao OSRM: %w", err)
	}
	defer resp.Body.Close()

	var osrmResponse DurationMatrix
	err = json.NewDecoder(resp.Body).Decode(&osrmResponse)
	if err != nil {
		return nil, fmt.Errorf("falha ao decodificar a resposta do OSRM: %w", err)
	}

	if osrmResponse.Code != "Ok" {
		return nil, fmt.Errorf("OSRM retornou um erro: %s", osrmResponse.Code)
	}

	fmt.Println("Matriz de durações obtida com sucesso.")
	return osrmResponse.Durations, nil
}
