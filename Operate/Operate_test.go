package Operate

import (
	"fmt"
	"testing"
)

func TestBlindTrig(t *testing.T) {
	cmd := BlindTrig()
	actual := cmd.Compose()
	expected := "T1"

	if actual != expected {
		t.Errorf("expected %v, found %v", expected, actual)
	}

	if !cmd.Interpret("T1").Ok() {
		t.Errorf("response error")
	}
}

func TestReadResult(t *testing.T) {
	cmd := ReadResult()
	actual := cmd.Compose()
	expected := "RT"

	if actual != expected {
		t.Errorf("expected %v, found %v", expected, actual)
	}

	if !cmd.Interpret("RT,32767,OK").Ok() {
		t.Errorf("response error")
	}

}

func TestTrig(t *testing.T) {
	cmd := Trig()
	actual := cmd.Compose()
	expected := "T2"

	if actual != expected {
		t.Errorf("expected %v, found %v", expected, actual)
	}

	if !cmd.Interpret("RT,32767,OK").Ok() {
		t.Errorf("response error")
	}
}

func TestProgramRead(t *testing.T) {
	cmd := ProgramRead()
	actual := cmd.Compose()
	expected := "PR"

	if actual != expected {
		t.Errorf("expected %v, found %v", expected, actual)
	}

	if !cmd.Interpret("PR,099").Ok() {
		t.Errorf("response error")
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

	if !cmd.Interpret("PW").Ok() {
		t.Errorf("response error")
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

	cmd, _ := ThresholdRead(0, false)
	if !cmd.Interpret("DR,64,1,9999999").Ok() {
		t.Errorf("response error")
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

	cmd, _ := ThresholdWrite(0, false, 9999999)
	if !cmd.Interpret("DW,64").Ok() {
		t.Errorf("response error")
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

	if !result.Interpret("CR,64, ").Ok() {
		t.Errorf("response error")
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

	if !result.Interpret("CW,64").Ok() {
		t.Errorf("response error")
	}
}
