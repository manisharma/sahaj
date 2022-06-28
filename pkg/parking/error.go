package parking

import "errors"

var (
	ErrNoSpace            = errors.New(" No space available")
	ErrInvalidAction      = errors.New(" Invalid action")
	ErrInvalidTicket      = errors.New(" Invalid ticket")
	ErrExitTime           = errors.New(" Invalid exit time")
	ErrChargeNotSupported = errors.New(" Invalid charge type not supported")
	ErrVehicleNotAllowed  = errors.New(" The vehicle is not allowed to be parked")
	ErrVehicleMismatch    = errors.New(" The vehicle on ticket is not the vehicle which was parked")
)
