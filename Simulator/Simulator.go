package Simulator

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"iv3/Camera"
)

var localhostAddr = [...]int{127, 0, 0, 1}

// str representation of command messges and responses
// in map format {prefix : response}
const responseData string = `
T1,T1
RT,RT,00000,NG
T2,RT,00000,NG
PR,PR,099
PW,PW
DR,DR,03,0,0023250
DW,DW,01,0,80
CR,CR,01,123456789ABCDE  
CW,CW,0123456789
CNR,CNR,01,06,10
CNW,CNW,01,06,10`

type simulator struct {
	camera    *Camera.Camera
	responses map[string]string
}

func NewSimulator(c *Camera.Camera, localhost bool) *simulator {
	lines := strings.Split(strings.Trim(responseData, "\n"), "\n")
	responses := make(map[string]string)
	for _, line := range lines {
		prefix, response, _ := strings.Cut(line, ",")
		responses[prefix] = response
	}
	c.SetPort(3333)
	sim := simulator{
		camera:    c,
		responses: responses,
	}
	if localhost {
		sim.camera.SetIp(localhostAddr[:])
	}

	return &sim
}

// Run ...
func (sim *simulator) Run() {
	listener, err := net.Listen("tcp", sim.camera.GetAddress())
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go sim.handleRequest(conn)
	}
}

func (sim *simulator) handleRequest(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString(byte(sim.camera.GetDelimiter()))
		if err != nil {
			conn.Close()
			return
		}
		msg, _ = strings.CutSuffix(msg, "\r")
		fields := strings.Split(msg, ",")
		prefix := ""
		if fields[0] == "ER" {
			prefix = fields[1]
		} else {
			prefix = fields[0]
		}
		fmt.Printf("received message %s\n", msg)
		conn.Write([]byte(sim.responses[prefix] + string(sim.camera.GetDelimiter())))
	}
}
