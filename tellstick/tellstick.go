package tellstick

import (
	"fmt"
	"strconv"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

type TellstickLibrary struct {
}

type tellstickMethod uint32

const (
	turnOn  tellstickMethod = 1
	turnOff tellstickMethod = 2
	bell    tellstickMethod = 4
	toggle  tellstickMethod = 8
	dim     tellstickMethod = 16
	learn   tellstickMethod = 32
	all     tellstickMethod = turnOn | turnOff | bell | toggle | dim | learn
)

// Parameters lists the supported parametes for Tellstick devices
var Parameters = []string{"devices", "house", "unit", "code", "system", "units", "fade"}

func NewTellstickLibrary() (*TellstickLibrary, error) {
	if tellstickSupported == false {
		return nil, fmt.Errorf("%s", tellstickErrorReason)
	}
	return &TellstickLibrary{}, nil
}

// GetDeviceIds will always return a valid slice even on errors
func (tl *TellstickLibrary) GetDeviceIds() ([]int, error) {
	numDevices := tdGetNumberOfDevices()
	ids := make([]int, numDevices, numDevices)
	for i := 0; i < numDevices; i++ {
		id := tdGetDeviceId(i)
		if id == -1 {
			return ids, fmt.Errorf("unable to get device ID for %d. Reason: %s", i, tdGetErrorString())
		}
		ids[i] = int(id)
	}
	return ids, nil
}

func (tl *TellstickLibrary) GetName(id int) string {
	return tdGetName(id)
}

func (tl *TellstickLibrary) SetName(id int, name string) error {
	if tdSetName(id, name) == false {
		return fmt.Errorf("unable to set device %d name to '%s'. Reason: %s'", id, name, tdGetErrorString())
	}
	return nil
}

func (tl *TellstickLibrary) supportsMethod(id int, method tellstickMethod) bool {
	ret := tdMethods(id, int(method))
	if ret != int(method) {
		return false
	}
	return true
}

func (tl *TellstickLibrary) SupportsOnOff(id int) bool {
	return tl.supportsMethod(id, turnOn) && tl.supportsMethod(id, turnOn)
}

func (tl *TellstickLibrary) SupportsDim(id int) bool {
	return tl.supportsMethod(id, dim)
}

func (tl *TellstickLibrary) SupportsLearn(id int) bool {
	return tl.supportsMethod(id, learn)
}

func (tl *TellstickLibrary) NewDevice() (int, error) {
	id := tdAddDevice()
	if id < 0 {
		return 0, fmt.Errorf("unable to add device. Reason: %s", tdGetErrorString())
	}
	return id, nil
}

func (tl *TellstickLibrary) RemoveDevice(id int) error {
	if tdRemoveDevice(id) == false {
		return fmt.Errorf("unable to remove device %d. Reason: %s", id, tdGetErrorString())
	}
	return nil
}

func (tl *TellstickLibrary) GetProtocol(id int) string {
	return tdGetProtocol(id)
}

func (tl *TellstickLibrary) SetProtocol(id int, protocol string) error {
	if tdSetProtocol(id, protocol) == false {
		return fmt.Errorf("unable to set protocol %s to device %d. Reason: %s", protocol, id, tdGetErrorString())
	}
	return nil
}

func (tl *TellstickLibrary) GetModel(id int) string {
	return tdGetModel(id)
}

func (tl *TellstickLibrary) SetModel(id int, model string) error {
	if tdSetModel(id, model) == false {
		return fmt.Errorf("unable to set model %s to device %d. Reason: %s", model, id, tdGetErrorString())
	}
	return nil
}

func (tl *TellstickLibrary) GetParameters(id int) map[string]string {
	result := make(map[string]string)
	for _, parameter := range Parameters {
		value := tdGetDeviceParameter(id, parameter, "")
		result[parameter] = value
	}
	return result
}

func (tl *TellstickLibrary) SetParameters(id int, paramAndValues map[string]string) error {
	for parameter, value := range paramAndValues {
		if stringInSlice(parameter, Parameters) == false {
			return fmt.Errorf("unknown parameter '%s'", parameter)
		}
		if tdSetDeviceParameter(id, parameter, value) == false {
			return fmt.Errorf("unable to set parameter '%s' to '%s'. Reason: %s",
				parameter,
				value,
				tdGetErrorString())
		}
	}
	return nil
}

func (tl *TellstickLibrary) TurnOn(id int) error {
	return checkResult(tdTurnOn(id))
}

func (tl *TellstickLibrary) TurnOff(id int) error {
	return checkResult(tdTurnOff(id))
}

func (tl *TellstickLibrary) Dim(id int, level byte) error {
	return checkResult(tdDim(id, level))
}

func (tl *TellstickLibrary) Learn(id int) error {
	return checkResult(tdLearn(id))
}

func (tl *TellstickLibrary) LastCmdWasOn(id int) bool {
	return tdLastSentCommand(id, int(turnOn)) == int(turnOn)
}

func (tl *TellstickLibrary) LastDimValue(id int) byte {
	valueStr := tdLastSentValue(id)
	value, err := strconv.ParseInt(valueStr, 10, 32)
	if err != nil {
		value = 0
	}
	return byte(value)
}

func (tl *TellstickLibrary) MinDimLevel() int {
	return 0
}

func (tl *TellstickLibrary) MaxDimLevel() int {
	return 255
}

func checkResult(retVal int) error {
	errString := "unknown response code"
	// Tellstick used signed 32 bit value whi
	switch retVal {
	case 0: // TELLSTICK_SUCCESS
		errString = ""
	case -1: // TELLSTICK_ERROR_NOT_FOUND
		errString = "not found"
	case -2: // TELLSTICK_ERROR_PERMISSION_DENIED
		errString = "permission denied"
	case -3: // TELLSTICK_ERROR_DEVICE_NOT_FOUND
		errString = "device not found"
	case -4: // TELLSTICK_ERROR_METHOD_NOT_SUPPORTED
		errString = "method not supported"
	case -5: // TELLSTICK_ERROR_COMMUNICATION
		errString = "communication error"
	case -6: // TELLSTICK_ERROR_CONNECTING_SERVICE
		errString = "connecting service error"
	case -7: // TELLSTICK_ERROR_UNKNOWN_RESPONSE
		errString = "unknown response"
	case -8: // TELLSTICK_ERROR_SYNTAX
		errString = "syntax error"
	case -9: // TELLSTICK_ERROR_BROKEN_PIPE
		errString = "broken pipe"
	case -10: // TELLSTICK_ERROR_COMMUNICATING_SERVICE
		errString = "communication service error"
	case -99: // TELLSTICK_ERROR_UNKNOWN
		errString = "unknown error"
	default:
		errString = fmt.Sprintf("unknown response code %d", retVal)
	}
	if errString != "" {
		return fmt.Errorf(errString)
	}
	return nil
}
