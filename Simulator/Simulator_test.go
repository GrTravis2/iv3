package Simulator

import (
	"fmt"
	"iv3/Camera"
	"iv3/Messenger"
	"iv3/Operate"
	"testing"
)

func TestSuiteSimulator(t *testing.T) {
	msgr := Messenger.NewMessenger("test", Camera.NewCamera([]int{}))
	sim := NewSimulator(msgr.Cameras["test"], true)
	go sim.Run()

	response, err := msgr.Send("test", Operate.BlindTrig())
	if err != nil {
		t.Errorf("error sending message: %v", err)
	}
	fmt.Printf("response: %v", response)
}
