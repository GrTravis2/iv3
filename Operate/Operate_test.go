package Operate

import (
	"fmt"
	"testing"
)

func TestBlindTrig(t *testing.T) {
	actual := BlindTrig().Compose()
	expected := "T1"

	if actual != expected {
		t.Errorf("expected %v, found %v", expected, actual)
	}
}

func TestReadResult(t *testing.T) {
	actual := ReadResult().Compose()
	expected := "RT"

	if actual != expected {
		t.Errorf("expected %v, found %v", expected, actual)
	}
}

func TestTrig(t *testing.T) {
	actual := Trig().Compose()
	expected := "T2"

	if actual != expected {
		t.Errorf("expected %v, found %v", expected, actual)
	}
}

func TestProgramRead(t *testing.T) {
	actual := ProgramRead().Compose()
	expected := "PR"

	if actual != expected {
		t.Errorf("expected %v, found %v", expected, actual)
	}
}

func TestProgramWrite(t *testing.T) {
	cmd, err := ProgramWrite(0)
	actual := cmd.Compose()
	expected := "PW,0"

	if actual != expected || err != nil {
		t.Errorf("expected %v, found %v", expected, actual)
	}

	cmd, err = ProgramWrite(-1)
	actual = cmd.Compose()
	expected = "err"
	if err == nil {
		t.Errorf("expected %v, found %v", expected, actual)
	}

	cmd, err = ProgramWrite(128)
	actual = cmd.Compose()
	expected = "err"
	if err == nil {
		t.Errorf("expected %v, found %v", expected, actual)
	}
}

func TestThresholdRead(t *testing.T) {
	var tests = []struct {
		num   int
		upper bool
		ans   string
	}{
		//input values, result
		{0, true, "DR,0,1"},    // -> good!
		{-1, false, "DR,-1,0"}, // -> low val
		{65, true, "DR,-1,1"},  // -> high val
	}

	for _, input := range tests {
		testName := fmt.Sprintf("input%v,%v", input.num, input.ans)
		t.Run(testName, func(t *testing.T) {
			out, _ := ThresholdRead(input.num, input.upper)
			if out.Compose() != input.ans {
				t.Errorf("Expected %v, found %v", input.ans, out.Compose())
			}
		})
	}
}

func TestThresholdWrite(t *testing.T) {
	var tests = []struct {
		num      int
		upper    bool
		newlimit int
		ans      string
	}{
		//input values, result
		{0, true, 0, "DW,0,1,0"},           // -> good!
		{-1, false, -1, "DW,-1,0,-1"},      // -> low val
		{65, true, 10000000, "DW,-1,1,-1"}, // -> high val
	}

	for _, input := range tests {
		testName := fmt.Sprintf("input%v,%v", input.num, input.ans)
		t.Run(testName, func(t *testing.T) {
			out, _ := ThresholdWrite(input.num, input.upper, input.newlimit)
			if out.Compose() != input.ans {
				t.Errorf("Expected %v found %v", input.ans, out.Compose())
			}
		})
	}
}

func TestTextRead(t *testing.T) {
	result, err := TextRead(1)
	ans := "CR,1"
	s := result.Compose()
	if result.Compose() != ans && err != nil {
		t.Errorf("expected %v found %v", ans, s)
	}

	result, err = TextRead(0)
	ans = "CR,-1"
	s = result.Compose()
	if result.Compose() != ans || err == nil {
		t.Errorf("expected %v found %v", ans, s)
	}

	result, err = TextRead(65)
	ans = "CR,-1"
	s = result.Compose()
	if result.Compose() != ans || err == nil {
		t.Errorf("expected %v found %v", ans, s)
	}
}

func TestTextWrite(t *testing.T) {
	result, err := TextWrite(1, "0")
	ans := "CW,1,0               " // empty chars are replaced with ' ' char
	s := result.Compose()
	if result.Compose() != ans && err != nil {
		t.Errorf("expected %v found %v", ans, s)
	}

	result, err = TextWrite(0, "AAAAAAAAAAAAAAAAB") // one extra char, should cut it
	ans = "CW,-1,AAAAAAAAAAAAAAAA"
	s = result.Compose()
	if result.Compose() != ans || err == nil {
		t.Errorf("expected %v found %v", ans, s)
	}

	result, err = TextWrite(65, "AAAAAAAAAAAAAAA") // 15 chars, should add one ' '
	ans = "CW,-1,AAAAAAAAAAAAAAA "
	s = result.Compose()
	if result.Compose() != ans || err == nil {
		t.Errorf("expected %v found %v", ans, s)
	}
}
