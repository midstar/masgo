package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

const maxStringSize = 256

func uintptrToString(u uintptr) string {
	pointer := unsafe.Pointer(u)
	bytes := (*[maxStringSize]byte)(pointer)[:]

	// Find 0 termination of string
	n := -1
	for i, b := range bytes {
		if b == 0 {
			break
		}
		n = i
	}

	return string(bytes[:n+1])
}

func stringToUintptr(s string) uintptr {
	bytes := []byte(s)
	// Add 0 termination
	bytes = append(bytes, 0)
	return uintptr(unsafe.Pointer(&bytes[0]))
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

type TellstickLibrary struct {
	tdGetNumberOfDevices *syscall.LazyProc
	tdGetDeviceID        *syscall.LazyProc
	tdGetName            *syscall.LazyProc
	tdSetName            *syscall.LazyProc
	tdAddDevice          *syscall.LazyProc
	tdReleaseString      *syscall.LazyProc
	tdGetErrorString     *syscall.LazyProc
	tdMethods            *syscall.LazyProc
	tdGetDeviceParameter *syscall.LazyProc
	tdSetDeviceParameter *syscall.LazyProc
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
	libraryName := "TelldusCore.dll"
	library := syscall.NewLazyDLL(libraryName)
	err := library.Load()
	if err != nil {
		return nil, err
	}

	return &TellstickLibrary{
			tdGetNumberOfDevices: library.NewProc("tdGetNumberOfDevices"),
			tdGetDeviceID:        library.NewProc("tdGetDeviceId"),
			tdGetName:            library.NewProc("tdGetName"),
			tdSetName:            library.NewProc("tdSetName"),
			tdAddDevice:          library.NewProc("tdAddDevice"),
			tdReleaseString:      library.NewProc("tdReleaseString"),
			tdGetErrorString:     library.NewProc("tdGetErrorString"),
			tdMethods:            library.NewProc("tdMethods"),
			tdGetDeviceParameter: library.NewProc("tdGetDeviceParameter"),
			tdSetDeviceParameter: library.NewProc("tdSetDeviceParameter")},
		nil
}

func (tl *TellstickLibrary) GetDeviceIds() ([]int, error) {
	numDevices, _, _ := tl.tdGetNumberOfDevices.Call()
	ids := make([]int, numDevices, numDevices)
	var i uintptr
	for i = 0; i < numDevices; i++ {
		idRet, _, _ := tl.tdGetDeviceID.Call(i)
		id := int32(idRet)
		if id == -1 {
			return nil, fmt.Errorf("unable to get device ID for %d. Reason: %s", i, tl.getErrorString())
		}
		ids[i] = int(id)
	}
	return ids, nil
}

func (tl *TellstickLibrary) GetName(id int) string {
	nameRet, _, _ := tl.tdGetName.Call(uintptr(id))
	name := uintptrToString(nameRet)
	tl.tdReleaseString.Call(nameRet)
	return name
}

func (tl *TellstickLibrary) SetName(id int, name string) error {
	successRet, _, _ := tl.tdSetName.Call(uintptr(id), stringToUintptr(name))
	if successRet == 0 {
		return fmt.Errorf("unable to set device %d name to '%s'. Reason: %s'", id, name, tl.getErrorString())
	}
	return nil
}

func (tl *TellstickLibrary) supportsMethod(id int, method tellstickMethod) bool {
	ret, _, _ := tl.tdMethods.Call(uintptr(id), uintptr(method))
	if ret != uintptr(method) {
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
	idRet, _, _ := tl.tdAddDevice.Call()
	id := int(int32(idRet))
	if id < 0 {
		return 0, fmt.Errorf("unable to add device. Reason: %s", tl.getErrorString())
	}
	return id, nil
}

func (tl *TellstickLibrary) GetParameters(id int) map[string]string {
	result := make(map[string]string)
	for _, parameter := range parameters {
		valueRet, _, _ := tl.tdGetDeviceParameter.Call(uintptr(id),
			stringToUintptr(parameter),
			stringToUintptr(""))
		value := uintptrToString(valueRet)
		tl.tdReleaseString.Call(valueRet)
		result[parameter] = value

	}
	return result
}

func (tl *TellstickLibrary) SetParameters(id int, paramAndValues map[string]string) error {
	for parameter, value := range paramAndValues {
		if stringInSlice(parameter, parameters) == false {
			return fmt.Errorf("unknown parameter '%s'", parameter)
		}
		resultRet, _, _ := tl.tdSetDeviceParameter.Call(uintptr(id),
			stringToUintptr(parameter),
			stringToUintptr(value))
		if resultRet != 1 {
			return fmt.Errorf("unable to set parameter '%s' to '%s'. Reason: %s",
				parameter,
				value,
				tl.getErrorString())
		}
	}
	return nil
}

func (tl *TellstickLibrary) getErrorString() string {
	errorStringRet, _, _ := tl.tdGetErrorString.Call()
	errorString := uintptrToString(errorStringRet)
	tl.tdReleaseString.Call(errorStringRet)
	return errorString
}

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
