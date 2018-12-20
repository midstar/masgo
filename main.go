package main

import (
	"fmt"

	"github.com/midstar/masgo/tellstick"
)

func main() {
	tl, err := tellstick.NewTellstickLibrary()
	if err != nil {
		panic(fmt.Sprintf("Error: %s\n", err))
	}
	ids, _ := tl.GetDeviceIds()
	for _, id := range ids {
		fmt.Printf("Id %d: Name: '%s' OnOff: %t Dim: %t Learn: %t\n", id, tl.GetName(id), tl.SupportsOnOff(id), tl.SupportsDim(id), tl.SupportsLearn(id))
	}
}
