package Operate

import (
	"fmt"
	"strconv"
	"strings"
)

type blindTrig struct{}

// Command camera to take image
func BlindTrig() *blindTrig {
	return &blindTrig{}
}

func (cmd *blindTrig) Compose() string {
	return "T1"
}

func (cmd *blindTrig) Interpret(s string) *blindTrigResponse {
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

func (cmd *readResult) Interpret(s string) *readResultResponse {
	err := false
	if s[:strings.Index(s, ",")] == "ER" { // err msg returned
		err = true
		_, s, _ = strings.Cut(s, ",")
	}
	prefix, data, _ := strings.Cut(s, ",")
	r := readResultResponse{
		err:    err,
		prefix: prefix,
		data:   ToolResult(data),
	}

	return &r
}

type readResultResponse struct {
	err    bool
	prefix string
	data   []toolResult
}

type toolResult struct {
	resultNum int  // -> in range [0, 32767]
	ok        bool // -> image pass/fail
}

func (r *readResultResponse) Ok() bool {
	return !r.err
}

func ToolResult(s string) []toolResult {
	tools := make([]toolResult, int((strings.Count(s, ",")+1)/2))
	data := strings.Split(s, ",")
	l := len(data)
	for i := 0; i < l; i += 2 {
		count, _ := strconv.Atoi(data[i])
		ok := false
		if data[i+1] == "OK" {
			ok = true
		}
		tools = append(tools[1:], toolResult{
			resultNum: count,
			ok:        ok,
		})
	}

	return tools
}

type trig struct{}

// Command camera to take image and return result
func Trig() *trig {
	return &trig{}
}

func (cmd *trig) Compose() string {
	return "T2"
}

func (cmd *trig) Interpret(s string) *trigResponse {
	err := false
	if s[:strings.Index(s, ",")] == "ER" { // err msg returned
		err = true
		_, s, _ = strings.Cut(s, ",")
	}
	prefix, data, _ := strings.Cut(s, ",")
	t := trigResponse{
		err:    err,
		prefix: prefix,
		data:   ToolResult(data),
	}

	return &t
}

type trigResponse struct {
	err    bool
	prefix string
	data   []toolResult
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

func (cmd *programRead) Interpret(s string) *programReadResponse {
	data := strings.Split(s, ",")
	num, _ := strconv.Atoi(data[1])
	r := programReadResponse{
		prefix:        data[0],
		programNumber: num,
	}

	return &r
}

type programReadResponse struct {
	prefix        string
	programNumber int
}

func (r programReadResponse) Ok() bool {
	ok := true
	if r.prefix != "PR" {
		ok = false
	}
	return ok
}

type programWrite struct {
	num int
}

// Set camera program, program number must be in range [0, 127]
func ProgramWrite(num int) (*programWrite, error) {
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

	return &pw, err
}

func (cmd *programWrite) Compose() string {
	return fmt.Sprintf("PW,%v", cmd.num)
}

func (cmd *programWrite) Interpret(s string) *programWriteResponse {
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
	toolNum    int
	upperLimit bool
}

// Read a program limit of the specified tool and limit type (upper/lower)
func ThresholdRead(toolNum int, upperLimit bool) (*thresholdRead, error) {
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

	return &tr, err
}

func (cmd *thresholdRead) Compose() string {
	b := "0"
	if cmd.upperLimit {
		b = "1"
	}

	return fmt.Sprintf("DR,%v,%v", cmd.toolNum, b)
}

func (cmd *thresholdRead) Interpret(s string) *thresholdReadResponse {
	data := strings.Split(s, ",")
	prefix := data[0]
	toolNum, _ := strconv.Atoi(data[1])
	upper := false
	if data[2] == "1" {
		upper = true
	}
	tLimit, _ := strconv.Atoi(data[3])

	r := thresholdReadResponse{
		prefix:         prefix,
		toolNum:        toolNum,
		upperLimit:     upper,
		thresholdLimit: tLimit,
	}

	return &r
}

type thresholdReadResponse struct { //docs says it wont error, can always add later if needed
	prefix         string
	toolNum        int
	upperLimit     bool
	thresholdLimit int
}

func (r *thresholdReadResponse) Ok() bool {
	ok := false
	if r.prefix == "DR" {
		ok = true
	}

	return ok
}

type thresholdWrite struct {
	toolNum    int
	upperLimit bool
	newLimit   int
}

// Set the threshold of the given tool and limit to entered limit
func ThresholdWrite(toolNum int, upperLimit bool, newLimit int) (*thresholdWrite, error) {
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

	return &tw, err
}

func (cmd *thresholdWrite) Compose() string {
	b := "0"
	if cmd.upperLimit {
		b = "1"
	}
	return fmt.Sprintf("DW,%v,%v,%v", cmd.toolNum, b, cmd.newLimit)
}

func (cmd *thresholdWrite) Interpret(s string) *thresholdWriteResponse {
	data := strings.Split(s, ",")
	prefix := data[0]
	toolNum, _ := strconv.Atoi(data[1])

	r := thresholdWriteResponse{
		prefix:  prefix,
		toolNum: toolNum,
	}

	return &r
}

type thresholdWriteResponse struct {
	prefix  string
	toolNum int
}

func (r *thresholdWriteResponse) Ok() bool {
	ok := false
	if r.prefix == "DW" {
		ok = true
	}

	return ok
}

type textRead struct {
	toolNum int
}

// Read current master text value for specified tool in range [1, 64]
func TextRead(toolNum int) (*textRead, error) {
	num := -1
	var err error = nil
	if 0 < toolNum && toolNum < 65 {
		num = toolNum
	} else {
		err = fmt.Errorf("invalid tool number %v, should be in range [1, 64]", toolNum)
	}
	return &textRead{num}, err
}

func (cmd *textRead) Compose() string {
	return fmt.Sprintf("CR,%v", cmd.toolNum)
}

func (cmd *textRead) Interpret(s string) *textReadResponse {
	data := strings.Split(s, ",")
	prefix := data[0]
	toolNum, _ := strconv.Atoi(data[1])
	masterText := strings.Split(data[2], "")
	r := textReadResponse{
		prefix:     prefix,
		toolNum:    toolNum,
		masterText: masterText,
	}

	return &r
}

type textReadResponse struct {
	prefix     string
	toolNum    int
	masterText []string
}

func (r *textReadResponse) Ok() bool {
	ok := false
	if r.prefix == "CR" {
		ok = true
	}

	return ok
}

type textWrite struct {
	toolNum    int
	masterText []string
}

const MAX_TEXT_LEN = 16

// Set master text for specified tool number, text must be less than 16 characters
func TextWrite(toolNum int, text string) (*textWrite, error) {
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

	return &tw, err
}

func (cmd *textWrite) Compose() string {
	return fmt.Sprintf("CW,%v,%v", cmd.toolNum, strings.Join(cmd.masterText, ""))
}

func (cmd *textWrite) Interpret(s string) *textWriteResponse {
	data := strings.Split(s, ",")
	prefix := data[0]
	toolNum, _ := strconv.Atoi(data[1])

	r := textWriteResponse{
		prefix:  prefix,
		toolNum: toolNum,
	}

	return &r
}

type textWriteResponse struct {
	prefix  string
	toolNum int
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
