package airport

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
			name: "Airport Parking Lot must be created, with specified fee and inventory",
			args: args{
				fee: parking.Fee{
					Charge:   parking.ChargeType_PerDay,
					Vehicles: []parking.Vehicle{},
				},
				inventory: map[parking.VehicleType]internal.Inventory{},
			},
			want: &ParkingLot{
				parking: internal.Parking{
					Inventory: map[parking.VehicleType]internal.Inventory{},
					Fee: parking.Fee{
						Charge:   parking.ChargeType_PerDay,
						Vehicles: []parking.Vehicle{},
					},
				},
			},
		},
		{
			name: "trying to instantiate Airport Parking Lot with Bus/Truck as Vehicle should panic",
			args: args{
				fee: parking.Fee{
					Charge:   parking.ChargeType_PerDay,
					Vehicles: []parking.Vehicle{},
				},
				inventory: map[parking.VehicleType]internal.Inventory{
					parking.VehicleType_BusTruck: {
						Total: 1000,
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != nil {
				got := New(tt.args.fee, tt.args.inventory)
				assert.Equal(t, len(tt.want.parking.Inventory), len(tt.want.parking.Inventory), "supported inventory quantity must match, expected %v got %v", len(tt.want.parking.Inventory), len(tt.want.parking.Inventory))
				assert.Equal(t, tt.want.parking.Fee.Charge, got.parking.Fee.Charge, "ChargeType must match, expected %v got %v", tt.want.parking.Fee.Charge, got.parking.Fee.Charge)
				assert.Equal(t, len(tt.want.parking.Fee.Vehicles), len(got.parking.Fee.Vehicles), "supported vehicle quantity must match, expected %v got %v", len(tt.want.parking.Fee.Vehicles), len(got.parking.Fee.Vehicles))
			} else {
				bool := assert.Panics(t, func() { New(tt.args.fee, tt.args.inventory) }, "must panic")
				assert.True(t, bool, "New() should have panicked")
			}
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
			name:   "Type of Parking Lot should always be Airport",
			fields: *New(parking.Fee{}, map[parking.VehicleType]internal.Inventory{}),
			want:   parking.ModelType_Airport,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fields.GetType()
			assert.Equal(t, tt.want, got, "must be Airport")
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
						Charge: parking.ChargeType_PerDay,
						Vehicles: []parking.Vehicle{
							{
								Kind: parking.VehicleType_Motorcycle,
								Rates: []parking.Rate{
									{
										Rate: 0,
									},
								},
							},
						},
					},
				},
				record: map[string]time.Time{
					"001Motorcycle": internal.Now().Add(-55 * time.Minute),
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
					ExitDateTime:  internal.Now().Add(-55 * time.Minute),
					Fees:          0,
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
		{
			name: "only PerDay Charge is supported",
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
										Rate: 0,
									},
								},
							},
						},
					},
				},
				record: map[string]time.Time{
					"001Motorcycle": internal.Now().Add(-55 * time.Minute),
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
				ParkingTicket:  nil,
				ParkingReceipt: nil,
				Err:            parking.ErrChargeNotSupported,
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
			name: "Motorcycle parked for 55 mins. Fees: 0",
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
									From: 0,
									Till: 1,
									Rate: 0,
								},
								{
									From: 1,
									Till: 8,
									Rate: 40,
								},
								{
									From: 8,
									Till: 24,
									Rate: 60,
								},
								{
									From: 24,
									Till: 0,
									Rate: 80,
								},
							},
						},
					},
				},
				entryTime: internal.Now(),
				exitTime:  internal.Now().Add(55 * time.Minute),
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "Motorcycle parked for 14 hours and 59 mins. Fees: 60",
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
									From: 0,
									Till: 1,
									Rate: 0,
								},
								{
									From: 1,
									Till: 8,
									Rate: 40,
								},
								{
									From: 8,
									Till: 24,
									Rate: 60,
								},
								{
									From: 24,
									Till: 0,
									Rate: 80,
								},
							},
						},
					},
				},
				entryTime: internal.Now(),
				exitTime:  internal.Now().Add(899 * time.Minute),
			},
			want:    60,
			wantErr: false,
		},
		{
			name: "Motorcycle parked for 1 day and 12 hours. Fees: 160",
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
									From: 0,
									Till: 1,
									Rate: 0,
								},
								{
									From: 1,
									Till: 8,
									Rate: 40,
								},
								{
									From: 8,
									Till: 24,
									Rate: 60,
								},
								{
									From: 24,
									Till: 0,
									Rate: 80,
								},
							},
						},
					},
				},
				entryTime: internal.Now(),
				exitTime:  internal.Now().Add(2160 * time.Minute),
			},
			want:    160,
			wantErr: false,
		},
		{
			name: "Car parked for 50 mins. Fees: 60",
			args: args{
				action: parking.Action{
					ActionType:  parking.ActionType_UnPark,
					VehicleType: parking.VehicleType_CarSuv,
				},
				fee: parking.Fee{
					Charge: parking.ChargeType_PerDay,
					Vehicles: []parking.Vehicle{
						{
							Kind: parking.VehicleType_CarSuv,
							Rates: []parking.Rate{
								{
									From: 12,
									Till: 24,
									Rate: 80,
								},
								{
									From: 24,
									Till: 0,
									Rate: 100,
								},
								{
									From: 0,
									Till: 12,
									Rate: 60,
								},
							},
						},
					},
				},
				entryTime: internal.Now(),
				exitTime:  internal.Now().Add(50 * time.Minute),
			},
			want:    60,
			wantErr: false,
		},
		{
			name: "SUV parked for 23 hours and 59 mins. Fees: 80",
			args: args{
				action: parking.Action{
					ActionType:  parking.ActionType_UnPark,
					VehicleType: parking.VehicleType_CarSuv,
				},
				fee: parking.Fee{
					Charge: parking.ChargeType_PerDay,
					Vehicles: []parking.Vehicle{
						{
							Kind: parking.VehicleType_CarSuv,
							Rates: []parking.Rate{
								{
									From: 12,
									Till: 24,
									Rate: 80,
								},
								{
									From: 24,
									Till: 0,
									Rate: 100,
								},
								{
									From: 0,
									Till: 12,
									Rate: 60,
								},
							},
						},
					},
				},
				entryTime: internal.Now(),
				exitTime:  internal.Now().Add(1439 * time.Minute),
			},
			want:    80,
			wantErr: false,
		},
		{
			name: "Car parked for 3 days and 1 hour. Fees: 400",
			args: args{
				action: parking.Action{
					ActionType:  parking.ActionType_UnPark,
					VehicleType: parking.VehicleType_CarSuv,
				},
				fee: parking.Fee{
					Charge: parking.ChargeType_PerDay,
					Vehicles: []parking.Vehicle{
						{
							Kind: parking.VehicleType_CarSuv,
							Rates: []parking.Rate{
								{
									From: 12,
									Till: 24,
									Rate: 80,
								},
								{
									From: 24,
									Till: 0,
									Rate: 100,
								},
								{
									From: 0,
									Till: 12,
									Rate: 60,
								},
							},
						},
					},
				},
				entryTime: internal.Now(),
				exitTime:  internal.Now().Add(4380 * time.Minute),
			},
			want:    400,
			wantErr: false,
		},
		{
			name: "Truck parking not allowed",
			args: args{
				action: parking.Action{
					ActionType:  parking.ActionType_UnPark,
					VehicleType: parking.VehicleType_BusTruck,
				},
				fee: parking.Fee{
					Charge:   parking.ChargeType_PerDay,
					Vehicles: []parking.Vehicle{},
				},
				entryTime: internal.Now(),
				exitTime:  internal.Now().Add(4380 * time.Minute),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "PerHour rates not allowed",
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
									From: 12,
									Till: 24,
									Rate: 80,
								},
								{
									From: 24,
									Till: 0,
									Rate: 100,
								},
								{
									From: 0,
									Till: 12,
									Rate: 60,
								},
							},
						},
					},
				},
				entryTime: internal.Now(),
				exitTime:  internal.Now().Add(4380 * time.Minute),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Fee can not be calculated during parking",
			args: args{
				action: parking.Action{
					ActionType:  parking.ActionType_Park,
					VehicleType: parking.VehicleType_CarSuv,
				},
				fee: parking.Fee{
					Charge: parking.ChargeType_PerHour,
					Vehicles: []parking.Vehicle{
						{
							Kind: parking.VehicleType_CarSuv,
							Rates: []parking.Rate{
								{
									From: 12,
									Till: 24,
									Rate: 80,
								},
								{
									From: 24,
									Till: 0,
									Rate: 100,
								},
								{
									From: 0,
									Till: 12,
									Rate: 60,
								},
							},
						},
					},
				},
				entryTime: internal.Now(),
				exitTime:  internal.Now().Add(4380 * time.Minute),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Exit time can not be before Entry time",
			args: args{
				action: parking.Action{
					ActionType:  parking.ActionType_Park,
					VehicleType: parking.VehicleType_CarSuv,
				},
				fee: parking.Fee{
					Charge: parking.ChargeType_PerHour,
					Vehicles: []parking.Vehicle{
						{
							Kind: parking.VehicleType_CarSuv,
							Rates: []parking.Rate{
								{
									From: 12,
									Till: 24,
									Rate: 80,
								},
								{
									From: 24,
									Till: 0,
									Rate: 100,
								},
								{
									From: 0,
									Till: 12,
									Rate: 60,
								},
							},
						},
					},
				},
				entryTime: internal.Now(),
				exitTime:  internal.Now().Add(-1 * time.Nanosecond),
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
			name: "009Car/Suv",
			args: args{
				ticketNo:    "009",
				vehicleType: parking.VehicleType_CarSuv,
			},
			want: "009Car/Suv",
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
