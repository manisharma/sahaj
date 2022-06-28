package mall

import (
	"fmt"
	"sahaj/internal"
	"sahaj/pkg/parking"
	"sort"
	"strings"
	"time"
)

type ParkingLot struct {
	parking   internal.Parking     // base parking model
	record    map[string]time.Time // parking record
	receiptNo uint                 // tracks upcoming receipt
	padWidth  uint                 // for printing receipt & ticket number
}

func New(fee parking.Fee, inventory map[parking.VehicleType]internal.Inventory) *ParkingLot {
	return &ParkingLot{
		parking: internal.Parking{
			Inventory: inventory,
			Fee:       fee,
		},
		record:    map[string]time.Time{},
		receiptNo: 0,
		padWidth:  3,
	}
}

func (p *ParkingLot) GetType() parking.ModelType {
	return parking.ModelType_Mall
}

func (p *ParkingLot) Do(action parking.Action) parking.Result {
	res := parking.Result{}
	switch action.ActionType {
	case parking.ActionType_Park:
		res.ParkingTicket, res.Err = p.generateParkingTicket(action)
	case parking.ActionType_UnPark:
		res.ParkingReceipt, res.Err = p.generateParkingReceipt(action)
	default:
		res.Err = parking.ErrInvalidAction
	}
	return res
}

func (p *ParkingLot) generateParkingTicket(action parking.Action) (*parking.Ticket, error) {
	inv, ok := p.parking.Inventory[action.VehicleType]
	if !ok {
		return nil, parking.ErrVehicleNotAllowed
	}
	if len(p.record) < int(inv.Total) {
		tktNo := fmt.Sprintf(fmt.Sprintf("%%0%dd", p.padWidth), len(p.record)+1)
		key := getRecordKey(tktNo, action.VehicleType)
		p.record[key] = internal.Now()
		return &parking.Ticket{
			TicketNumber:  tktNo,
			SpotNumber:    uint(len(p.record)),
			EntryDateTime: internal.Now(),
		}, nil
	}
	return nil, parking.ErrNoSpace
}

func (p *ParkingLot) generateParkingReceipt(action parking.Action) (*parking.Receipt, error) {
	if action.TicketNumer == nil {
		return nil, parking.ErrInvalidTicket
	}
	key := getRecordKey(*action.TicketNumer, action.VehicleType)
	if parkingTime, ok := p.record[key]; ok {
		fee, err := calculateFee(action, p.parking.Fee, parkingTime, internal.Now())
		if err != nil {
			return nil, err
		}
		p.receiptNo++
		receipt := &parking.Receipt{
			ReceiptNumber: fmt.Sprintf(fmt.Sprintf("R-%%0%dd", p.padWidth), p.receiptNo),
			EntryDateTime: parkingTime,
			ExitDateTime:  internal.Now(),
			Fees:          fee,
		}
		delete(p.record, key)
		return receipt, nil
	}
	return nil, parking.ErrVehicleMismatch
}

func calculateFee(action parking.Action, fee parking.Fee, entryTime, exitTime time.Time) (uint, error) {
	var fees uint
	if action.ActionType != parking.ActionType_UnPark {
		return fees, parking.ErrInvalidAction
	}
	if exitTime.Before(entryTime) {
		return fees, parking.ErrExitTime
	}
	parkingDuration := exitTime.Sub(entryTime)
	switch fee.Charge {
	case parking.ChargeType_PerHour: // only PerHour Charge is supported by Mall
		var rates []parking.Rate
		for _, vehicle := range fee.Vehicles {
			if vehicle.Kind == action.VehicleType {
				rates = vehicle.Rates
				break
			}
		}
		if len(rates) == 0 {
			return fees, parking.ErrInvalidTicket
		}

		// make sure the rates are in increasing order of from time
		ratesSortedByTime := parking.SortRatesByStartTime(rates)
		sort.Sort(ratesSortedByTime)

		hr := uint(parkingDuration.Hours())
		min := uint(parkingDuration.Minutes())
		min = min - hr*60
		fees = hr * ratesSortedByTime[0].Rate
		if min > 0 {
			fees += ratesSortedByTime[0].Rate
		}
	case parking.ChargeType_PerDay:
		fallthrough
	default:
		return fees, parking.ErrChargeNotSupported
	}
	return fees, nil
}

func getRecordKey(ticketNo string, vehicleType parking.VehicleType) string {
	var key strings.Builder
	key.WriteString(ticketNo)
	key.WriteString(vehicleType.String())
	return key.String()
}
