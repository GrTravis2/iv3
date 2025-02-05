package Operate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/GrTravis2/iv3/Camera"
	"github.com/GrTravis2/iv3/Messenger"
)

type blindTrig struct{}

// Command camera to take image
func BlindTrig() *blindTrig {
	return &blindTrig{}
}

func (cmd *blindTrig) Compose() string {
	return "T1"
}

func (cmd *blindTrig) Interpret(s string) Messenger.Response {
	fields := strings.Split(s, ",")
	r := blindTrigResponse{
		prefix: fields[0],
	}

	return &r
}

type blindTrigResponse struct {
	prefix string
}

func (r *blindTrigResponse) Ok() bool {
	ok := true
	if r.prefix != "T1" {
		ok = false
	}
	return ok
}

type readResult struct{}

// Read last camera image result
func ReadResult() *readResult {
	return &readResult{}
}
func (cmd *readResult) Compose() string {
	return "RT"
}

func (cmd *readResult) Interpret(s string) Messenger.Response {
	err := false
	if s[:strings.Index(s, ",")] == "ER" { // err msg returned
		err = true
		_, s, _ = strings.Cut(s, ",")
	}
	prefix, data, _ := strings.Cut(s, ",")
	r := readResultResponse{
		err:    err,
		prefix: prefix,
		data:   Camera.ToolResult(data),
	}

	return &r
}

type readResultResponse struct {
	err    bool
	prefix string
	data   Camera.ToolResults
}

func (r *readResultResponse) Ok() bool {
	return !r.err
}

type trig struct{}

// Command camera to take image and return result
func Trig() *trig {
	return &trig{}
}

func (cmd *trig) Compose() string {
	return "T2"
}

func (cmd *trig) Interpret(s string) Messenger.Response {
	err := false
	if s[:strings.Index(s, ",")] == "ER" { // err msg returned
		err = true
		_, s, _ = strings.Cut(s, ",")
	}
	prefix, data, _ := strings.Cut(s, ",")
	t := trigResponse{
		err:    err,
		prefix: prefix,
		data:   Camera.ToolResult(data),
	}

	return &t
}

type trigResponse struct {
	err    bool
	prefix string
	data   Camera.ToolResults
}

func (r *trigResponse) Ok() bool {
	return !r.err
}

type programRead struct{}

// Read camera's current program
func ProgramRead() *programRead {
	return &programRead{}
}

func (cmd *programRead) Compose() string {
	return "PR"
}

func (cmd *programRead) Interpret(s string) Messenger.Response {
	data := strings.Split(s, ",")
	num, _ := strconv.Atoi(data[1]) // -> data coming from camera, trust its valid
	pNum, _ := Camera.MakeProgramNumber(num)

	r := programReadResponse{
		prefix:        data[0],
		programNumber: pNum,
	}

	return &r
}

type programReadResponse struct {
	prefix        string
	programNumber Camera.ProgramNumber
}

func (r programReadResponse) Ok() bool {
	ok := true
	if r.prefix != "PR" {
		ok = false
	}
	return ok
}

type programWrite struct {
	num Camera.ProgramNumber
}

// Set camera program, program number must be in range [0, 127]
func ProgramWrite(num int) (*programWrite, error) {
	val, err := Camera.MakeProgramNumber(num)
	if err != nil {
		val = -1
	}
	pw := programWrite{
		val,
	}

	return &pw, err
}

func (cmd *programWrite) Compose() string {
	return fmt.Sprintf("PW,%v", cmd.num)
}

func (cmd *programWrite) Interpret(s string) Messenger.Response {
	err := false
	if s[:2] == "ER" { // err msg returned
		err = true
		_, s, _ = strings.Cut(s, ",")
	}
	prefix, _, _ := strings.Cut(s, ",")
	r := programWriteResponse{
		err:    err,
		prefix: prefix,
	}

	return &r
}

type programWriteResponse struct {
	err    bool
	prefix string
}

func (r *programWriteResponse) Ok() bool {
	return !r.err
}

type thresholdRead struct {
	toolNum    Camera.ToolNumber
	upperLimit Camera.UpperLimit
}

// Read a program limit of the specified tool and limit type (upper/lower)
// Ranges: toolNum [0, 64]
func ThresholdRead(toolNum int, upperLimit bool) (*thresholdRead, error) {
	num, err := Camera.MakeToolNumber(toolNum)
	if err != nil {
		num = -1
	}

	tr := thresholdRead{
		num,
		Camera.MakeUpperLimit(upperLimit),
	}

	return &tr, err
}

func (cmd *thresholdRead) Compose() string {
	b := "0"
	if cmd.upperLimit {
		b = "1"
	}

	return fmt.Sprintf("DR,%v,%v", cmd.toolNum, b)
}

func (cmd *thresholdRead) Interpret(s string) Messenger.Response {
	data := strings.Split(s, ",")
	prefix := data[0]

	num, _ := strconv.Atoi(data[1])
	toolNum, _ := Camera.MakeToolNumber(num)

	var upper bool
	if data[2] == "0" {
		upper = false
	} else if data[2] == "1" {
		upper = true
	}

	tLimit, _ := strconv.Atoi(data[3])
	threshold, _ := Camera.MakeThreshold(tLimit)

	r := thresholdReadResponse{
		prefix:         prefix,
		toolNum:        toolNum,
		upperLimit:     Camera.MakeUpperLimit(upper),
		thresholdLimit: threshold,
	}

	return &r
}

type thresholdReadResponse struct { //docs says it wont error, can always add later if needed
	prefix         string
	toolNum        Camera.ToolNumber
	upperLimit     Camera.UpperLimit
	thresholdLimit Camera.Threshold
}

func (r *thresholdReadResponse) Ok() bool {
	ok := false
	if r.prefix == "DR" {
		ok = true
	}

	return ok
}

