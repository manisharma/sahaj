package main

import (
	"log"
	"os"
	"sahaj/internal"
	"sahaj/pkg/parking"
	parkingFactory "sahaj/pkg/parking_factory"
)

func main() {
	f, err := os.Open("sample.json")
	if err != nil {
		log.Fatalf("os.Open() failed, err:%v", err.Error())
	}
	defer f.Close()

	feeModels, err := parking.GetFeeModels(f)
	if err != nil {
		log.Fatalf("parking.GetFeeModels() failed, err:%v", err.Error())
	}

	feeModel := feeModels[parking.ModelType_Mall]
	parkingLot := parkingFactory.New(parking.ModelType_Mall, feeModel.Fee, map[parking.VehicleType]internal.Inventory{
		parking.VehicleType_Motorcycle: {
			Total: 2,
		},
	})

	// park motorcycle
	result := parkingLot.Do(parking.Action{
		ActionType:  parking.ActionType_Park,
		VehicleType: parking.VehicleType_Motorcycle,
	})
	if result.ParkingReceipt != nil {
		log.Printf("ParkingReceipt : %+v", *result.ParkingReceipt)
	}
	if result.ParkingTicket != nil {
		log.Printf("ParkingTicket : %+v", *result.ParkingTicket)
	}
	if result.Err != nil {
		log.Printf("Err : %+v", result.Err)
	}

	// un park motorcycle
	result = parkingLot.Do(parking.Action{
		ActionType:  parking.ActionType_UnPark,
		VehicleType: parking.VehicleType_Motorcycle,
		TicketNumer: &result.ParkingTicket.TicketNumber,
	})
	if result.ParkingReceipt != nil {
		log.Printf("ParkingReceipt : %+v", *result.ParkingReceipt)
	}
	if result.ParkingTicket != nil {
		log.Printf("ParkingTicket : %+v", *result.ParkingTicket)
	}
	if result.Err != nil {
		log.Printf("Err : %+v", result.Err)
	}
}
