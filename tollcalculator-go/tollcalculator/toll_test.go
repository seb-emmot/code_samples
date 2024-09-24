package tollcalculator

import (
	"testing"
	"time"
)

func AreEqual(a, b int) bool {
	return a == b
}

func isWeekend(t time.Time) bool {
	return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
}

func TestGetToll(t *testing.T) {
	type test struct {
		input    time.Time
		expected int
	}

	cases := []test{
		{input: time.Date(2023, 10, 04, 5, 0, 0, 0, time.Local), expected: 0},
		{input: time.Date(2023, 10, 04, 6, 0, 0, 0, time.Local), expected: 8},
		{input: time.Date(2023, 10, 04, 6, 30, 0, 0, time.Local), expected: 13},
		{input: time.Date(2023, 10, 04, 7, 30, 0, 0, time.Local), expected: 18},
		{input: time.Date(2023, 10, 04, 8, 0, 0, 0, time.Local), expected: 13},
		{input: time.Date(2023, 10, 04, 12, 0, 0, 0, time.Local), expected: 8},
		{input: time.Date(2023, 10, 04, 15, 0, 0, 0, time.Local), expected: 13},
		{input: time.Date(2023, 10, 04, 16, 0, 0, 0, time.Local), expected: 18},
		{input: time.Date(2023, 10, 04, 17, 0, 0, 0, time.Local), expected: 13},
		{input: time.Date(2023, 10, 04, 18, 0, 0, 0, time.Local), expected: 8},
		{input: time.Date(2023, 10, 04, 18, 30, 0, 0, time.Local), expected: 0},
	}

	tc, err := NewTollCalculator(isWeekend)

	if err != nil {
		t.Errorf("error %s", err)
	}

	for _, tCase := range cases {
		act, err := tc.getTollFee(tCase.input, Other)

		if err != nil {
			t.Errorf("error %s", err)
		}

		if !AreEqual(act, tCase.expected) {
			t.Errorf("expected %d got %d for time %s", tCase.expected, act, tCase.input)
		}
	}
}

func TestGetTollFee_Holidays(t *testing.T) {
	tc, err := NewTollCalculator(isWeekend)

	if err != nil {
		t.Errorf("error %s", err)
	}

	for i := 0; i < 24; i++ {
		input := time.Date(2024, 9, 22, i, 0, 0, 0, time.Local) // holiday
		exp := 0

		act, err := tc.getTollFee(input, Other)

		if err != nil {
			t.Errorf("error %s", err)
		}

		if !AreEqual(act, exp) {
			t.Errorf("Expected %d got %d", exp, act)
		}
	}
}

func TestGetTollFee_Vehicles(t *testing.T) {
	tc, err := NewTollCalculator(isWeekend)

	if err != nil {
		t.Errorf("error %s", err)
	}

	for i := 0; i < 7; i++ {
		timestamp := time.Date(2024, 9, 20, 15, 0, 0, 0, time.Local) // weekday

		exp := 0

		// The only vehicletype which should give a tollfee
		if VehicleType(i) == Other {
			exp = 13
		}

		act, err := tc.getTollFee(timestamp, VehicleType(i))

		if err != nil {
			t.Errorf("error %s", err)
		}

		if !AreEqual(act, exp) {
			t.Errorf("Expected %d got %d for Vehicle %d", exp, act, VehicleType(i))
		}
	}
}

func TestGetTollFees_Single(t *testing.T) {
	tc, err := NewTollCalculator(isWeekend)

	if err != nil {
		t.Errorf("error %s", err)
	}

	passes := []time.Time{
		time.Date(2024, 9, 20, 6, 0, 0, 0, time.Local), // 8
	}

	exp := 8

	act, err := tc.GetTollFees(passes, Other)

	if err != nil {
		t.Errorf("error %s", err)
	}

	if !AreEqual(act, exp) {
		t.Errorf("Expected %d got %d", exp, act)
	}
}

