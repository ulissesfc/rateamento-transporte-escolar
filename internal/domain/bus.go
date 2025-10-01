package domain

type Bus struct {
	Type     int
	Quantity int
	Capacity int
}

type BusSolution struct {
	Type     int
	Quantity int
}

type Vehicle struct {
	Id       int
	Type     int
	Capacity int
	Cluster  *Cluster // O cluster/rota que será atribuído a este veículo
}

type Travel struct {
	Vehicle Vehicle
	Cluster *Cluster
}

func CreateFleetVehicles(Buses []Bus) []Vehicle {
	var fleet []Vehicle
	vehicleID := 1
	for _, busType := range Buses {
		for i := 0; i < busType.Quantity; i++ {
			fleet = append(fleet, Vehicle{
				Id:       vehicleID,
				Type:     busType.Type,
				Capacity: busType.Capacity,
			})
			vehicleID++
		}
	}
	return fleet
}
