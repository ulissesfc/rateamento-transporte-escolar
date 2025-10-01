package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/ulissesfc/rateamento-transporte-escolar.git/internal/adapter/database"
	osrm "github.com/ulissesfc/rateamento-transporte-escolar.git/internal/adapter/http/orsm"
	"github.com/ulissesfc/rateamento-transporte-escolar.git/internal/application"
	"github.com/ulissesfc/rateamento-transporte-escolar.git/internal/db"
	"github.com/ulissesfc/rateamento-transporte-escolar.git/internal/domain"
)

func PPLI() {

	busProblem := domain.BusProblem{
		Buses: []domain.Bus{
			{Type: 1, Quantity: 3, Capacity: 22},
			{Type: 2, Quantity: 1, Capacity: 22},
			{Type: 3, Quantity: 2, Capacity: 23},
			{Type: 4, Quantity: 4, Capacity: 25},
			{Type: 5, Quantity: 2, Capacity: 28},
			{Type: 6, Quantity: 1, Capacity: 30},
		},
		TotalDemand: 121,
		Slack:       6,
	}

	busService := application.BusService{
		DataPath: "ILP-Solver",
	}

	solution, err := busService.Solve(busProblem)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("--- Resultado da Otimização ---")
		totalBuses := 0
		totalCapacity := 0
		for _, s := range solution {
			for _, bus := range busProblem.Buses {
				if bus.Type == s.Type {
					fmt.Printf("Tipo %d| %d ônibus | Cap. %d | Cap total: %d\n", s.Type, s.Quantity, bus.Capacity, bus.Capacity*s.Quantity)
					totalCapacity += bus.Capacity * s.Quantity
					break
				}
			}
			totalBuses += s.Quantity
		}
		fmt.Printf("Total de ônibus: %d\n", totalBuses)
		fmt.Printf("Capacidade total: %d\n", totalCapacity)
	}
}

func getDurantionsMatrix(repository database.Repository) (*[][]float64, error) {
	students, err := repository.GetStudents(125)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var coords []osrm.Coord
	for _, student := range students {
		//fmt.Printf("lat: %f | long: %f | address: %s\n", student.Latitude, student.Longitude, student.Address)
		coords = append(coords, osrm.Coord{Latitude: student.Latitude, Longitude: student.Longitude})
	}

	osrmClient := osrm.NewClient("http://localhost:5000")
	matrix, err := osrmClient.RequestDurationMatrix(coords)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &matrix, nil
}

func Seeds(students []domain.Student, matrix *[][]float64) []domain.Student {

	targetSeedsValue := 7
	manualSelectSeedsByID := []int{96, 94}
	virtualDepositByID := 109

	seedsProblem := domain.SeedsProblem{
		Students:         students,
		DurationsMatrix:  matrix,
		TargetSeedsValue: targetSeedsValue,
		StarterSeedsIDs:  manualSelectSeedsByID,
		VirtualDepositID: virtualDepositByID,
	}

	seedsService := application.SeedsService{}

	studentsSeeds, err := seedsService.SelectSeeds(&seedsProblem)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	/*
		for _, studentSeed := range studentsSeeds {
			fmt.Println(`"id"=`, studentSeed.Id, "or")
		}
	*/
	return studentsSeeds

}

