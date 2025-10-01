package application

import (
	"sort"

	"github.com/ulissesfc/rateamento-transporte-escolar.git/internal/domain"
)

type TravelService struct {
}

func (t *TravelService) AllocateBuses(problem domain.TravelProblem) ([]domain.Travel, []*domain.Cluster) {
	var viagensAlocadas []domain.Travel
	var clustersNaoAlocados []*domain.Cluster

	// 1. Ordena os clusters por demanda, do MAIOR para o MENOR.
	sort.Slice(problem.Clusters, func(i, j int) bool {
		return problem.Clusters[i].TotalDemand > problem.Clusters[j].TotalDemand
	})

	// Copia a frota para uma lista de veículos disponíveis que podemos modificar.
	veiculosDisponiveis := make([]domain.Vehicle, len(problem.Fleet))
	copy(veiculosDisponiveis, problem.Fleet)

	// 2. Itera sobre cada cluster (do maior para o menor) e tenta encontrar um ônibus.
	for _, cluster := range problem.Clusters {
		melhorVeiculoIdx := -1
		minCapacityCompativel := 9999

		// Procura pelo ônibus com a MENOR capacidade que ainda serve o cluster (Best Fit).
		for i, veiculo := range veiculosDisponiveis {
			if veiculo.Capacity >= cluster.TotalDemand {
				if veiculo.Capacity < minCapacityCompativel {
					minCapacityCompativel = veiculo.Capacity
					melhorVeiculoIdx = i
				}
			}
		}

		// 3. Verifica se um ônibus foi encontrado.
		if melhorVeiculoIdx != -1 {
			// SUCESSO: Aloca o ônibus ao cluster.
			veiculoAlocado := veiculosDisponiveis[melhorVeiculoIdx]
			viagensAlocadas = append(viagensAlocadas, domain.Travel{
				Vehicle: veiculoAlocado,
				Cluster: cluster,
			})
			// Remove o ônibus da lista de disponíveis (forma eficiente de remover de uma slice).
			veiculosDisponiveis[melhorVeiculoIdx] = veiculosDisponiveis[len(veiculosDisponiveis)-1]
			veiculosDisponiveis = veiculosDisponiveis[:len(veiculosDisponiveis)-1]
		} else {
			// FALHA: Nenhum ônibus disponível comporta este cluster.
			clustersNaoAlocados = append(clustersNaoAlocados, cluster)
		}
	}

	return viagensAlocadas, clustersNaoAlocados
}
