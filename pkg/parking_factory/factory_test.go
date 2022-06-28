package parking_factory

import (
	"reflect"
	"sahaj/internal"
	"sahaj/internal/airport"
	"sahaj/internal/mall"
	"sahaj/internal/stadium"
	"sahaj/pkg/parking"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		modelType parking.ModelType
		fee       parking.Fee
		inventory map[parking.VehicleType]internal.Inventory
	}
	tests := []struct {
		name string
		args args
		want parking.ParkingLot
	}{
		{
			name: "ModelType Mall should create Mall Parking Lot",
			args: args{
				modelType: parking.ModelType_Mall,
				fee:       parking.Fee{},
				inventory: map[parking.VehicleType]internal.Inventory{},
			},
			want: &mall.ParkingLot{},
		},
		{
			name: "ModelType Airport should create Airport Parking Lot",
			args: args{
				modelType: parking.ModelType_Airport,
				fee:       parking.Fee{},
				inventory: map[parking.VehicleType]internal.Inventory{},
			},
			want: &airport.ParkingLot{},
		},
		{
			name: "ModelType Stadium should create Stadium Parking Lot",
			args: args{
				modelType: parking.ModelType_Stadium,
				fee:       parking.Fee{},
				inventory: map[parking.VehicleType]internal.Inventory{},
			},
			want: &stadium.ParkingLot{},
		},
		{
			name: "Invalid ModelType should not create any Parking Lot",
			args: args{
				fee:       parking.Fee{},
				inventory: map[parking.VehicleType]internal.Inventory{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.modelType, tt.args.fee, tt.args.inventory); !reflect.DeepEqual(got, tt.want) {
				assert.IsType(t, tt.want, got)
			}
		})
	}
}