func Clusters(students []domain.Student, matrix *[][]float64, repository database.Repository) {

	Buses := []domain.Bus{
		{Type: 1, Quantity: 2, Capacity: 22},
		{Type: 2, Quantity: 1, Capacity: 22},
		{Type: 4, Quantity: 1, Capacity: 25},
		{Type: 5, Quantity: 2, Capacity: 28},
		{Type: 6, Quantity: 1, Capacity: 30},
	}

	fleetVehicle := domain.CreateFleetVehicles(Buses)

	clusterService := application.ClusterService{}
	savingService := application.SavingService{}

	studentsSeeds := Seeds(students, matrix)

	//seeds := studentsSeeds
	seeds := []domain.Student{{Id: 96}, {Id: 63}, {Id: 63}, {Id: 83}, {Id: 62}, {Id: 35}, {Id: 28}}

	clusterProblem1 := domain.ClusterProblem1{
		Students:        students,
		StudentsSeeds:   seeds,
		DurationsMatrix: matrix,
	}

	fmt.Println("Calculando Solução inicial...")
	initialSolution := clusterService.CreateInitialSolution(&clusterProblem1)
	fmt.Println("tamanho solução inicial", len(initialSolution))

	savingProblem := domain.SavingProblem{
		Students:        students,
		StudentsSeeds:   studentsSeeds,
		DurationsMatrix: matrix,
		InitialSolution: initialSolution,
	}

	fmt.Println("Calculando Economias...")
	savings := savingService.CalculatesSavings(&savingProblem)
	fmt.Println("tamanho savings", len(savings))

	clusterProblem2 := domain.ClusterProblem2{
		InitialSolution:   initialSolution,
		Savings:           savings,
		VehiclesAvailable: fleetVehicle,
	}

	fmt.Println("Gerando Clusters...")
	clusters := clusterService.GenerateClusters(&clusterProblem2)

	for i, cluster := range clusters {
		fmt.Printf("Cluster %d | id do aluno semente: %d | Demanda: %d \nIDs de Sequenciamento: ", i, cluster.SeedID, cluster.TotalDemand)
		fmt.Println(cluster.StudentIDs)
		fmt.Printf("\n\n")
	}

	travelProblem := domain.TravelProblem{
		Clusters: clusters,
		Fleet:    fleetVehicle,
	}

	travelService := application.TravelService{}
	trips, unallocatedClusters := travelService.AllocateBuses(travelProblem)

	for _, travel := range trips {
		fmt.Println("Veiculo:", travel.Vehicle.Id, "capacidade: ", travel.Vehicle.Capacity, "Numero de alunos alocados", travel.Cluster.TotalDemand)
	}
	for _, cluster := range unallocatedClusters {
		fmt.Println("numero de alunos no cluster não alocado: ", cluster.TotalDemand)
	}

	Routes(students, trips, repository)
}

func Routes(students []domain.Student, trips []domain.Travel, repository database.Repository) {

	garagem := domain.Location{ID: -1, Lat: -26.26832, Lon: -48.85045} // Coordenadas da garagem
	ufsc := domain.Location{ID: -2, Lat: -26.23478, Lon: -48.88377}    // Coordenadas da UFSC Joinville

	osrmClient := osrm.NewClient("http://localhost:5000")
	routeService := application.RouteService{OsrmClient: osrmClient}

	for _, viagem := range trips {
		fmt.Printf("\n--- Gerando Rota para Viagem do Veículo %d (Demanda: %d) ---\n", viagem.Vehicle.Id, viagem.Cluster.TotalDemand)

		rotaFinal, err := routeService.GenerateRouteForCluster(domain.RouteProblem{
			Cluster:         viagem.Cluster,
			AllStudents:     students,
			GaragemLocation: garagem,
			UFSCLocation:    ufsc,
		})
		if err != nil {
			fmt.Printf("Erro ao gerar rota: %v\n", err)
			continue
		}

		fmt.Println("Sequência de Coleta Otimizada:")
		fmt.Println(rotaFinal)

		var coords []osrm.Coord
		selectStudents := []domain.Student{}
		coords = append(coords, osrm.Coord{Latitude: garagem.Lat, Longitude: garagem.Lon})
		for _, studentId := range rotaFinal {
			for _, student := range students {
				if student.Id == studentId {
					coords = append(coords, osrm.Coord{Latitude: student.Latitude, Longitude: student.Longitude})
					selectStudents = append(selectStudents, student)
					break
				}
			}
		}
		coords = append(coords, osrm.Coord{Latitude: ufsc.Lat, Longitude: ufsc.Lon})

		osrmRoute, _ := osrmClient.GetRoutes(coords)
		//fmt.Println("osrm route: ", osrmRoute)

		routeID, err := repository.InsertRoute(db.InsertRouteParams{
			Name:     fmt.Sprintf("viagem do veículo: %d", viagem.Vehicle.Id),
			Polyline: osrmRoute[0].Polyline,
			Distance: float32(osrmRoute[0].Distance),
			Duration: float32(osrmRoute[0].Duration),
		})
		if err != nil {
			fmt.Println(err)
			return
		}

		for i, student := range selectStudents {
			repository.InsertNode(db.InsertNodeParams{
				Name:      student.Name,
				RouteID:   int32(routeID),
				Sequence:  int32(i + 1),
				Latitude:  student.Latitude,
				Longitude: student.Longitude,
			})
		}

	}

}

func NewPostgresConnection(user, pass, host, dbname string) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, pass, host, dbname)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Opcional: testar a conexão
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {

	dbcon, err := NewPostgresConnection("admin", "admin123", "localhost:5436", "transport_routing")
	if err != nil {
		panic(err)
	}
	defer dbcon.Close()

	repo := database.NewRepository(dbcon)
	students, _ := repo.GetStudents(125)

	durantionsMatrix, err := getDurantionsMatrix(*repo)
	if err != nil {
		panic(err)
	}

	PPLI()
	//Seeds(students, durantionsMatrix)
	Clusters(students, durantionsMatrix, *repo)
	//GetUsers()

}
