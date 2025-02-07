package Simulator

import (
	"testing"

	"iv3/Camera"
	"iv3/Messenger"
	"iv3/Operate"
)

func TestSuiteSimulator(t *testing.T) {
	msgr := Messenger.NewMessenger("test", Camera.NewCamera([]int{}))
	sim := NewSimulator(msgr.Cameras["test"], true)
	go sim.Run()

	_, err := msgr.Send("test", Operate.BlindTrig())
	if err != nil {
		t.Errorf("error sending message: %v", err)
	}

	_, err = msgr.Send("test", Operate.ReadResult())
	if err != nil {
		t.Errorf("error sending message: %v", err)
	}

	_, err = msgr.Send("test", Operate.Trig())
	if err != nil {
		t.Errorf("error sending message: %v", err)
	}

	_, err = msgr.Send("test", Operate.ProgramRead())
	if err != nil {
		t.Errorf("error sending message: %v", err)
	}

	// Beginging of messages w/ args, require extra cmd var
	var cmd Messenger.Message = nil

	cmd, _ = Operate.ProgramWrite(0)
	_, err = msgr.Send("test", cmd)
	if err != nil {
		t.Errorf("error sending message: %v", err)
	}

	cmd, _ = Operate.ThresholdRead(0, false)
	_, err = msgr.Send("test", cmd)
	if err != nil {
		t.Errorf("error sending message: %v", err)
	}

	cmd, _ = Operate.ThresholdWrite(0, true, 0)
	_, err = msgr.Send("test", cmd)
	if err != nil {
		t.Errorf("error sending message: %v", err)
	}

	cmd, _ = Operate.TextRead(1)
	_, err = msgr.Send("test", cmd)
	if err != nil {
		t.Errorf("error sending message: %v", err)
	}

	cmd, _ = Operate.TextWrite(1, " :) ")
	_, err = msgr.Send("test", cmd)
	if err != nil {
		t.Errorf("error sending message: %v", err)
	}

}
