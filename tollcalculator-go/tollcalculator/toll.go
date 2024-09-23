package tollcalculator

import (
	"fmt"
	"time"
)

type VehicleType int

const (
	MotorBike VehicleType = iota
	Tractor
	Emergency
	Diplomat
	Foreign
	Miliary
	Other
)

type FeeSpec struct {
	Start, End string
	Fee        int
}

type FeeInterval struct {
	Start time.Duration
	End   time.Duration
	Fee   int
}

// An interval is defined as closed on start and open on end. [interval)
func (t FeeInterval) Contain(c time.Duration) bool {
	return t.Start <= c && c < t.End
}

func NewTimeInterval(start, end string, fee int) (FeeInterval, error) {
	s, err := time.ParseDuration(start)

	if err != nil {
		return FeeInterval{}, err
	}

	e, err := time.ParseDuration(end)

	if err != nil {
		return FeeInterval{}, err
	}

	return FeeInterval{Start: s, End: e, Fee: fee}, nil
}

func GetTollFees(passes []time.Time, v VehicleType) (int, error) {

	if len(passes) == 0 {
		return 0, nil
	}

	start := passes[0]
	pending := 0
	total := 0
	for _, pass := range passes {
		fee, err := getTollFee(pass, v)
		if err != nil {
			return 0, err
		}

		grace, err := time.ParseDuration("1h")

		if err != nil {
			return 0, err
		}

		if start.Add(grace).After(pass) {
			// still in 1 hour period
		} else {
			// 1hr has passed.
			total = total + pending
			pending = 0
			start = pass
		}

		if fee > pending {
			pending = fee
		}
	}

	total += pending

	if total > 60 {
		total = 60
	}

	return total, nil
}

func getTollFee(t time.Time, v VehicleType) (int, error) {
	isTollFree, err := isTollFreeVehicle(v)

	if err != nil {
		return 0, err
	}

	if isTollFreeDate(t) || isTollFree {
		return 0, nil
	}

	// this could be sent from Json or similar.
	specs := []FeeSpec{
		FeeSpec{Start: "0h", End: "6h", Fee: 0},
		FeeSpec{Start: "6h", End: "6h30m", Fee: 8},
		FeeSpec{Start: "6h30m", End: "7h", Fee: 13},
		FeeSpec{Start: "7h", End: "8h", Fee: 18},
		FeeSpec{Start: "8h", End: "8h30m", Fee: 13},
		FeeSpec{Start: "8h30m", End: "15h", Fee: 8},
		FeeSpec{Start: "15h", End: "15h30m", Fee: 13},
		FeeSpec{Start: "15h30m", End: "17h", Fee: 18},
		FeeSpec{Start: "17h", End: "18h", Fee: 13},
		FeeSpec{Start: "18h", End: "18h30m", Fee: 8},
		FeeSpec{Start: "18h30m", End: "24h", Fee: 0},
	}

	intervals := []FeeInterval{}

	for _, spec := range specs {
		i, err := NewTimeInterval(spec.Start, spec.End, spec.Fee)

		if err != nil {
			return 0, err
		}

		intervals = append(intervals, i)
	}

	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	dur := t.Sub(midnight)

	for _, interval := range intervals {
		if interval.Contain(dur) {
			return interval.Fee, nil
		}
	}

	return 0, nil
}

func isTollFreeDate(t time.Time) bool {
	return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
}

func isTollFreeVehicle(v VehicleType) (bool, error) {
	m := map[VehicleType]bool{
		MotorBike: true,
		Tractor:   true,
		Emergency: true,
		Diplomat:  true,
		Foreign:   true,
		Miliary:   true,
		Other:     false,
	}

	res, ok := m[v]

	if !ok {
		return false, fmt.Errorf("not a VehicleType")
	}

	return res, nil
}
