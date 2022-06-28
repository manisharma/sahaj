package internal

import "sahaj/pkg/parking"

// Parking represents a Parking Lot
type Parking struct {
	Inventory map[parking.VehicleType]Inventory
	Fee       parking.Fee
}

// Inventory represents actual parking spot
type Inventory struct {
	Total uint
}
