package main

import (
	"fmt"
	"os"
	"testing"
)

var tl *TellstickLibrary

func TestMain(m *testing.M) {
	var err error
	tl, err = NewTellstickLibrary()
	if err != nil {
		panic(fmt.Sprintf("Tellstick library not installed. Unable to run tests. Reason: %s", err))
	}
	retCode := m.Run()
	os.Exit(retCode)
}

func TestGetDevices(t *testing.T) {
	ids, err := tl.GetDeviceIds()
	if err != nil {
		t.Fatalf("Failed to get devices. Reason: %s", err)
	}
	for _, id := range ids {
		t.Logf("Id %d: Name: '%s' OnOff: %t Dim: %t Learn: %t\n",
			id,
			tl.GetName(id),
			tl.SupportsOnOff(id),
			tl.SupportsDim(id),
			tl.SupportsLearn(id))
		t.Logf("    Parameters:\n")
		parameters := tl.GetParameters(id)
		for key, value := range parameters {
			t.Logf("      %s = %s\n", key, value)
		}
	}
}

func TestCreateDevice(t *testing.T) {
	id, err := tl.NewDevice()
	if err != nil {
		t.Fatalf("Failed to create device. Reason: %s", err)
	}
	err = tl.SetName(id, "TestCreateDevice")
	if err != nil {
		t.Fatalf("Failed to set name on device %d. Reason: %s", id, err)
	}
}
