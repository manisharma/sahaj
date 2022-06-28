package parking

import (
	"time"
)

// Action encapsulates a basic opration on a Parking Lot
type Action struct {
	ActionType  ActionType
	VehicleType VehicleType
	TicketNumer *string
}

// Result encapsulates result of an Action on a Parking Lot
type Result struct {
	ParkingTicket  *Ticket
	ParkingReceipt *Receipt
	Err            error
}

// Ticket represents a Parking ticket
type Ticket struct {
	TicketNumber  string
	SpotNumber    uint
	EntryDateTime time.Time
}

// Receipt represents a receipt a User recieves after surrendring the Parking Ticket
type Receipt struct {
	ReceiptNumber               string
	EntryDateTime, ExitDateTime time.Time
	Fees                        uint
}

// FeeModels encapsulates all FeeModel on which a Parking Lot works
type FeeModels map[ModelType]FeeModel

// FeeModel represents a basic unit of Parking Lot
type FeeModel struct {
	Model ModelType `json:"model"`
	Fee   Fee       `json:"Fee"`
}

type Fee struct {
	Charge   ChargeType `json:"charge"`
	Vehicles []Vehicle  `json:"vehicles"`
}

type Vehicle struct {
	Kind  VehicleType `json:"kind"`
	Rates []Rate      `json:"rates"`
}

type Rate struct {
	From uint `json:"from"`
	Till uint `json:"till"`
	Rate uint `json:"rate"`
}

// Rates can be jumbled, we need to traverse from lowest to highest time
type SortRatesByStartTime []Rate

func (a SortRatesByStartTime) Len() int           { return len(a) }
func (a SortRatesByStartTime) Less(i, j int) bool { return a[i].From < a[j].From }
func (a SortRatesByStartTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
