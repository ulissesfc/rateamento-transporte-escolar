package application

import (
	"fmt"

	"github.com/ulissesfc/rateamento-transporte-escolar.git/internal/domain"
)

type SeedsService struct {
}

func (s *SeedsService) SelectSeeds(problem *domain.SeedsProblem) ([]domain.Student, error) {
	if len(problem.Students) == 0 || len(*problem.DurationsMatrix) == 0 {
		return nil, fmt.Errorf("listas de alunos e durações não podem ser vazias")
	}

	// Mapeamentos para busca rápida por índice e para verificar se um ponto já é semente
	alunoIndexMap := make(map[int]int)
	for i, a := range problem.Students {
		alunoIndexMap[a.Id] = i
	}

	jaEhSemente := make(map[int]bool)
	var sementesSelecionadas []domain.Student

	// 1. Adiciona as sementes iniciais (processo manual do artigo)
	for _, id := range problem.StarterSeedsIDs {
		if idx, ok := alunoIndexMap[id]; ok {
			sementesSelecionadas = append(sementesSelecionadas, problem.Students[idx])
			jaEhSemente[idx] = true
		}
	}

	fmt.Printf("Sementes iniciais (manuais): %d\n", len(sementesSelecionadas))

	// 2. Encontra a primeira semente algorítmica (mais distante do depósito central)
	idxDepositoCentral, ok := alunoIndexMap[problem.VirtualDepositID]
	if !ok {
		return nil, fmt.Errorf("ID do depósito central (%d) não encontrado", problem.VirtualDepositID)
	}

	maiorDuracao := -1.0
	idxProximaSemente := -1
	for i := range problem.Students {
		if jaEhSemente[i] {
			continue // Pula se já for semente
		}
		duracao := (*problem.DurationsMatrix)[idxDepositoCentral][i]
		if duracao > maiorDuracao {
			maiorDuracao = duracao
			idxProximaSemente = i
		}
	}

	if idxProximaSemente != -1 {
		sementesSelecionadas = append(sementesSelecionadas, problem.Students[idxProximaSemente])
		jaEhSemente[idxProximaSemente] = true
	}

	// 3. Encontra as sementes restantes (maximizar a soma das distâncias)
	for len(sementesSelecionadas) < problem.TargetSeedsValue {
		maiorSomaDuracoes := -1.0
		idxMelhorCandidato := -1

		// Para cada ponto que ainda não é semente...
		for i := range problem.Students {
			if jaEhSemente[i] {
				continue
			}

			somaDuracoesAtual := 0.0
			// ... some a duração dele para todas as sementes já selecionadas.
			for _, semente := range sementesSelecionadas {
				idxSemente := alunoIndexMap[semente.Id]
				somaDuracoesAtual += (*problem.DurationsMatrix)[i][idxSemente]
			}

			if somaDuracoesAtual > maiorSomaDuracoes {
				maiorSomaDuracoes = somaDuracoesAtual
				idxMelhorCandidato = i
			}
		}

		if idxMelhorCandidato != -1 {
			sementesSelecionadas = append(sementesSelecionadas, problem.Students[idxMelhorCandidato])
			jaEhSemente[idxMelhorCandidato] = true
			fmt.Printf("Semente %d/%d encontrada: Aluno ID %d\n", len(sementesSelecionadas), problem.TargetSeedsValue, problem.Students[idxMelhorCandidato].Id)
		} else {
			break // Não há mais candidatos
		}
	}

	return sementesSelecionadas, nil
}
