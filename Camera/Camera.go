package Camera

import "fmt"

type Camera struct {
	ip        []int
	port      int
	delimiter rune
}

// init iv3 camera with default values, change thru setters
func New() *Camera {
	c := Camera{
		ip:        []int{192, 168, 1, 1},
		port:      8500,
		delimiter: '\r',
	}

	return &c
}

// getters
func (c *Camera) getIp() []int {
	return c.ip
}

func (c *Camera) getPort() int {
	return c.port
}

func (c *Camera) getDelimiter() rune {
	return c.delimiter
}

// setters
func (c *Camera) SetIp(vals []int) bool {
	//validate format
	ok := true
	if len(vals) != 4 {
		ok = false
	}
	if ok {
		index, value := 0, 0
		for i, val := range vals {
			if val < 0 || val > 255 {
				index = i
				value = val
				ok = false
				break
			}
		}
		if ok {
			//input should have 4 valid values
			c.ip = vals
		} else {
			fmt.Printf("Invalid number %v at position %v\n", value, index)
		}
	} else {
		fmt.Printf("Invalid number of values, please try again. There should be 4 values between [0, 255] - You entered: %v\n", vals)
	}

	return ok
}

func (c *Camera) SetPort(newPort int) bool {
	ok := true
	if newPort > 1023 && newPort < 65536 {
		c.port = newPort
	} else {
		ok = false
		fmt.Printf("Invalid value for new port number, value must be in range [1024, 65535] - you entered %v\n", newPort)
	}

	return ok
}

func (c *Camera) SetDelimiter(delim rune) bool {
	c.delimiter = delim
	return true
}
