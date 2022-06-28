package mall

import (
	"sahaj/internal"
	"sahaj/pkg/parking"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		fee       parking.Fee
		inventory map[parking.VehicleType]internal.Inventory
	}
	tests := []struct {
		name string
		args args
		want *ParkingLot
	}{
		{
			name: "Mall Parking Lot must be created, with specified fee and inventory",
			args: args{
				fee: parking.Fee{
					Charge:   parking.ChargeType_PerHour,
					Vehicles: []parking.Vehicle{},
				},
				inventory: map[parking.VehicleType]internal.Inventory{},
			},
			want: &ParkingLot{
				parking: internal.Parking{
					Inventory: map[parking.VehicleType]internal.Inventory{},
					Fee: parking.Fee{
						Charge:   parking.ChargeType_PerHour,
						Vehicles: []parking.Vehicle{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.fee, tt.args.inventory)
			assert.Equal(t, len(tt.want.parking.Inventory), len(tt.want.parking.Inventory), "supported inventory quantity must match, expected %v got %v", len(tt.want.parking.Inventory), len(tt.want.parking.Inventory))
			assert.Equal(t, tt.want.parking.Fee.Charge, got.parking.Fee.Charge, "ChargeType must match, expected %v got %v", tt.want.parking.Fee.Charge, got.parking.Fee.Charge)
			assert.Equal(t, len(tt.want.parking.Fee.Vehicles), len(got.parking.Fee.Vehicles), "supported vehicle quantity must match, expected %v got %v", len(tt.want.parking.Fee.Vehicles), len(got.parking.Fee.Vehicles))
		})
	}
}

func TestParkingLot_GetType(t *testing.T) {
	tests := []struct {
		name   string
		fields ParkingLot
		want   parking.ModelType
	}{
		{
			name:   "Type of Parking Lot should always be Mall",
			fields: *New(parking.Fee{}, map[parking.VehicleType]internal.Inventory{}),
			want:   parking.ModelType_Mall,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fields.GetType()
			assert.Equal(t, tt.want, got, "must be Mall")
		})
	}
}

func TestParkingLot_Do(t *testing.T) {
	type args struct {
		action parking.Action
	}
	tests := []struct {
		name   string
		fields ParkingLot
		args   args
		want   parking.Result
	}{
		{
			name: "Motercycle should get parked",
			fields: *New(parking.Fee{}, map[parking.VehicleType]internal.Inventory{
				parking.VehicleType_Motorcycle: {
					Total: 1,
				},
			}),
			args: args{
				action: parking.Action{
					ActionType:  parking.ActionType_Park,
					VehicleType: parking.VehicleType_Motorcycle,
				},
			},
			want: parking.Result{
				ParkingTicket: &parking.Ticket{
					TicketNumber:  "001",
					SpotNumber:    1,
					EntryDateTime: internal.Now(),
				},
				ParkingReceipt: nil,
				Err:            nil,
			},
		},
		{
			name: "Motercycle should not get parked if no space in Parking lot",
			fields: *New(parking.Fee{}, map[parking.VehicleType]internal.Inventory{
				parking.VehicleType_Motorcycle: {
					Total: 0,
				},
			}),
			args: args{
				action: parking.Action{
					ActionType:  parking.ActionType_Park,
					VehicleType: parking.VehicleType_Motorcycle,
				},
			},
			want: parking.Result{
				ParkingTicket:  nil,
				ParkingReceipt: nil,
				Err:            parking.ErrNoSpace,
			},
		},
		{
			name: "Motercycle should get un-parked",
			fields: ParkingLot{
				parking: internal.Parking{
					Inventory: map[parking.VehicleType]internal.Inventory{},
					Fee: parking.Fee{
						Charge: parking.ChargeType_PerHour,
						Vehicles: []parking.Vehicle{
							{
								Kind: parking.VehicleType_Motorcycle,
								Rates: []parking.Rate{
									{
										Rate: 10,
									},
								},
							},
						},
					},
				},
				record: map[string]time.Time{
					"001Motorcycle": internal.Now().Add(-1 * time.Hour),
				},
				receiptNo: 0,
				padWidth:  3,
			},
			args: args{
				action: parking.Action{
					ActionType:  parking.ActionType_UnPark,
					VehicleType: parking.VehicleType_Motorcycle,
					TicketNumer: internal.ToStringPtr("001"),
				},
			},
			want: parking.Result{
				ParkingTicket: nil,
				ParkingReceipt: &parking.Receipt{
					ReceiptNumber: "R-001",
					EntryDateTime: internal.Now(),
					ExitDateTime:  internal.Now().Add(1 * time.Hour),
					Fees:          10,
				},
				Err: nil,
			},
		},
		{
			name: "Invalid parking Action",
			fields: *New(parking.Fee{}, map[parking.VehicleType]internal.Inventory{
				parking.VehicleType_Motorcycle: {
					Total: 1,
				},
			}),
			args: args{
				action: parking.Action{
					VehicleType: parking.VehicleType_Motorcycle,
				},
			},
			want: parking.Result{
				ParkingTicket:  nil,
				ParkingReceipt: nil,
				Err:            parking.ErrInvalidAction,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fields.Do(tt.args.action)
			if tt.want.ParkingTicket != nil {
				assert.NotNil(t, got.ParkingTicket, "nil ParkingTicket")
				assert.Nil(t, got.ParkingReceipt, "ParkingReceipt must be nil")
				assert.Nil(t, got.Err, "Err must be nil")
				assert.Equal(t, tt.want.ParkingTicket.SpotNumber, got.ParkingTicket.SpotNumber, "SpotNumber must match, want %v, got %v", tt.want.ParkingTicket.SpotNumber, got.ParkingTicket.SpotNumber)
			}
			if tt.want.ParkingReceipt != nil {
				assert.NotNil(t, got.ParkingReceipt, "nil ParkingReceipt")
				assert.Nil(t, got.ParkingTicket, "ParkingTicket must be nil")
				assert.Nil(t, got.Err, "Err must be nil")
				assert.Equal(t, tt.want.ParkingReceipt.ReceiptNumber, got.ParkingReceipt.ReceiptNumber, "ReceiptNumber must match, want %v, got %v", tt.want.ParkingReceipt.ReceiptNumber, got.ParkingReceipt.ReceiptNumber)
				assert.Equal(t, tt.want.ParkingReceipt.Fees, got.ParkingReceipt.Fees, "Fees must match, want %v, got %v", tt.want.ParkingReceipt.Fees, got.ParkingReceipt.Fees)
			}
			assert.Equal(t, tt.want.Err, got.Err, "Err must match, want %v, got %v", tt.want.Err, got.Err)
		})
	}
}

