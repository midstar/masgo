package main

import (
	"fmt"

	"github.com/midstar/masgo/tellstick"
)

type DeviceLibrary interface {
	GetDeviceIds() ([]int, error)
	GetName(id int) string
	SetName(id int, name string) error
	SupportsOnOff(id int) bool
	SupportsDim(id int) bool
	SupportsLearn(id int) bool
	NewDevice() (int, error)
	RemoveDevice(id int) error
	GetProtocol(id int) string
	SetProtocol(id int, protocol string) error
	GetModel(id int) string
	SetModel(id int, model string) error
	GetParameters(id int) map[string]string
	SetParameters(id int, paramAndValues map[string]string) error
	TurnOn(id int) error
	TurnOff(id int) error
	Dim(id int, level byte) error
	Learn(id int) error
	LastCmdWasOn(id int) bool
	LastDimValue(id int) byte
	MinDimLevel() int
	MaxDimLevel() int
}

func main() {
	var dl DeviceLibrary
	dl, err := tellstick.NewTellstickLibrary()
	if err != nil {
		panic(fmt.Sprintf("Error: %s\n", err))
	}
	groups := createGroups(dl)

	webAPI := CreateWebAPI(9834, dl, groups)
	httpServerDone := webAPI.Start()
	<-httpServerDone // Block until http server is done
	/*
		ids, _ := dl.GetDeviceIds()
		for _, id := range ids {
			fmt.Printf("Id %d: Name: '%s' OnOff: %t Dim: %t Learn: %t\n", id,
				dl.GetName(id), dl.SupportsOnOff(id), dl.SupportsDim(id), dl.SupportsLearn(id))
		}
	*/
}
