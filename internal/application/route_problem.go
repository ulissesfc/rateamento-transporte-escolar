package application

import (
	"fmt"

	osrm "github.com/ulissesfc/rateamento-transporte-escolar.git/internal/adapter/http/orsm"
	"github.com/ulissesfc/rateamento-transporte-escolar.git/internal/domain"
)

type RouteService struct {
	OsrmClient *osrm.Client
}

func (rs *RouteService) GenerateRouteForCluster(problem domain.RouteProblem) ([]int, error) {

	// 1. Criar uma lista de todos os locais relevantes para este cluster
	locations := []domain.Location{problem.GaragemLocation, problem.UFSCLocation}
	studentsInCluster := make(map[int]domain.Student)
	for _, studentID := range problem.Cluster.StudentIDs {
		for _, s := range problem.AllStudents {
			if s.Id == studentID {
				locations = append(locations, domain.Location{ID: s.Id, Lat: s.Latitude, Lon: s.Longitude})
				studentsInCluster[s.Id] = s
				break
			}
		}
	}

	// 2. Criar mapas de tradução ID <-> Índice da Matriz
	idParaIndice := make(map[int]int)
	for i, loc := range locations {
		idParaIndice[loc.ID] = i
	}

	// 3. Chamar OSRM para obter a matriz de N+2 x N+2
	var coords []osrm.Coord
	for _, coord := range locations {
		coords = append(coords, osrm.Coord{Latitude: coord.Lat, Longitude: coord.Lon})
	}
	matrix, err := rs.OsrmClient.RequestDurationMatrix(coords)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter matriz de duração para o cluster: %w", err)
	}

	// --- PASSO B: Definir os pontos de início e fim da rota (fiel ao artigo) ---
	idxGaragem := idParaIndice[problem.GaragemLocation.ID]
	idxUfsc := idParaIndice[problem.UFSCLocation.ID]

	primeiroAlunoID := -1
	ultimoAlunoID := -1
	minDuracaoGaragem := 1e9
	minDuracaoUfsc := 1e9

	for studentID := range studentsInCluster {
		idxAluno := idParaIndice[studentID]
		if matrix[idxGaragem][idxAluno] < minDuracaoGaragem {
			minDuracaoGaragem = matrix[idxGaragem][idxAluno]
			primeiroAlunoID = studentID
		}
		if matrix[idxAluno][idxUfsc] < minDuracaoUfsc {
			minDuracaoUfsc = matrix[idxAluno][idxUfsc]
			ultimoAlunoID = studentID
		}
	}

	if primeiroAlunoID == -1 { // Se o cluster estiver vazio, retorne rota vazia
		if len(studentsInCluster) > 0 {
			return nil, fmt.Errorf("não foi possível definir o ponto de início para o cluster")
		}
		return []int{}, nil
	}
	if ultimoAlunoID == -1 {
		ultimoAlunoID = primeiroAlunoID // Garante que temos um ponto final
	}

	// --- PASSO C: Inicializar a rota e a lista de não roteirizados ---
	rotaOrdenada := []int{primeiroAlunoID}
	if primeiroAlunoID != ultimoAlunoID {
		rotaOrdenada = append(rotaOrdenada, ultimoAlunoID)
	}

	naoRoteirizados := make(map[int]bool)
	for studentID := range studentsInCluster {
		if studentID != primeiroAlunoID && studentID != ultimoAlunoID {
			naoRoteirizados[studentID] = true
		}
	}

	// --- CORREÇÃO PARA ROTAS COM APENAS UM PONTO INICIAL ---
	// Se a rota inicial tem apenas um ponto (início == fim),
	// precisamos adicionar o segundo ponto para criar a primeira "aresta".
	if len(rotaOrdenada) == 1 && len(naoRoteirizados) > 0 {
		pontoInicialID := rotaOrdenada[0]
		idxPontoInicial := idParaIndice[pontoInicialID]

		melhorProximoID := -1
		menorDuracaoParaProximo := 1e9

		// Encontra o aluno não roteirizado mais próximo do ponto inicial
		for id := range naoRoteirizados {
			idxCandidato := idParaIndice[id]
			duracao := matrix[idxPontoInicial][idxCandidato]
			if duracao < menorDuracaoParaProximo {
				menorDuracaoParaProximo = duracao
				melhorProximoID = id
			}
		}

		if melhorProximoID != -1 {
			rotaOrdenada = append(rotaOrdenada, melhorProximoID)
			delete(naoRoteirizados, melhorProximoID)
		}
	}

	// --- PASSO D: Loop de Inserção mais Econômica ---
	for len(naoRoteirizados) > 0 {
		melhorCustoInsercao := 1e9
		melhorAlunoParaInserir := -1
		melhorPosicao := -1

		for alunoK_ID := range naoRoteirizados {
			idxK := idParaIndice[alunoK_ID]
			for i := 0; i < len(rotaOrdenada)-1; i++ {
				alunoI_ID := rotaOrdenada[i]
				alunoJ_ID := rotaOrdenada[i+1]
				idxI := idParaIndice[alunoI_ID]
				idxJ := idParaIndice[alunoJ_ID]

				custoInsercao := matrix[idxI][idxK] + matrix[idxK][idxJ] - matrix[idxI][idxJ]
				if custoInsercao < melhorCustoInsercao {
					melhorCustoInsercao = custoInsercao
					melhorAlunoParaInserir = alunoK_ID
					melhorPosicao = i + 1
				}
			}
		}

		if melhorAlunoParaInserir != -1 {
			rotaOrdenada = append(rotaOrdenada[:melhorPosicao], append([]int{melhorAlunoParaInserir}, rotaOrdenada[melhorPosicao:]...)...)
			delete(naoRoteirizados, melhorAlunoParaInserir)
		} else {
			break
		}
	}

	return rotaOrdenada, nil
}
