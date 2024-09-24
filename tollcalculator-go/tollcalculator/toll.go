package tollcalculator

import (
	"fmt"
	"slices"
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

type TollCalculator struct {
	intervals       []FeeInterval
	gracePeriod     string
	holidayProvider func(time.Time) bool
}

func newTimeInterval(start, end string, fee int) (FeeInterval, error) {
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

// Returns a new TollCalculator using the provided holidayProvider implementation.
func NewTollCalculator(holidayProvider func(time.Time) bool) (TollCalculator, error) {
	tc := TollCalculator{
		gracePeriod:     "1h",
		holidayProvider: holidayProvider,
	}
	specs := []FeeSpec{
		{Start: "0h", End: "6h", Fee: 0},
		{Start: "6h", End: "6h30m", Fee: 8},
		{Start: "6h30m", End: "7h", Fee: 13},
		{Start: "7h", End: "8h", Fee: 18},
		{Start: "8h", End: "8h30m", Fee: 13},
		{Start: "8h30m", End: "15h", Fee: 8},
		{Start: "15h", End: "15h30m", Fee: 13},
		{Start: "15h30m", End: "17h", Fee: 18},
		{Start: "17h", End: "18h", Fee: 13},
		{Start: "18h", End: "18h30m", Fee: 8},
		{Start: "18h30m", End: "24h", Fee: 0},
	}

	intervals := []FeeInterval{}

	for _, spec := range specs {
		i, err := newTimeInterval(spec.Start, spec.End, spec.Fee)

		if err != nil {
			return TollCalculator{}, err
		}

		intervals = append(intervals, i)
	}

	tc.intervals = intervals

	return tc, nil
}

// Calculates total tolls attributable to a vehicle type given an array of passes.
func (tc TollCalculator) GetTollFees(passes []time.Time, v VehicleType) (int, error) {

	slices.SortFunc(passes, func(a, b time.Time) int {
		return a.Compare(b)
	})

	if len(passes) == 0 {
		return 0, nil
	}

	start := passes[0]
	pending := 0
	total := 0
	for _, pass := range passes {
		fee, err := tc.getTollFee(pass, v)
		if err != nil {
			return 0, err
		}

		grace, err := time.ParseDuration(tc.gracePeriod)

		if err != nil {
			return 0, err
		}

		if start.Add(grace).After(pass) {
			// still in 1 hour period
		} else {
			// grace period has passed.
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

func (tc TollCalculator) getTollFee(t time.Time, v VehicleType) (int, error) {
	isTollFree, err := isTollFreeVehicle(v)

	if err != nil {
		return 0, err
	}

	if tc.isTollFreeDate(t) || isTollFree {
		return 0, nil
	}

	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	dur := t.Sub(midnight)

	for _, interval := range tc.intervals {
		if interval.Contain(dur) {
			return interval.Fee, nil
		}
	}

	return 0, nil
}

func (tc TollCalculator) isTollFreeDate(t time.Time) bool {
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
