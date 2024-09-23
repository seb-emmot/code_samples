package tollcalculator

import (
	"testing"
	"time"
)

func AreEqual(a, b int) bool {
	return a == b
}

func TestGetToll(t *testing.T) {
	type test struct {
		input    time.Time
		expected int
	}

	cases := []test{
		test{input: time.Date(2023, 10, 04, 5, 0, 0, 0, time.Local), expected: 0},
		test{input: time.Date(2023, 10, 04, 6, 0, 0, 0, time.Local), expected: 8},
		test{input: time.Date(2023, 10, 04, 6, 30, 0, 0, time.Local), expected: 13},
		test{input: time.Date(2023, 10, 04, 7, 30, 0, 0, time.Local), expected: 18},
		test{input: time.Date(2023, 10, 04, 8, 0, 0, 0, time.Local), expected: 13},
		test{input: time.Date(2023, 10, 04, 12, 0, 0, 0, time.Local), expected: 8},
		test{input: time.Date(2023, 10, 04, 15, 0, 0, 0, time.Local), expected: 13},
		test{input: time.Date(2023, 10, 04, 16, 0, 0, 0, time.Local), expected: 18},
		test{input: time.Date(2023, 10, 04, 17, 0, 0, 0, time.Local), expected: 13},
		test{input: time.Date(2023, 10, 04, 18, 0, 0, 0, time.Local), expected: 8},
		test{input: time.Date(2023, 10, 04, 18, 30, 0, 0, time.Local), expected: 0},
	}

	for _, tCase := range cases {
		act, err := getTollFee(tCase.input, Other)

		if err != nil {
			t.Errorf("error %s", err)
		}

		if !AreEqual(act, tCase.expected) {
			t.Errorf("expected %d got %d for time %s", tCase.expected, act, tCase.input)
		}
	}
}

func TestGetTollFee_Holidays(t *testing.T) {
	for i := 0; i < 24; i++ {
		input := time.Date(2024, 9, 22, i, 0, 0, 0, time.Local) // holiday
		exp := 0

		act, err := getTollFee(input, Other)

		if err != nil {
			t.Errorf("error %s", err)
		}

		if !AreEqual(act, exp) {
			t.Errorf("Expected %d got %d", exp, act)
		}
	}
}

func TestGetTollFee_Vehicles(t *testing.T) {
	for i := 0; i < 7; i++ {
		timestamp := time.Date(2024, 9, 20, 15, 0, 0, 0, time.Local) // weekday

		exp := 0

		// The only vehicletype which should give a tollfee
		if VehicleType(i) == Other {
			exp = 13
		}

		act, err := getTollFee(timestamp, VehicleType(i))

		if err != nil {
			t.Errorf("error %s", err)
		}

		if !AreEqual(act, exp) {
			t.Errorf("Expected %d got %d for Vehicle %d", exp, act, VehicleType(i))
		}
	}
}

func TestGetTollFees_Single(t *testing.T) {
	passes := []time.Time{
		time.Date(2024, 9, 20, 6, 0, 0, 0, time.Local), // 8
	}

	exp := 8

	act, err := GetTollFees(passes, Other)

	if err != nil {
		t.Errorf("error %s", err)
	}

	if !AreEqual(act, exp) {
		t.Errorf("Expected %d got %d", exp, act)
	}
}

func TestGetTollFees_Simple(t *testing.T) {
	passes := []time.Time{
		time.Date(2024, 9, 20, 5, 0, 0, 0, time.Local),
		time.Date(2024, 9, 20, 6, 0, 0, 0, time.Local), // 8
		time.Date(2024, 9, 20, 6, 15, 0, 0, time.Local),
	}

	exp := 8

	act, err := GetTollFees(passes, Other)

	if err != nil {
		t.Errorf("error %s", err)
	}

	if !AreEqual(act, exp) {
		t.Errorf("Expected %d got %d", exp, act)
	}
}

func TestGetTollFees_MultiplePeriods(t *testing.T) {
	passes := []time.Time{
		time.Date(2024, 9, 20, 5, 0, 0, 0, time.Local),
		time.Date(2024, 9, 20, 6, 0, 0, 0, time.Local), // 8
		time.Date(2024, 9, 20, 7, 0, 0, 0, time.Local), // 18
	}

	exp := 26

	act, err := GetTollFees(passes, Other)

	if err != nil {
		t.Errorf("error %s", err)
	}

	if !AreEqual(act, exp) {
		t.Errorf("Expected %d got %d", exp, act)
	}
}

func TestGetTollFees_MaxFee(t *testing.T) {
	passes := []time.Time{
		time.Date(2024, 9, 20, 6, 0, 0, 0, time.Local),
		time.Date(2024, 9, 20, 7, 0, 0, 0, time.Local),
		time.Date(2024, 9, 20, 8, 30, 0, 0, time.Local),
		time.Date(2024, 9, 20, 10, 0, 0, 0, time.Local),
		time.Date(2024, 9, 20, 12, 0, 0, 0, time.Local),
		time.Date(2024, 9, 20, 15, 0, 0, 0, time.Local),
	}

	exp := 60

	act, err := GetTollFees(passes, Other)

	if err != nil {
		t.Errorf("error %s", err)
	}

	if !AreEqual(act, exp) {
		t.Errorf("Expected %d got %d", exp, act)
	}
}

func TestInterval(t *testing.T) {
	i, err := NewTimeInterval("1h", "1h31m", 30)

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
