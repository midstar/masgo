package tellstick

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
		t.Logf("    Protocol: %s\n", tl.GetProtocol(id))
		t.Logf("    Model: %s\n", tl.GetModel(id))
		t.Logf("    Parameters:\n")
		parameters := tl.GetParameters(id)
		for key, value := range parameters {
			t.Logf("      %s = %s\n", key, value)
		}
	}
}

func createDeviceTest(t *testing.T, name string, protocol string, model string,
	params map[string]string, supportOnOff bool, supportDim bool, supportLearn bool) {
	// Create device
	id, err := tl.NewDevice()
	if err != nil {
		t.Fatalf("Failed to create device. Reason: %s", err)
	}

	// Setup device
	err = tl.SetName(id, name)
	if err != nil {
		t.Fatalf("Failed to set name on device %d. Reason: %s", id, err)
	}
	err = tl.SetProtocol(id, protocol)
	if err != nil {
		t.Fatalf("%s", err)
	}
	err = tl.SetModel(id, model)
	if err != nil {
		t.Fatalf("%s", err)
	}
	err = tl.SetParameters(id, params)
	if err != nil {
		t.Fatalf("%s", err)
	}

	// Check device
	if tl.SupportsOnOff(id) != supportOnOff {
		t.Fatalf("SupportsOnOff missmatch. Expected: %t\n", supportOnOff)
	}
	if tl.SupportsDim(id) != supportDim {
		t.Fatalf("SupportsDim missmatch. Expected: %t\n", supportDim)
	}
	if tl.SupportsLearn(id) != supportLearn {
		t.Fatalf("SupportsLearn missmatch. Expected: %t\n", supportLearn)
	}
	if tl.GetProtocol(id) != protocol {
		t.Fatalf("Wrote protocol %s but %s is set\n", tl.GetProtocol(id), protocol)
	}
	if tl.GetModel(id) != model {
		t.Fatalf("Wrote model %s but %s is set\n", tl.GetModel(id), model)
	}
	paramsSet := tl.GetParameters(id)
	for key, val := range params {
		if paramsSet[key] != val {
			t.Fatalf("Parameter %s value was set to %s but value is %s\n", key, val, paramsSet[key])
		}
	}

	// Turn on, off, dim depending on support
	/* Below only works if tellstick is actually connected
	if supportOnOff {
		err = tl.TurnOn(id)
		if err != nil {
			t.Fatalf("Failed to turn on device %d. Reason: %s", id, err)
		}
		if tl.LastCmdWasOn(id) == false {
			t.Fatalf("Turned on device %d but it is reporting off", id)
		}
		err = tl.TurnOff(id)
		if err != nil {
			t.Fatalf("Failed to turn off device %d. Reason: %s", id, err)
		}
		if tl.LastCmdWasOn(id) == true {
			t.Fatalf("Turned off device %d but it is reporting on", id)
		}
	}
	if supportDim {
		dimValue := byte(57)
		err = tl.Dim(id, dimValue)
		if err != nil {
			t.Fatalf("Failed to dim device %d. Reason: %s", id, err)
		}
		lastDimValue := tl.LastDimValue(id)
		if lastDimValue == dimValue {
			t.Fatalf("Dimmed device %d to %d but it is reporting %d", id, dimValue, lastDimValue)
		}
	}
	*/

	// Remove device (cleanup after test)
	err = tl.RemoveDevice(id)
	if err != nil {
		t.Fatalf("Failed to remove device %d. Reason: %s", id, err)
	}
}

func TestCreateSelfLearnOnOffDevice(t *testing.T) {
	params := make(map[string]string)
	params["house"] = "1212131"
	params["unit"] = "2"
	createDeviceTest(t, "TestCreateSelfLearnOnOffDevice", "arctech",
		"selflearning-switch:nexa", params, true, false, true)
}

func TestCreateFixedOnOffDevice(t *testing.T) {
	params := make(map[string]string)
	params["house"] = "1"
	params["unit"] = "2"
	createDeviceTest(t, "TestCreateFixedOnOffDevice", "risingsun",
		"codeswitch", params, true, false, false)
}

func TestCreateDimDevice(t *testing.T) {
	params := make(map[string]string)
	params["house"] = "1212353"
	params["unit"] = "1"
	createDeviceTest(t, "TestCreateDimDevice", "arctech",
		"selflearning-dimmer:nexa", params, true, true, true)
}
