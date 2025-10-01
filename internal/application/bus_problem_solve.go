package application

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ulissesfc/rateamento-transporte-escolar.git/internal/domain"
)

type BusService struct {
	DataPath string // diretório onde ficam problem.mod e solution.txt
}

// Gera o arquivo problem.mod dinamicamente
func (s *BusService) writeProblemFile(problem domain.BusProblem) error {
	// garante que o diretório existe
	if _, err := os.Stat(s.DataPath); os.IsNotExist(err) {
		if err := os.MkdirAll(s.DataPath, 0755); err != nil {
			return err
		}
	}

	file, err := os.Create(s.DataPath + "/problem.mod")
	if err != nil {
		return err
	}
	defer file.Close()

	// 1. Declaração das variáveis (todas primeiro)
	for _, bus := range problem.Buses {
		fmt.Fprintf(file, "var x%d >= 0, <= %d, integer;\n", bus.Type, bus.Quantity)
	}

	// 2. Objetivo
	fmt.Fprint(file, "minimize obj: ")
	for i, bus := range problem.Buses {
		if i > 0 {
			fmt.Fprint(file, " + ")
		}
		fmt.Fprintf(file, "x%d", bus.Type)
	}
	fmt.Fprintln(file, ";")

	// 3. Restrição de demanda
	fmt.Fprint(file, "subject to demand: ")
	for i, bus := range problem.Buses {
		if i > 0 {
			fmt.Fprint(file, " + ")
		}
		fmt.Fprintf(file, "%d*x%d", bus.Capacity-problem.Slack, bus.Type)
	}
	fmt.Fprintf(file, " >= %d;\n", problem.TotalDemand)

	// 4. Restrição especial x1 >= 2
	fmt.Fprintln(file, "s.t. x1_min: x1 >= 2;")

	return nil
}

// Lê solution.txt gerado pelo GLPK
func (s *BusService) readsGeneratedSolution() ([]domain.BusSolution, error) {
	file, err := os.Open(s.DataPath + "/solution.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var solution []domain.BusSolution
	scanner := bufio.NewScanner(file)
	inColumns := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Quando chega na tabela de variáveis
		if strings.HasPrefix(line, "No. Column name") {
			inColumns = true
			// pula linha de cabeçalho
			scanner.Scan()
			continue
		}

		// Processa linhas de variáveis
		if inColumns {
			if line == "" || strings.HasPrefix(line, "Integer feasibility") {
				break // fim da tabela
			}

			fields := strings.Fields(line)
			if len(fields) >= 4 {
				var busType int
				var quantity int
				_, err1 := fmt.Sscanf(fields[1], "x%d", &busType)
				_, err2 := fmt.Sscanf(fields[3], "%d", &quantity) // pega o Activity
				if err1 == nil && err2 == nil && quantity > 0 {
					solution = append(solution, domain.BusSolution{
						Type:     busType,
						Quantity: quantity,
					})
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return solution, nil
}

// Resolve o PPLI usando GLPK no container Docker
func (s *BusService) Solve(problem domain.BusProblem) ([]domain.BusSolution, error) {
	if err := s.writeProblemFile(problem); err != nil {
		return nil, fmt.Errorf("failed to write problem.mod: %w", err)
	}

	cmd := exec.Command(
		"docker", "exec", "glpk",
		"glpsol", "-m", "/app/problem.mod", "-o", "/app/solution.txt",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error running GLPK: %w\n%s", err, string(output))
	}

	solution, err := s.readsGeneratedSolution()
	if err != nil {
		return nil, fmt.Errorf("failed to read solution.txt: %w", err)
	}

	return solution, nil
}