func Test_calculateFee(t *testing.T) {
	type args struct {
		action    parking.Action
		fee       parking.Fee
		entryTime time.Time
		exitTime  time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    uint
		wantErr bool
	}{
		{
			name: "Motorcycle parked for 3 hours and 30 mins. Fees: 40",
			args: args{
				action: parking.Action{
					ActionType:  parking.ActionType_UnPark,
					VehicleType: parking.VehicleType_Motorcycle,
				},
				fee: parking.Fee{
					Charge: parking.ChargeType_PerHour,
					Vehicles: []parking.Vehicle{
						{
							Kind: parking.VehicleType_Motorcycle,
							Rates: []parking.Rate{
								{
									Rate: 10,
								},
							},
						},
					},
				},
				entryTime: internal.Now(),
				exitTime:  internal.Now().Add(210 * time.Minute),
			},
			want:    40,
			wantErr: false,
		},
		{
			name: "Car parked for 6 hours and 1 min. Fees: 140",
			args: args{
				action: parking.Action{
					ActionType:  parking.ActionType_UnPark,
					VehicleType: parking.VehicleType_CarSuv,
				},
				fee: parking.Fee{
					Charge: parking.ChargeType_PerHour,
					Vehicles: []parking.Vehicle{
						{
							Kind: parking.VehicleType_CarSuv,
							Rates: []parking.Rate{
								{
									Rate: 20,
								},
							},
						},
					},
				},
				entryTime: internal.Now(),
				exitTime:  internal.Now().Add(361 * time.Minute),
			},
			want:    140,
			wantErr: false,
		},
		{
			name: "Truck parked for 1 hour and 59 mins. Fees: 100",
			args: args{
				action: parking.Action{
					ActionType:  parking.ActionType_UnPark,
					VehicleType: parking.VehicleType_BusTruck,
				},
				fee: parking.Fee{
					Charge: parking.ChargeType_PerHour,
					Vehicles: []parking.Vehicle{
						{
							Kind: parking.VehicleType_BusTruck,
							Rates: []parking.Rate{
								{
									Rate: 50,
								},
							},
						},
					},
				},
				entryTime: internal.Now(),
				exitTime:  internal.Now().Add(119 * time.Minute),
			},
			want:    100,
			wantErr: false,
		},
		{
			name: "UnPark Motorcycle in 30 minutes",
			args: args{
				action: parking.Action{
					ActionType:  parking.ActionType_UnPark,
					VehicleType: parking.VehicleType_Motorcycle,
				},
				fee: parking.Fee{
					Charge: parking.ChargeType_PerHour,
					Vehicles: []parking.Vehicle{
						{
							Kind: parking.VehicleType_Motorcycle,
							Rates: []parking.Rate{
								{
									Rate: 10,
								},
							},
						},
					},
				},
				entryTime: internal.Now(),
				exitTime:  internal.Now().Add(30 * time.Minute),
			},
			want:    10,
			wantErr: false,
		},
		{
			name: "UnPark Motorcycle in 1 hour and 30 minutes",
			args: args{
				action: parking.Action{
					ActionType:  parking.ActionType_UnPark,
					VehicleType: parking.VehicleType_Motorcycle,
				},
				fee: parking.Fee{
					Charge: parking.ChargeType_PerHour,
					Vehicles: []parking.Vehicle{
						{
							Kind: parking.VehicleType_Motorcycle,
							Rates: []parking.Rate{
								{
									Rate: 10,
								},
							},
						},
					},
				},
				entryTime: internal.Now(),
				exitTime:  internal.Now().Add(90 * time.Minute),
			},
			want:    20,
			wantErr: false,
		},
		{
			name: "Invalid action",
			args: args{
				action: parking.Action{
					ActionType:  parking.ActionType_Park,
					VehicleType: parking.VehicleType_Motorcycle,
				},
				fee: parking.Fee{
					Charge: parking.ChargeType_PerHour,
					Vehicles: []parking.Vehicle{
						{
							Kind: parking.VehicleType_Motorcycle,
							Rates: []parking.Rate{
								{
									Rate: 10,
								},
							},
						},
					},
				},
				entryTime: internal.Now(),
				exitTime:  internal.Now().Add(30 * time.Minute),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Invalid exit time",
			args: args{
				action: parking.Action{
					ActionType:  parking.ActionType_UnPark,
					VehicleType: parking.VehicleType_Motorcycle,
				},
				fee: parking.Fee{
					Charge: parking.ChargeType_PerHour,
					Vehicles: []parking.Vehicle{
						{
							Kind: parking.VehicleType_Motorcycle,
							Rates: []parking.Rate{
								{
									Rate: 10,
								},
							},
						},
					},
				},
				entryTime: internal.Now(),
				exitTime:  internal.Now().Add(-30 * time.Minute),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Invalid charge",
			args: args{
				action: parking.Action{
					ActionType:  parking.ActionType_UnPark,
					VehicleType: parking.VehicleType_Motorcycle,
				},
				fee: parking.Fee{
					Charge: parking.ChargeType_PerDay,
					Vehicles: []parking.Vehicle{
						{
							Kind: parking.VehicleType_Motorcycle,
							Rates: []parking.Rate{
								{
									Rate: 10,
								},
							},
						},
					},
				},
				entryTime: internal.Now(),
				exitTime:  internal.Now().Add(30 * time.Minute),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Invalid ticket",
			args: args{
				action: parking.Action{
					ActionType:  parking.ActionType_UnPark,
					VehicleType: parking.VehicleType_BusTruck,
				},
				fee: parking.Fee{
					Charge: parking.ChargeType_PerDay,
					Vehicles: []parking.Vehicle{
						{
							Kind: parking.VehicleType_Motorcycle,
							Rates: []parking.Rate{
								{
									Rate: 10,
								},
							},
						},
					},
				},
				entryTime: internal.Now(),
				exitTime:  internal.Now().Add(30 * time.Minute),
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateFee(tt.args.action, tt.args.fee, tt.args.entryTime, tt.args.exitTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateFee() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calculateFee() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getRecordKey(t *testing.T) {
	type args struct {
		ticketNo    string
		vehicleType parking.VehicleType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "001Motorcycle",
			args: args{
				ticketNo:    "001",
				vehicleType: parking.VehicleType_Motorcycle,
			},
			want: "001Motorcycle",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRecordKey(tt.args.ticketNo, tt.args.vehicleType); got != tt.want {
				t.Errorf("getRecordKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