func TestGetTollFees_Simple(t *testing.T) {
	tc, err := NewTollCalculator(isWeekend)

	if err != nil {
		t.Errorf("error %s", err)
	}

	passes := []time.Time{
		time.Date(2024, 9, 20, 5, 0, 0, 0, time.Local),
		time.Date(2024, 9, 20, 6, 0, 0, 0, time.Local), // 8
		time.Date(2024, 9, 20, 6, 15, 0, 0, time.Local),
	}

	exp := 8

	act, err := tc.GetTollFees(passes, Other)

	if err != nil {
		t.Errorf("error %s", err)
	}

	if !AreEqual(act, exp) {
		t.Errorf("Expected %d got %d", exp, act)
	}
}

func TestGetTollFees_MultiplePeriods(t *testing.T) {
	tc, err := NewTollCalculator(isWeekend)

	if err != nil {
		t.Errorf("error %s", err)
	}

	passes := []time.Time{
		time.Date(2024, 9, 20, 5, 0, 0, 0, time.Local),
		time.Date(2024, 9, 20, 6, 0, 0, 0, time.Local), // 8
		time.Date(2024, 9, 20, 7, 0, 0, 0, time.Local), // 18
	}

	exp := 26

	act, err := tc.GetTollFees(passes, Other)

	if err != nil {
		t.Errorf("error %s", err)
	}

	if !AreEqual(act, exp) {
		t.Errorf("Expected %d got %d", exp, act)
	}
}

func TestGetTollFees_UnOrdered(t *testing.T) {
	tc, err := NewTollCalculator(isWeekend)

	if err != nil {
		t.Errorf("error %s", err)
	}

	passes := []time.Time{
		time.Date(2024, 9, 20, 7, 0, 0, 0, time.Local), // 18
		time.Date(2024, 9, 20, 6, 0, 0, 0, time.Local), // 8
		time.Date(2024, 9, 20, 5, 0, 0, 0, time.Local),
	}

	exp := 26

	act, err := tc.GetTollFees(passes, Other)

	if err != nil {
		t.Errorf("error %s", err)
	}

	if !AreEqual(act, exp) {
		t.Errorf("Expected %d got %d", exp, act)
	}
}

func TestGetTollFees_MaxFee(t *testing.T) {
	tc, err := NewTollCalculator(isWeekend)

	if err != nil {
		t.Errorf("error %s", err)
	}

	passes := []time.Time{
		time.Date(2024, 9, 20, 6, 0, 0, 0, time.Local),
		time.Date(2024, 9, 20, 7, 0, 0, 0, time.Local),
		time.Date(2024, 9, 20, 8, 30, 0, 0, time.Local),
		time.Date(2024, 9, 20, 10, 0, 0, 0, time.Local),
		time.Date(2024, 9, 20, 12, 0, 0, 0, time.Local),
		time.Date(2024, 9, 20, 15, 0, 0, 0, time.Local),
	}

	exp := 60

	act, err := tc.GetTollFees(passes, Other)

	if err != nil {
		t.Errorf("error %s", err)
	}

	if !AreEqual(act, exp) {
		t.Errorf("Expected %d got %d", exp, act)
	}
}

func TestInterval(t *testing.T) {
	i, err := newTimeInterval("1h", "1h31m", 30)

	if err != nil {
		t.Error(err)
	}

	c, err := time.ParseDuration("1h30m")

	if err != nil {
		t.Error(err)
	}

	act := i.Contain(c)
	exp := true

	if act != exp {
		t.Errorf("Expected %t got %t", exp, act)
	}
}

func TestInterval_Start(t *testing.T) {
	i, err := newTimeInterval("1h", "1h30m", 30)

	if err != nil {
		t.Error(err)
	}

	// start of the interval.
	c, err := time.ParseDuration("1h")

	if err != nil {
		t.Error(err)
	}

	act := i.Contain(c)
	exp := true

	if act != exp {
		t.Errorf("Expected %t got %t", exp, act)
	}
}

func TestInterval_End(t *testing.T) {
	i, err := newTimeInterval("1h", "1h30m", 30)

	if err != nil {
		t.Error(err)
	}

	// end of the interval.
	c, err := time.ParseDuration("1h30m")

	if err != nil {
		t.Error(err)
	}

	act := i.Contain(c)
	exp := false

	if act != exp {
		t.Errorf("Expected %t got %t", exp, act)
	}
}
