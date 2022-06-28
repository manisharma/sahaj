package parking

// ParkingLot represents the contract needed for a parking lot
type ParkingLot interface {
	GetType() ModelType
	Do(action Action) Result
}
