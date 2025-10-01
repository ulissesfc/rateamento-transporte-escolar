package application

import (
	"sort"

	"github.com/ulissesfc/rateamento-transporte-escolar.git/internal/domain"
)

type ClusterService struct {
}

/*
func (c *ClusterService) CreateInitialSolution(problem *domain.ClusterProblem1) map[int]*domain.Cluster {

	solucao := make(map[int]*domain.Cluster)
	alunoIndexMap := make(map[int]int)
	for i, a := range problem.Students {
		alunoIndexMap[a.Id] = i
	}

	for _, aluno := range problem.Students {
		// Ignora os próprios alunos-semente por enquanto
		isSeed := false
		for _, s := range problem.StudentsSeeds {
			if s.Id == aluno.Id {
				isSeed = true
				break
			}
		}
		if isSeed {
			continue
		}

		// Encontra a semente mais próxima para este aluno
		sementeMaisProximaID := -1
		menorDuracao := 1e9 // infinito

		idxAluno := alunoIndexMap[aluno.Id]

		for _, semente := range problem.StudentsSeeds {
			idxSemente := alunoIndexMap[semente.Id]
			duracao := (*problem.DurationsMatrix)[idxAluno][idxSemente]
			if duracao < menorDuracao {
				menorDuracao = duracao
				sementeMaisProximaID = semente.Id
			}
		}

		// Cria um cluster/rota individual para este aluno
		solucao[aluno.Id] = &domain.Cluster{
			SeedID:      sementeMaisProximaID,
			StudentIDs:  []int{aluno.Id}, // Rota contém apenas ele mesmo
			TotalDemand: 1,               // Assumindo que cada aluno tem demanda 1
		}
	}

	return solucao
}


func (c *ClusterService) GenerateClusters(problem *domain.ClusterProblem2) []*domain.Cluster {
	// Ordena as economias da maior para a menor
	sort.Slice(problem.Savings, func(i, j int) bool {
		return problem.Savings[i].Value > problem.Savings[j].Value
	})

	rotasAtuais := problem.InitialSolution

	// Para rastrear a qual cluster cada aluno pertence
	alunoParaCluster := make(map[int]*domain.Cluster)
	for id, cluster := range rotasAtuais {
		alunoParaCluster[id] = cluster
	}

	for _, s := range problem.Savings {
		clusterI := alunoParaCluster[s.StudentI]
		clusterJ := alunoParaCluster[s.StudentJ]

		// --- Validações ---
		// 1. Os alunos já estão no mesmo cluster? Se sim, pule.
		if clusterI == clusterJ {
			continue
		}

		// 2. Os alunos são as extremidades de suas rotas?
		// (No nosso modelo simplificado, sempre são, até a primeira fusão)
		alunoI_ehExtremidade := (clusterI.StudentIDs[0] == s.StudentI || clusterI.StudentIDs[len(clusterI.StudentIDs)-1] == s.StudentI)
		alunoJ_ehExtremidade := (clusterJ.StudentIDs[0] == s.StudentJ || clusterJ.StudentIDs[len(clusterJ.StudentIDs)-1] == s.StudentJ)
		if !alunoI_ehExtremidade || !alunoJ_ehExtremidade {
			continue
		}

		// 3. A capacidade do maior ônibus disponível suporta a fusão?
		demandaTotal := clusterI.TotalDemand + clusterJ.TotalDemand
		capacidadeDisponivel := false
		// Aqui, precisaríamos de uma lógica mais complexa para alocar ônibus,
		// mas para simplificar, vamos verificar contra a capacidade máxima da frota.
		maxCapacity := 0
		for _, b := range problem.VehiclesAvailable {
			if b.Capacity > maxCapacity {
				maxCapacity = b.Capacity
			}
		}
		if demandaTotal <= maxCapacity {
			capacidadeDisponivel = true
		}

		if !capacidadeDisponivel {
			continue
		}

		// --- FUSÃO ---
		// A lógica aqui seria complexa: qual ponta ligar com qual?
		// Simplificação: apenas juntamos as listas.
		novaListaAlunos := append(clusterI.StudentIDs, clusterJ.StudentIDs...)
		clusterI.StudentIDs = novaListaAlunos
		clusterI.TotalDemand = demandaTotal

		// Atualiza o mapeamento: agora todos os alunos de J apontam para o cluster I
		for _, alunoID := range clusterJ.StudentIDs {
			alunoParaCluster[alunoID] = clusterI
		}

		// Remove o cluster J, que foi absorvido
		// (em uma implementação real, marcaríamos como inativo)
	}

	// Coleta os clusters finais que restaram
	clustersFinaisMap := make(map[*domain.Cluster]bool)
	for _, cluster := range alunoParaCluster {
		clustersFinaisMap[cluster] = true
	}

	var clustersFinais []*domain.Cluster
	for c := range clustersFinaisMap {
		clustersFinais = append(clustersFinais, c)
	}

	return clustersFinais
}
*/

