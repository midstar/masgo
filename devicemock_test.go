package main

import (
	"fmt"

	"github.com/midstar/masgo/tellstick"
)

type DeviceMock struct {
	id           int
	name         string
	model        string
	protocol     string
	parameters   map[string]string
	supportOnOff bool
	supportDim   bool
	supportLearn bool
	isOn         bool
	dimLevel     byte
	learnCount   int
}

type DeviceMockLibrary struct {
	devices map[int]*DeviceMock // Keyed on id
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func NewDeviceMockLibrary() *DeviceMockLibrary {
	return &DeviceMockLibrary{
		devices: make(map[int]*DeviceMock)}
}

// GetDeviceIds will always return a valid slice even on errors
func (tl *DeviceMockLibrary) GetDeviceIds() ([]int, error) {
	var ids []int
	for key := range tl.devices {
		ids = append(ids, key)
	}
	return ids, nil
}

func (tl *DeviceMockLibrary) GetName(id int) string {
	if val, ok := tl.devices[id]; ok {
		return val.name
	}
	return ""
}

func (tl *DeviceMockLibrary) SetName(id int, name string) error {
	if val, ok := tl.devices[id]; ok {
		val.name = name
		return nil
	}
	return fmt.Errorf("no device exist with id %d", id)
}

func (tl *DeviceMockLibrary) SupportsOnOff(id int) bool {
	if val, ok := tl.devices[id]; ok {
		return val.supportOnOff
	}
	return false
}

func (tl *DeviceMockLibrary) SupportsDim(id int) bool {
	if val, ok := tl.devices[id]; ok {
		return val.supportDim
	}
	return false
}

func (tl *DeviceMockLibrary) SupportsLearn(id int) bool {
	if val, ok := tl.devices[id]; ok {
		return val.supportLearn
	}
	return false
}

func (tl *DeviceMockLibrary) NewDevice() (int, error) {
	// Find next id
	ok := true
	id := 0
	for ok == true {
		id++
		_, ok = tl.devices[id]
	}
	tl.devices[id] = &DeviceMock{
		parameters: make(map[string]string)}
	return id, nil
}

func (tl *DeviceMockLibrary) RemoveDevice(id int) error {
	delete(tl.devices, id)
	return nil
}

func (tl *DeviceMockLibrary) GetProtocol(id int) string {
	if val, ok := tl.devices[id]; ok {
		return val.protocol
	}
	return ""
}

func (tl *DeviceMockLibrary) SetProtocol(id int, protocol string) error {
	if val, ok := tl.devices[id]; ok {
		val.protocol = protocol
		return nil
	}
	return fmt.Errorf("no device exist with id %d", id)
}

func (tl *DeviceMockLibrary) GetModel(id int) string {
	if val, ok := tl.devices[id]; ok {
		return val.model
	}
	return ""
}

func (tl *DeviceMockLibrary) SetModel(id int, model string) error {
	if val, ok := tl.devices[id]; ok {
		val.model = model
		return nil
	}
	return fmt.Errorf("no device exist with id %d", id)
}

func (tl *DeviceMockLibrary) GetParameters(id int) map[string]string {
	if val, ok := tl.devices[id]; ok {
		return val.parameters
	}
	return make(map[string]string) // Empty map
}

func (tl *DeviceMockLibrary) SetParameters(id int, paramAndValues map[string]string) error {
	var val *DeviceMock
	var ok bool
	if val, ok = tl.devices[id]; ok == false {
		return fmt.Errorf("no device exist with id %d", id)
	}
	for parameter, value := range paramAndValues {
		if stringInSlice(parameter, tellstick.Parameters) == false {
			return fmt.Errorf("unknown parameter '%s'", parameter)
		}
		val.parameters[parameter] = value
	}
	return nil
}

func (tl *DeviceMockLibrary) TurnOn(id int) error {
	if val, ok := tl.devices[id]; ok {
		if val.supportOnOff == false {
			return fmt.Errorf("on off not supported for device with id %d", id)
		}
		val.isOn = true
		return nil
	}
	return fmt.Errorf("no device exist with id %d", id)
}

func (tl *DeviceMockLibrary) TurnOff(id int) error {
	if val, ok := tl.devices[id]; ok {
		if val.supportOnOff == false {
			return fmt.Errorf("on off not supported for device with id %d", id)
		}
		val.isOn = false
		return nil
	}
	return fmt.Errorf("no device exist with id %d", id)
}

func (tl *DeviceMockLibrary) Dim(id int, level byte) error {
	if val, ok := tl.devices[id]; ok {
		if val.supportDim == false {
			return fmt.Errorf("dim not supported for device with id %d", id)
		}
		val.dimLevel = level
		return nil
	}
	return fmt.Errorf("no device exist with id %d", id)
}

func (tl *DeviceMockLibrary) Learn(id int) error {
	if val, ok := tl.devices[id]; ok {
		if val.supportLearn == false {
			return fmt.Errorf("learn not supported for device with id %d", id)
		}
		val.learnCount++
		return nil
	}
	return fmt.Errorf("no device exist with id %d", id)
}

func (tl *DeviceMockLibrary) LastCmdWasOn(id int) bool {
	if val, ok := tl.devices[id]; ok {
		return val.isOn
	}
	return false
}

func (tl *DeviceMockLibrary) LastDimValue(id int) byte {
	if val, ok := tl.devices[id]; ok {
		return val.dimLevel
	}
	return 0
}
