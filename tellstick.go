package tellstick

import "fmt"

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

var parameters = []string{"devices", "house", "unit", "code", "system", "units", "fade"}

func NewTellstickLibrary() (*TellstickLibrary, error) {
	if tellstickSupported == false {
		return nil, fmt.Errorf("%s", tellstickErrorReason)
	}
	return &TellstickLibrary{}, nil
}

func (tl *TellstickLibrary) GetDeviceIds() ([]int, error) {
	numDevices := tdGetNumberOfDevices()
	ids := make([]int, numDevices, numDevices)
	for i := 0; i < numDevices; i++ {
		id := tdGetDeviceId(i)
		if id == -1 {
			return nil, fmt.Errorf("unable to get device ID for %d. Reason: %s", i, tdGetErrorString())
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
	for _, parameter := range parameters {
		value := tdGetDeviceParameter(id, parameter, "")
		result[parameter] = value
	}
	return result
}

func (tl *TellstickLibrary) SetParameters(id int, paramAndValues map[string]string) error {
	for parameter, value := range paramAndValues {
		if stringInSlice(parameter, parameters) == false {
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

/*
func main() {
	tl, err := NewTellstickLibrary()
	if err != nil {
		panic(fmt.Sprintf("Error: %s\n", err))
	}
	ids, _ := tl.GetDeviceIds()
	for _, id := range ids {
		fmt.Printf("Id %d: Name: '%s' OnOff: %t Dim: %t Learn: %t\n", id, tl.GetName(id), tl.SupportsOnOff(id), tl.SupportsDim(id), tl.SupportsLearn(id))
	}
	parameters := tl.GetParameters(1)
	for key, value := range parameters {
		fmt.Printf("%s = %s\n", key, value)
	}
	tl.SetName(11, "elvan")
}
*/
