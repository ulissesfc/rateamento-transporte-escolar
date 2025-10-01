package application

import "github.com/ulissesfc/rateamento-transporte-escolar.git/internal/domain"

type SavingService struct {
}

func (s *SavingService) CalculatesSavings(problem *domain.SavingProblem) []domain.Saving {
	var savings []domain.Saving
	alunoIndexMap := make(map[int]int)
	for i, a := range problem.Students {
		alunoIndexMap[a.Id] = i
	}

	//fórmula = s(i,j,k) = c(i,k_i) + c(j,k_j) - c(i,j)
	//Onde k_i é a semente mais próxima de i, e k_j é a mais próxima de j.
	//Esta é uma simplificação da fórmula original para o estado inicial.

	// Itera sobre todos os pares de alunos (que não são sementes)
	for id1, cluster1 := range problem.InitialSolution {
		for id2, cluster2 := range problem.InitialSolution {
			if id1 >= id2 {
				continue
			}

			// Para cada semente, calcula a economia de juntar id1 e id2 nela
			for _, semente := range problem.StudentsSeeds {
				idx1 := alunoIndexMap[id1]
				idx2 := alunoIndexMap[id2]
				idxSemente1 := alunoIndexMap[cluster1.SeedID]
				idxSemente2 := alunoIndexMap[cluster2.SeedID]

				// Custo original: (Semente1 -> Aluno1 -> Semente1) + (Semente2 -> Aluno2 -> Semente2)
				// Como estamos usando OSRM, a viagem de volta pode ser diferente,
				// mas para simplificar, usaremos 2 * duração.
				custoOriginal := 2*(*problem.DurationsMatrix)[idxSemente1][idx1] + 2*(*problem.DurationsMatrix)[idxSemente2][idx2]

				// Novo custo: SementeK -> Aluno1 -> Aluno2 -> SementeK
				idxSementeK := alunoIndexMap[semente.Id]
				novoCusto := (*problem.DurationsMatrix)[idxSementeK][idx1] + (*problem.DurationsMatrix)[idx1][idx2] + (*problem.DurationsMatrix)[idx2][idxSementeK]

				economia := custoOriginal - novoCusto

				if economia > 0 {
					savings = append(savings, domain.Saving{
						Value:    economia,
						StudentI: id1,
						StudentJ: id2,
						SeedK:    semente.Id,
					})
				}
			}
		}
	}
	return savings
}
