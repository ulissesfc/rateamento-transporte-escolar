package domain

type Saving struct {
	Value    float64 // Valor da economia (em segundos de duração)
	StudentI int     // ID do primeiro aluno
	StudentJ int     // ID do segundo aluno
	SeedK    int     // ID da semente para a fusão
}
