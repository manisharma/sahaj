package parking_factory

import (
	"sahaj/internal"
	"sahaj/internal/airport"
	"sahaj/internal/mall"
	"sahaj/internal/stadium"
	"sahaj/pkg/parking"
)

// New creates a new Parking Lot
func New(modelType parking.ModelType, fee parking.Fee, inventory map[parking.VehicleType]internal.Inventory) parking.ParkingLot {
	switch modelType {
	case parking.ModelType_Mall:
		return mall.New(fee, inventory)
	case parking.ModelType_Airport:
		return airport.New(fee, inventory)
	case parking.ModelType_Stadium:
		return stadium.New(fee, inventory)
	}
	return nil
}
