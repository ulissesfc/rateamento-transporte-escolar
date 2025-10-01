package application

import (
	"fmt"

	osrm "github.com/ulissesfc/rateamento-transporte-escolar.git/internal/adapter/http/orsm"
	"github.com/ulissesfc/rateamento-transporte-escolar.git/internal/domain"
)

// RouteService contém a lógica para gerar a rota final.
type RouteService struct {
	OsrmClient *osrm.Client // Supondo que você tenha um client OSRM
}

// GenerateRouteForCluster implementa a "Inserção mais Econômica" para um único cluster.
func (rs *RouteService) GenerateRouteForCluster(problem domain.RouteProblem) ([]int, error) {

	// --- PASSO A: Montar os pontos e obter a matriz de custos específica ---

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

	var coords []osrm.Coord
	for _, coord := range locations {
		coords = append(coords, osrm.Coord{Latitude: coord.Lat, Longitude: coord.Lon})
	}

	// 3. Chamar OSRM para obter a matriz de N+2 x N+2
	matrix, err := rs.OsrmClient.RequestDurationMatrix(coords) // Você precisará criar este método no seu client
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

		// Encontra aluno mais próximo da garagem
		if matrix[idxGaragem][idxAluno] < minDuracaoGaragem {
			minDuracaoGaragem = matrix[idxGaragem][idxAluno]
			primeiroAlunoID = studentID
		}

		// Encontra aluno mais próximo da UFSC
		if matrix[idxAluno][idxUfsc] < minDuracaoUfsc {
			minDuracaoUfsc = matrix[idxAluno][idxUfsc]
			ultimoAlunoID = studentID
		}
	}

	if primeiroAlunoID == -1 || ultimoAlunoID == -1 {
		return nil, fmt.Errorf("não foi possível definir pontos de início/fim para o cluster")
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

	// --- PASSO D: Loop de Inserção mais Econômica ---

	for len(naoRoteirizados) > 0 {
		melhorCustoInsercao := 1e9
		melhorAlunoParaInserir := -1
		melhorPosicao := -1

		// Para cada aluno 'k' ainda não roteirizado...
		for alunoK_ID := range naoRoteirizados {
			idxK := idParaIndice[alunoK_ID]

			// ...teste inseri-lo entre cada par (i, j) da rota atual.
			for i := 0; i < len(rotaOrdenada)-1; i++ {
				alunoI_ID := rotaOrdenada[i]
				alunoJ_ID := rotaOrdenada[i+1]

				idxI := idParaIndice[alunoI_ID]
				idxJ := idParaIndice[alunoJ_ID]

				// Custo de inserção = custo(i->k) + custo(k->j) - custo(i->j)
				custoInsercao := matrix[idxI][idxK] + matrix[idxK][idxJ] - matrix[idxI][idxJ]

				if custoInsercao < melhorCustoInsercao {
					melhorCustoInsercao = custoInsercao
					melhorAlunoParaInserir = alunoK_ID
					melhorPosicao = i + 1 // A inserção ocorre na posição APÓS o aluno 'i'
				}
			}
		}

		// --- PASSO E: Realizar a melhor inserção encontrada ---
		if melhorAlunoParaInserir != -1 {
			// Insere o aluno na melhor posição
			rotaOrdenada = append(rotaOrdenada[:melhorPosicao], append([]int{melhorAlunoParaInserir}, rotaOrdenada[melhorPosicao:]...)...)
			// Remove o aluno da lista de não roteirizados
			delete(naoRoteirizados, melhorAlunoParaInserir)
		} else {
			// Se nenhum aluno puder ser inserido, quebra o loop para evitar loop infinito
			break
		}
	}

	return rotaOrdenada, nil
}
