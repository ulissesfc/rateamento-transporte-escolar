package osrm

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Route struct {
	Polyline string  `json:"geometry"`
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
}

func (c *Client) GetRoutes(coords []Coord) ([]Route, error) {

	pointsBuilder := pointsBuilder(coords)
	requestURL := fmt.Sprintf("%s/route/v1/driving/%s?overview=full&geometries=polyline", c.baseURL, pointsBuilder)

	fmt.Println("Buscando rota do OSRM...")
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("falha na requisição ao OSRM: %w", err)
	}
	defer resp.Body.Close()

	var osrmResponse DistanceMatrix
	response := struct {
		StatusCode string  `json:"code"`
		Routes     []Route `json:"routes"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("falha ao decodificar a resposta do OSRM: %w", err)
	}

	if response.StatusCode != "Ok" {
		return nil, fmt.Errorf("OSRM retornou um erro: %s", osrmResponse.Code)
	}

	fmt.Println("Rota obtida com sucesso.")
	return response.Routes, nil
}
