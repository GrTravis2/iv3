package Messenger

import (
	"fmt"
	"net"
	"strings"

	"github.com/GrTravis2/iv3/Camera"
)

type Message interface {
	Compose() string
	Interpret(string) Response
}

type Response interface {
	Ok() bool
}

type Messenger struct {
	Cameras map[string]*Camera.Camera
}

func NewMessenger(name string, c *Camera.Camera) *Messenger {
	m := Messenger{
		Cameras: make(map[string]*Camera.Camera),
	}
	m.Add(name, c)

	return &m
}

func (m *Messenger) Send(name string, msg Message) (Response, error) {
	c := m.Cameras[name]
	var result string = ""
	conn, err := net.Dial("tcp", c.GetAddress())
	if err == nil {
		data := []byte(msg.Compose() + string(c.GetDelimiter()))
		_, err := conn.Write(data)
		if err == nil {
			//read response
			buffer := make([]byte, 1024)

			//read data from client
			n, err := conn.Read(buffer)
			if err == nil {
				result = strings.Trim(string(buffer[:n]), fmt.Sprintf("%v", c.GetDelimiter()))
			}
		}
	}
	conn.Close()

	return msg.Interpret(result), err
}

func (m *Messenger) Add(name string, c *Camera.Camera) {
	m.Cameras[name] = c
}
