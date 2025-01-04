package Operate

import (
	"fmt"
	"strings"
)

type blindTrig struct{}

// Command camera to take image
func BlindTrig() blindTrig {
	return blindTrig{}
}

func (cmd blindTrig) Compose() string {
	return "T1"
}

type readResult struct{}

// Read last camera image result
func ReadResult() readResult {
	return readResult{}
}
func (cmd readResult) Compose() string {
	return "RT"
}

type trig struct{}

// Command camera to take image and return result
func Trig() trig {
	return trig{}
}

func (cmd trig) Compose() string {
	return "T2"
}

type programRead struct{}

// Read camera's current program
func ProgramRead() programRead {
	return programRead{}
}

func (cmd programRead) Compose() string {
	return "PR"
}

type programWrite struct {
	num int
}

// Set camera program, program number must be in range [0, 127]
func ProgramWrite(num int) (programWrite, error) {
	val := -1
	var err error = nil
	if -1 < num && num < 128 {
		val = num
	} else {
		val = -1
		err = fmt.Errorf("invalid program number %v", num)
	}

	pw := programWrite{
		val,
	}

	return pw, err
}

func (cmd programWrite) Compose() string {
	return fmt.Sprintf("PW,%v", cmd.num)
}

type thresholdRead struct {
	toolNum    int
	upperLimit bool
}

// Read a program limit of the specified tool and limit type (upper/lower)
func ThresholdRead(toolNum int, upperLimit bool) (thresholdRead, error) {
	val := -1
	var err error = nil
	if -1 < toolNum && toolNum < 65 {
		val = toolNum
	} else {
		err = fmt.Errorf("invalid tool number %v", toolNum)
	}

	tr := thresholdRead{
		val,
		upperLimit,
	}

	return tr, err
}

func (cmd thresholdRead) Compose() string {
	b := "0"
	if cmd.upperLimit {
		b = "1"
	}

	return fmt.Sprintf("DR,%v,%v", cmd.toolNum, b)
}

type thresholdWrite struct {
	toolNum    int
	upperLimit bool
	newLimit   int
}

// Set the threshold of the given tool and limit to entered limit
func ThresholdWrite(toolNum int, upperLimit bool, newLimit int) (thresholdWrite, error) {
	tVal := -1
	var err error = nil
	if -1 < toolNum && toolNum < 65 {
		tVal = toolNum
	} else {
		err = fmt.Errorf("invalid tool number %v", toolNum)
	}
	newVal := -1
	if err == nil && -1 < newLimit && newLimit < 10000000 {
		newVal = newLimit
	} else {
		err = fmt.Errorf("invalid threshold value %v", newLimit)
	}
	tw := thresholdWrite{
		toolNum:    tVal,
		upperLimit: upperLimit,
		newLimit:   newVal,
	}

	return tw, err
}

func (cmd thresholdWrite) Compose() string {
	b := "0"
	if cmd.upperLimit {
		b = "1"
	}
	return fmt.Sprintf("DW,%v,%v,%v", cmd.toolNum, b, cmd.newLimit)
}

type textRead struct {
	toolNum int
}

func TextRead(toolNum int) (textRead, error) {
	num := -1
	var err error = nil
	if 0 < toolNum && toolNum < 65 {
		num = toolNum
	} else {
		err = fmt.Errorf("invalid tool number %v, should be in range [1, 64]", toolNum)
	}
	return textRead{num}, err
}

func (cmd textRead) Compose() string {
	return fmt.Sprintf("CR,%v", cmd.toolNum)
}

type textWrite struct {
	toolNum    int
	masterText []string
}

const MAX_TEXT_LEN = 16

// Set master text for specified tool number, text must be less than 16 characters
func TextWrite(toolNum int, text string) (textWrite, error) {
	masterText := make([]string, MAX_TEXT_LEN)
	if len(text) > 16 {
		fmt.Printf("length of input text %v, is too long. Only the first 16 chars will be sent.\n", text)
	}
	var err error = nil
	for i := range masterText {
		masterText[i] = " "
	}
	s := strings.Split(text, "")
	copy(masterText, s)

	num := -1
	if 0 < toolNum && toolNum < 65 {
		num = toolNum
	} else {
		err = fmt.Errorf("invalid tool number %v, should be in range [1, 64]", toolNum)
	}
	tw := textWrite{
		toolNum:    num,
		masterText: masterText,
	}

	return tw, err
}

func (cmd textWrite) Compose() string {
	return fmt.Sprintf("CW,%v,%v", cmd.toolNum, strings.Join(cmd.masterText, ""))
}

// **TODO** Nice to have functions for later :)

/*
type textLengthRead struct {}

func TextLengthRead () {}

func (cmd textLengthRead) Compose() string {}

type textLengthWrite struct {}

func TextLengthWrite () {}

func (cmd textLengthWrite) Compose() string {}

type savedFileRead struct {}

func SavedFileRead() {}

func (cmd savedFileRead) Compose() string {}

type savedFileWrite struct {}

func SavedFileWrite() {}

func (cmd savedFileWrite) Compose() string {}

*/
