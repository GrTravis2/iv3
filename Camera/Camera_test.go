package Camera

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	c := NewCamera([]int{0, 0, 0, 0})

	expectedIp := []int{0, 0, 0, 0}
	ip := c.GetIp()
	for i := range len(c.ip) {
		if expectedIp[i] != ip[i] {
			t.Errorf("expected default value: %v, actual value: ip[%v] = %v\n", expectedIp[i], i, c.ip[i])
		}
	}
	if c.GetPort() != 8500 {
		t.Errorf("Expected port: %v, found %v\n", 8500, c.GetPort())
	}
	if c.GetDelimiter() != '\r' {
		t.Errorf("Expected delimiter: %v, found %v\n", '\r', c.GetDelimiter())
	}

	c2 := NewCamera([]int{192, 168, 1, 1})
	expectedIp = []int{192, 168, 1, 1}
	ip = c2.GetIp()
	for i := range len(c.ip) {
		if expectedIp[i] != ip[i] {
			t.Errorf("expected default value: %v, actual value: ip[%v] = %v\n", expectedIp[i], i, c.ip[i])
		}
	}
}

func TestSetIp(t *testing.T) {
	c := NewCamera([]int{0, 0, 0, 0})
	var tests = []struct {
		values []int
		out    bool
	}{
		//input values, result
		{[]int{0, 0, 0, 0}, true},     // -> good!
		{[]int{-1, 0, 0, 0}, false},   // -> low val
		{[]int{256, 0, 0, 0}, false},  // -> high val
		{[]int{0, 0, 0, 0, 0}, false}, // -> too many vals
		{[]int{0, 0, 0}, false},       // -> too few vals
	}

	for _, input := range tests {
		testName := fmt.Sprintf("input%v,%v", input.values, input.out)
		t.Run(testName, func(t *testing.T) {
			out := c.SetIp(input.values)
			if out != input.out {
				t.Errorf("Expected %v, found %v", input.out, out)
			}
		})
	}
}

func TestSetPort(t *testing.T) {
	c := NewCamera([]int{0, 0, 0, 0})
	low := c.SetPort(0)
	if low == true {
		t.Error("setPort allows value under lower bound - Fail")
	}
	good := c.SetPort(1024)
	if good == false {
		t.Error("setPort does not set good value - Fail")
	}
	high := c.SetPort(65536)
	if high == true {
		t.Error("setPort allows value above upper bound - Fail")
	}
}

func TestSetDelimiter(t *testing.T) {
	//do nothing
}
