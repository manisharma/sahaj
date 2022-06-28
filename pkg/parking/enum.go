package parking

import "encoding/json"

type ModelType uint

const (
	ModelType_Mall ModelType = iota + 1
	ModelType_Stadium
	ModelType_Airport
)

func (f ModelType) String() string {
	return [...]string{"", "Mall", "Stadium", "Airport"}[f]
}

func (f *ModelType) FromString(val string) ModelType {
	return map[string]ModelType{
		"Mall":    ModelType_Mall,
		"Stadium": ModelType_Stadium,
		"Airport": ModelType_Airport,
	}[val]
}

func (f ModelType) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

func (f *ModelType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*f = f.FromString(s)
	return nil
}

type VehicleType uint

const (
	VehicleType_Motorcycle VehicleType = iota + 1
	VehicleType_CarSuv
	VehicleType_BusTruck
)

func (s VehicleType) String() string {
	return [...]string{"", "Motorcycle", "Car/Suv", "Bus/Truck"}[s]
}

func (s *VehicleType) FromString(val string) VehicleType {
	return map[string]VehicleType{
		"Motorcycle": VehicleType_Motorcycle,
		"Car/Suv":    VehicleType_CarSuv,
		"Bus/Truck":  VehicleType_BusTruck,
	}[val]
}

func (s VehicleType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *VehicleType) UnmarshalJSON(b []byte) error {
	var v string
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	*s = s.FromString(v)
	return nil
}

type ActionType uint

const (
	ActionType_Park ActionType = iota + 1
	ActionType_UnPark
)

func (s ActionType) String() string {
	return [...]string{"", "Park", "UnPark"}[s]
}

func (s *ActionType) FromString(val string) ActionType {
	return map[string]ActionType{
		"Park":   ActionType_Park,
		"UnPark": ActionType_UnPark,
	}[val]
}

func (s ActionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *ActionType) UnmarshalJSON(b []byte) error {
	var v string
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	*s = s.FromString(v)
	return nil
}

type ChargeType uint

const (
	ChargeType_PerHour ChargeType = iota + 1
	ChargeType_PerDay
)

func (s ChargeType) String() string {
	return [...]string{"", "PerHour", "PerDay"}[s]
}

func (s *ChargeType) FromString(val string) ChargeType {
	return map[string]ChargeType{
		"PerHour": ChargeType_PerHour,
		"PerDay":  ChargeType_PerDay,
	}[val]
}

func (s ChargeType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *ChargeType) UnmarshalJSON(b []byte) error {
	var v string
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	*s = s.FromString(v)
	return nil
}
