package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
		t.Fatalf("Unexpected status code for path %s: %d (%s)",
			path, resp.StatusCode, respToString(resp.Body))
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
		t.Fatalf("Unexpected status code for post path %s: %d (%s)",
			path, resp.StatusCode, respToString(resp.Body))
	}
	defer resp.Body.Close()
}

func postObject(t *testing.T, path string, v interface{}) {
	bodyBytes, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Unable to generate json for path %s. Reason: %s", path, err)
	}
	body := bytes.NewReader(bodyBytes)
	resp, err := http.Post(fmt.Sprintf("%s/%s", baseURL, path), "application/json", body)
	if err != nil {
		t.Fatalf("Unable to post path %s. Reason: %s", path, err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code for post path %s: %d (%s)",
			path, resp.StatusCode, respToString(resp.Body))
	}
	defer resp.Body.Close()
}

func putObject(t *testing.T, path string, v interface{}) {
	bodyBytes, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Unable to generate json for path %s. Reason: %s", path, err)
	}
	body := bytes.NewReader(bodyBytes)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", baseURL, path), body)
	if err != nil {
		t.Fatalf("Unable to create request for put path %s. Reason: %s", path, err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Unable to execute request for put path %s. Reason: %s", path, err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code for put path %s: %d (%s)",
			path, resp.StatusCode, respToString(resp.Body))
	}
	defer resp.Body.Close()
}

func respToString(response io.ReadCloser) string {
	defer response.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response)
	return buf.String()
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

func TestGetDeviceConfig(t *testing.T) {
	var config DeviceConfig
	getObject(t, "devices/2/config", &config)
	assertEqualsInt(t, "Invalid id received", 2, config.ID)
	assertEqualsStr(t, "Invalid name received", mock.devices[2].name, config.Name)
	assertEqualsStr(t, "Invalid protocol received", mock.devices[2].protocol, config.Protocol)
	assertEqualsStr(t, "Invalid model received", mock.devices[2].model, config.Model)
	assertEqualsStr(t, "Invalid parameter house received",
		mock.devices[2].parameters["house"], config.Parameters["house"])
	assertEqualsStr(t, "Invalid parameter unit received",
		mock.devices[2].parameters["unit"], config.Parameters["unit"])
}

func TestGetDeviceConfigs(t *testing.T) {
	var config []DeviceConfig
	getObject(t, "devices/config", &config)
	assertEqualsInt(t, "Invalid number or entries", len(mock.devices), len(config))
}

func TestPutDeviceConfig(t *testing.T) {
	var oldConfig DeviceConfig
	getObject(t, "devices/2/config", &oldConfig)
	newConfig := DeviceConfig{
		Name:     "newname",
		Protocol: "newprotocol",
		Model:    "newmodel",
		Parameters: map[string]string{
			"house": "newhouse",
			"unit":  "newunit"},
	}
	putObject(t, "devices/2/config", &newConfig)
	assertEqualsStr(t, "Invalid name", newConfig.Name, mock.devices[2].name)
	assertEqualsStr(t, "Invalid protocol", newConfig.Protocol, mock.devices[2].protocol)
	assertEqualsStr(t, "Invalid model", newConfig.Model, mock.devices[2].model)
	assertEqualsStr(t, "Invalid parameter house",
		newConfig.Parameters["house"], mock.devices[2].parameters["house"])
	assertEqualsStr(t, "Invalid parameter unit received",
		newConfig.Parameters["unit"], mock.devices[2].parameters["unit"])

	// Restore old config
	putObject(t, "devices/2/config", &oldConfig)
}

func TestDeleteDevice(t *testing.T) {
	// Copy device
	oldCopy := *mock.devices[2]

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/devices/2/config", baseURL), nil)
	if err != nil {
		t.Fatalf("Unable to create request for delete devices/2/config. Reason: %s", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Unable to execute request for delete devices/2/config. Reason: %s", err)
	}
	assertEqualsInt(t, "Unexpected status code",
		http.StatusOK, resp.StatusCode)

	_, exists := mock.devices[2]
	assertTrue(t, "Device 2 was never removed", exists == false)

	// Put back old device
	mock.devices[2] = &oldCopy

}

func TestTurnOnOff(t *testing.T) {
	mock.devices[2].isOn = false
	post(t, "devices/2/on")
	assertTrue(t, "Device 2 shall be on", mock.devices[2].isOn == true)
	post(t, "devices/2/off")
	assertTrue(t, "Device 2 shall be off", mock.devices[2].isOn == false)

	// Test for device not supporing on / off
	mock.devices[3].supportOnOff = false
	resp, _ := http.Post(fmt.Sprintf("%s/devices/3/on", baseURL), "", nil)
	assertEqualsInt(t, "Unexpected status code",
		http.StatusMethodNotAllowed, resp.StatusCode)
	mock.devices[3].supportOnOff = true
}

func TestDim(t *testing.T) {
	mock.devices[3].dimLevel = 0
	post(t, "devices/3/dim/1")
	assertEqualsInt(t, "Unexpected dim level", 1, int(mock.devices[3].dimLevel))
	post(t, "devices/3/dim/255")
	assertEqualsInt(t, "Unexpected dim level", 255, int(mock.devices[3].dimLevel))
	post(t, "devices/3/dim/0")
	assertEqualsInt(t, "Unexpected dim level", 0, int(mock.devices[3].dimLevel))

	// Test invalid values
	resp, _ := http.Post(fmt.Sprintf("%s/devices/3/dim/-1", baseURL), "", nil)
	assertEqualsInt(t, "Unexpected status code",
		http.StatusBadRequest, resp.StatusCode)
	resp, _ = http.Post(fmt.Sprintf("%s/devices/3/dim/256", baseURL), "", nil)
	assertEqualsInt(t, "Unexpected status code",
		http.StatusBadRequest, resp.StatusCode)
	resp, _ = http.Post(fmt.Sprintf("%s/devices/3/dim/abc", baseURL), "", nil)
	assertEqualsInt(t, "Unexpected status code",
		http.StatusBadRequest, resp.StatusCode)

	// Test for device not supporing dim
	resp, _ = http.Post(fmt.Sprintf("%s/devices/1/dim/1", baseURL), "", nil)
	assertEqualsInt(t, "Unexpected status code",
		http.StatusMethodNotAllowed, resp.StatusCode)
}
