package domain

type Cluster struct {
	SeedID      int
	StudentIDs  []int
	TotalDemand int
	Endpoint1ID int
	Endpoint2ID int
	VehicleID   int //ID do veículo alocado a este cluster
}

func (c *Cluster) Reverse() {
	// 1. Inverte a ordem da lista de IDs de alunos (in-place)
	for i, j := 0, len(c.StudentIDs)-1; i < j; i, j = i+1, j-1 {
		c.StudentIDs[i], c.StudentIDs[j] = c.StudentIDs[j], c.StudentIDs[i]
	}

	// 2. Troca os IDs dos endpoints para refletir a nova ordem
	c.Endpoint1ID, c.Endpoint2ID = c.Endpoint2ID, c.Endpoint1ID
}