type thresholdWrite struct {
	toolNum    Camera.ToolNumber
	upperLimit Camera.UpperLimit
	newLimit   Camera.Threshold
}

// Set the threshold of the given tool and limit to entered limit
// Ranges: toolNum [0, 64], newLimit [0, 9999999]
func ThresholdWrite(toolNum int, upperLimit bool, newLimit int) (*thresholdWrite, error) {
	var tNum Camera.ToolNumber
	var threshold Camera.Threshold
	var positionToolErr, toolNumErr, thresholdErr error = nil, nil, nil
	if toolNum == 0 {
		toolNum = -1
		positionToolErr = fmt.Errorf("invalid tool number, cannot target position tool 00")
	} else {
		tNum, toolNumErr = Camera.MakeToolNumber(toolNum)
		threshold, thresholdErr = Camera.MakeThreshold(newLimit)
	}
	tw := thresholdWrite{
		toolNum:    tNum,
		upperLimit: Camera.MakeUpperLimit(upperLimit),
		newLimit:   threshold,
	}

	err := errors.Join(positionToolErr, toolNumErr, thresholdErr)

	return &tw, err
}

func (cmd *thresholdWrite) Compose() string {
	b := "0"
	if cmd.upperLimit {
		b = "1"
	}
	return fmt.Sprintf("DW,%v,%v,%v", cmd.toolNum, b, cmd.newLimit)
}

func (cmd *thresholdWrite) Interpret(s string) Messenger.Response {
	data := strings.Split(s, ",")
	prefix := data[0]
	num, _ := strconv.Atoi(data[1])
	toolNum, _ := Camera.MakeToolNumber(num)

	r := thresholdWriteResponse{
		prefix:  prefix,
		toolNum: toolNum,
	}

	return &r
}

type thresholdWriteResponse struct {
	prefix  string
	toolNum Camera.ToolNumber
}

func (r *thresholdWriteResponse) Ok() bool {
	ok := false
	if r.prefix == "DW" {
		ok = true
	}

	return ok
}

type textRead struct {
	toolNum Camera.ToolNumber
}

// Read current master text value for specified tool in range [1, 64]
func TextRead(toolNum int) (*textRead, error) {
	var tNum Camera.ToolNumber
	var err error = nil
	if toolNum == 0 {
		tNum = -1
		err = fmt.Errorf("invalid tool number, cannot target position tool 00")
	} else {
		tNum, err = Camera.MakeToolNumber(toolNum)

	}

	return &textRead{tNum}, err
}

func (cmd *textRead) Compose() string {
	return fmt.Sprintf("CR,%v", cmd.toolNum)
}

func (cmd *textRead) Interpret(s string) Messenger.Response {
	data := strings.Split(s, ",")
	prefix := data[0]
	num, _ := strconv.Atoi(data[1])
	tNum, _ := Camera.MakeToolNumber(num)
	mText := Camera.MakeMasterText(data[2])
	r := textReadResponse{
		prefix:     prefix,
		toolNum:    tNum,
		masterText: mText,
	}

	return &r
}

type textReadResponse struct {
	prefix     string
	toolNum    Camera.ToolNumber
	masterText Camera.MasterText
}

func (r *textReadResponse) Ok() bool {
	ok := false
	if r.prefix == "CR" {
		ok = true
	}

	return ok
}

type textWrite struct {
	toolNum    Camera.ToolNumber
	masterText Camera.MasterText
}

// Set master text for specified tool number, text must be less than 16 characters
func TextWrite(toolNum int, text string) (*textWrite, error) {
	var tNum Camera.ToolNumber
	var positionToolErr, err error = nil, nil
	if toolNum == 0 {
		tNum = -1
		positionToolErr = fmt.Errorf("invalid tool number, cannot target position tool 00")
	} else {
		tNum, err = Camera.MakeToolNumber(toolNum)
	}
	mText := Camera.MakeMasterText(text)
	tw := textWrite{
		toolNum:    tNum,
		masterText: mText,
	}

	return &tw, errors.Join(positionToolErr, err)
}

func (cmd *textWrite) Compose() string {
	return fmt.Sprintf("CW,%v,%v", cmd.toolNum, cmd.masterText)
}

func (cmd *textWrite) Interpret(s string) Messenger.Response {
	data := strings.Split(s, ",")
	prefix := data[0]
	num, _ := strconv.Atoi(data[1])
	tNum, _ := Camera.MakeToolNumber(num)

	r := textWriteResponse{
		prefix:  prefix,
		toolNum: tNum,
	}

	return &r
}

type textWriteResponse struct {
	prefix  string
	toolNum Camera.ToolNumber
}

func (r *textWriteResponse) Ok() bool {
	ok := false
	if r.prefix == "CW" {
		ok = true
	}

	return ok
}

// **TODO** Nice to have functions for later :)

/*

type textLengthRead struct {}

// Get master text when scanning for length of text mode
func TextLengthRead() textLengthRead {}

func (cmd textLengthRead) Compose() string {}

type textLengthWrite struct {}

// Set length of text to be scanned for
func TextLengthWrite() textLengthWrite {}

func (cmd textLengthWrite) Compose() string {}

type savedFileRead struct {}

// Read the name of the file used to configure image transfer thru FTP/SD
func SavedFileRead() savedFileRead {}

func (cmd savedFileRead) Compose() string {}

type savedFileWrite struct {}

// Set the name of the file used to configure image transfer thru FTP/SD
func SavedFileWrite() savedFileWrite {}

func (cmd savedFileWrite) Compose() string {}

type registerMasterImage struct {}

// Update the program master image to current live image
func RegisterMasterImage() registerMasterImage {}

func (cmd registerMasterImage) Compose() string {}

*/
