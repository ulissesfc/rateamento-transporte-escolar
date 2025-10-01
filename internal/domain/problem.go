package domain

type BusProblem struct {
	Buses       []Bus
	TotalDemand int
	Slack       int
}

type SeedsProblem struct {
	Students         []Student
	DurationsMatrix  *[][]float64
	TargetSeedsValue int
	StarterSeedsIDs  []int
	VirtualDepositID int
}

type ClusterProblem1 struct {
	Students        []Student
	StudentsSeeds   []Student
	DurationsMatrix *[][]float64
}

type SavingProblem struct {
	Students        []Student
	StudentsSeeds   []Student
	DurationsMatrix *[][]float64
	InitialSolution map[int]*Cluster
}

type ClusterProblem2 struct {
	InitialSolution   map[int]*Cluster
	Savings           []Saving
	VehiclesAvailable []Vehicle
}

type TravelProblem struct {
	Clusters []*Cluster
	Fleet    []Vehicle
}

type RouteProblem struct {
	Cluster         *Cluster
	AllStudents     []Student
	GaragemLocation Location
	UFSCLocation    Location
}
