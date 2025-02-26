package Camera

import (
	"fmt"
	"strconv"
	"strings"
)

type Camera struct {
	ip        []int
	port      int
	delimiter rune
}

// init iv3 camera with default values, change thru setters
func NewCamera(ipAddress []int) *Camera {
	c := Camera{
		ip:        []int{0, 0, 0, 0},
		port:      8500,
		delimiter: '\r',
	}
	if c.SetIp(ipAddress) != nil {
		fmt.Printf("Invalid ipAddress %v, camera initialized with default ip %v\nip can be updated later.", ipAddress, c.ip)
	}

	return &c
}

// **getters**
func (c *Camera) GetIp() []int {
	return c.ip
}

func (c *Camera) GetPort() int {
	return c.port
}

func (c *Camera) GetDelimiter() rune {
	return c.delimiter
}

// needs test!
func (c *Camera) GetAddress() string {
	s := make([]string, len(c.ip))
	for i, num := range c.ip {
		s[i] = strconv.Itoa(num)
	}
	return strings.Join(s, ".") + ":" + strconv.Itoa(c.GetPort())
}

// **setters**

func (c *Camera) SetIp(vals []int) error {
	//validate format
	var err error = nil
	if len(vals) != 4 {
		err = fmt.Errorf("expected 4 value ip address, found %v", vals)
	} else {
		for _, val := range vals {
			if val < 0 || val > 255 {
				err = fmt.Errorf("expected ip value in range [0, 255], found %v", val)
				break
			}
		}
		if err == nil {
			//input should have 4 valid values
			c.ip = vals
		}
	}

	return err
}

// set camera's port value, must be in range [1024, 65535]
func (c *Camera) SetPort(newPort int) error {
	var err error = nil
	if newPort > 1023 && newPort < 65536 {
		c.port = newPort
	} else {
		err = fmt.Errorf("expected port in range [1024, 65535], found %v", newPort)
	}

	return err
}

func (c *Camera) SetDelimiter(delim rune) {
	c.delimiter = delim
}

// Camera attributes to be used in command packages

type toolResult struct {
	resultNum int  // -> in range [0, 32767]
	ok        bool // -> image pass/fail
}

type ToolResults []toolResult //

// tool result data for the selected program, each tool
// result contains result count and pass/fail for each tool
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

type ProgramNumber int

// Camera program number, must be in range [0, 127]
func MakeProgramNumber(n int) (ProgramNumber, error) {
	var err error = nil
	if n < 0 || 127 < n {
		err = fmt.Errorf("input n outside of allowable range [0, 127], found %v", n)
		n = -1
	}

	return ProgramNumber(n), err
}

type ToolNumber int

// camera internal program number in range [0, 64] where
// 0 is the position adjustment, and [1, 64] is generic tool number
func MakeToolNumber(n int) (ToolNumber, error) {
	var err error = nil

	if n < 0 || 64 < n {
		err = fmt.Errorf("input n outside of allowable range [0, 64], found %v", n)
		n = -1
	}

	return ToolNumber(n), err
}

type UpperLimit bool

// camera program pass/fail threshold targeting
// 0 if modifying lower limit, 1 for upper limt
func MakeUpperLimit(b bool) UpperLimit {
	return UpperLimit(b)
}

type Threshold int

// camera program pass/fail threshold in range [0, 9999999]
func MakeThreshold(n int) (Threshold, error) {
	var err error = nil
	if n < 0 || 9999999 < n {
		err = fmt.Errorf("input n outside of allowable range [0, 9999999], found %v", n)
		n = -1
	}

	return Threshold(n), err
}

type MasterText string

const MAX_TEXT_LEN = 16

// 16 char string with trailing spaces if text is shorter than 16 chars
func MakeMasterText(s string) MasterText {
	text := make([]rune, MAX_TEXT_LEN) // text will always be 16 chars
	for i := range text {
		if i < len(s) {
			text[i] = rune(s[i])
		} else {
			text[i] = ' '
		}
	}

	mText := string(text)

	return MasterText(mText)
}

type CharsRequired int

// specify's min/max number of characters the program scans for during text judgement
func MakeCharsRequired(n int) (CharsRequired, error) {
	var err error = nil
	if n < 1 || 16 < n {
		err = fmt.Errorf("input n outside of allowable range [1, 16], found %v", n)
		n = -1
	}

	return CharsRequired(n), err
}