func (c *ClusterService) CreateInitialSolution(problem *domain.ClusterProblem1) map[int]*domain.Cluster {

	solucao := make(map[int]*domain.Cluster)
	alunoIndexMap := make(map[int]int)
	for i, a := range problem.Students {
		alunoIndexMap[a.Id] = i
	}

	// Use um mapa para a verificação de sementes para melhor performance
	seedIDs := make(map[int]bool)
	for _, s := range problem.StudentsSeeds {
		seedIDs[s.Id] = true
	}

	for _, aluno := range problem.Students {
		// Se o aluno é uma semente, pule
		if seedIDs[aluno.Id] {
			continue
		}

		// Encontra a semente mais próxima para este aluno
		sementeMaisProximaID := -1
		menorDuracao := 1e9 // infinito
		idxAluno := alunoIndexMap[aluno.Id]

		for _, semente := range problem.StudentsSeeds {
			idxSemente := alunoIndexMap[semente.Id]
			duracao := (*problem.DurationsMatrix)[idxAluno][idxSemente]
			if duracao < menorDuracao {
				menorDuracao = duracao
				sementeMaisProximaID = semente.Id
			}
		}

		// Cria um cluster/rota individual para este aluno
		solucao[aluno.Id] = &domain.Cluster{
			SeedID:      sementeMaisProximaID,
			StudentIDs:  []int{aluno.Id}, // Rota contém apenas ele mesmo
			TotalDemand: 1,               // Assumindo que cada aluno tem demanda 1

			// --- CORREÇÃO: INICIALIZAR OS NOVOS CAMPOS ---
			Endpoint1ID: aluno.Id, // No início, o aluno é ambas as extremidades da rota.
			Endpoint2ID: aluno.Id, // Essencial para a lógica de fusão funcionar.
			VehicleID:   0,        // Nenhum veículo alocado ainda (use 0 ou -1 como indicador).
		}
	}

	return solucao
}

func (c *ClusterService) GenerateClusters(problem *domain.ClusterProblem2) []*domain.Cluster {
	sort.Slice(problem.Savings, func(i, j int) bool {
		return problem.Savings[i].Value > problem.Savings[j].Value
	})

	// Encontra a capacidade máxima de um veículo na frota.
	// O algoritmo usará isso como o limite máximo que um cluster pode atingir.
	maxCapacity := 0
	for _, v := range problem.VehiclesAvailable {
		if v.Capacity > maxCapacity {
			maxCapacity = v.Capacity
		}
	}

	// Mapeia o ID de um aluno para o cluster ao qual ele pertence.
	alunoParaCluster := make(map[int]*domain.Cluster)
	for _, cluster := range problem.InitialSolution {
		alunoParaCluster[cluster.StudentIDs[0]] = cluster
	}

	// --- Loop Principal de Fusão ---
	for _, s := range problem.Savings {
		clusterI := alunoParaCluster[s.StudentI]
		clusterJ := alunoParaCluster[s.StudentJ]

		if clusterI == nil || clusterJ == nil || clusterI == clusterJ {
			continue
		}

		alunoI_ehExtremidade := clusterI.Endpoint1ID == s.StudentI || clusterI.Endpoint2ID == s.StudentI
		alunoJ_ehExtremidade := clusterJ.Endpoint1ID == s.StudentJ || clusterJ.Endpoint2ID == s.StudentJ
		if !alunoI_ehExtremidade || !alunoJ_ehExtremidade {
			continue
		}

		// --- LÓGICA DE CAPACIDADE SIMPLIFICADA E CORRETA PARA ESTA FASE ---
		novaDemanda := clusterI.TotalDemand + clusterJ.TotalDemand
		if novaDemanda > maxCapacity {
			continue // A fusão criaria um cluster maior que o nosso maior ônibus.
		}

		// --- FUSÃO ---
		if clusterI.Endpoint2ID != s.StudentI {
			clusterI.Reverse()
		}
		if clusterJ.Endpoint1ID != s.StudentJ {
			clusterJ.Reverse()
		}

		clusterI.StudentIDs = append(clusterI.StudentIDs, clusterJ.StudentIDs...)
		clusterI.TotalDemand = novaDemanda
		clusterI.Endpoint2ID = clusterJ.Endpoint2ID

		// Atualiza o mapeamento para todos os alunos do cluster absorvido.
		for _, alunoID := range clusterJ.StudentIDs {
			alunoParaCluster[alunoID] = clusterI
		}
	}

	// --- Coleta dos Resultados Finais ---
	clustersFinaisMap := make(map[*domain.Cluster]bool)
	for _, cluster := range alunoParaCluster {
		if cluster != nil {
			clustersFinaisMap[cluster] = true
		}
	}

	var clustersFinais []*domain.Cluster
	for c := range clustersFinaisMap {
		clustersFinais = append(clustersFinais, c)
	}

	return clustersFinais
}
