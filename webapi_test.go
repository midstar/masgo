package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

var mock *DeviceMockLibrary
var baseURL string

const port int = 9843

func TestMain(m *testing.M) {
	mock = NewDeviceMockLibrary()
	webAPI := CreateWebAPI(port, mock)
	baseURL = fmt.Sprintf("http://localhost:%d", port)

	// Add some devices
	mock.devices[1] = &DeviceMock{
		id:       1,
		name:     "onename",
		model:    "onemodel",
		protocol: "oneprotocol",
		parameters: map[string]string{
			"house": "onehouse",
			"unit":  "oneunit"},
		supportOnOff: true,
		supportDim:   false,
		supportLearn: false,
		isOn:         false,
		dimLevel:     0}

	mock.devices[2] = &DeviceMock{
		id:       2,
		name:     "twoname",
		model:    "twomodel",
		protocol: "twoprotocol",
		parameters: map[string]string{
			"house": "twohouse",
			"unit":  "twounit"},
		supportOnOff: true,
		supportDim:   true,
		supportLearn: false,
		isOn:         false,
		dimLevel:     0}

	mock.devices[3] = &DeviceMock{
		id:       3,
		name:     "threename",
		model:    "threemodel",
		protocol: "threeprotocol",
		parameters: map[string]string{
			"house": "threehouse",
			"unit":  "threeunit"},
		supportOnOff: true,
		supportDim:   true,
		supportLearn: true,
		isOn:         false,
		dimLevel:     0}

	webAPI.Start()
	retCode := m.Run()
	webAPI.Stop()
	os.Exit(retCode)
}

func getObject(t *testing.T, path string, v interface{}) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", baseURL, path))
	if err != nil {
		t.Fatalf("Unable to get path %s. Reason: %s", path, err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code for path %s: %d", path, resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unable to read body for %s. Reason: %s", path, err)
	}
	err = json.Unmarshal(body, &v)
	if err != nil {
		t.Fatalf("Unable decode path %s. Reason: %s", path, err)
	}
}

func post(t *testing.T, path string) {
	resp, err := http.Post(fmt.Sprintf("%s/%s", baseURL, path), "", nil)
	if err != nil {
		t.Fatalf("Unable to post path %s. Reason: %s", path, err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code for post path %s: %d", path, resp.StatusCode)
	}
	defer resp.Body.Close()
}

func TestGetDeviceIds(t *testing.T) {
	var deviceIDs []int
	getObject(t, "devices", &deviceIDs)
	assertEqualsInt(t, "Invalid devices received", len(mock.devices), len(deviceIDs))
	for id := range mock.devices {
		idFound := false
		for i := 0; i < len(deviceIDs); i++ {
			if id == deviceIDs[i] {
				idFound = true
				break
			}
		}
		if !idFound {
			t.Fatalf("Device with id %d not recived", id)
		}
	}
}

func TestTurnOnOff(t *testing.T) {
	mock.devices[2].isOn = false
	post(t, "devices/2/on")
	assertTrue(t, "Device 2 shall be on", mock.devices[2].isOn == true)
	post(t, "devices/2/off")
	assertTrue(t, "Device 2 shall be off", mock.devices[2].isOn == false)
}
