package tellstick

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	tellstickSupported   = false
	tellstickErrorReason = ""
	library              *syscall.LazyDLL
	lazyProcs            = make(map[string]*syscall.LazyProc)
)

func init() {
	libraryName := "TelldusCore.dll"
	library = syscall.NewLazyDLL(libraryName)
	err := library.Load()
	if err != nil {
		tellstickErrorReason = fmt.Sprintf("Unable to load library. Reason: %s\n", err)
		return
	}
	tellstickSupported = true
}

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

func lazy(name string) *syscall.LazyProc {
	lazyProc := lazyProcs[name]
	if lazyProc == nil && tellstickSupported == false {
		panic(fmt.Sprintf("Tried to access Tellstick function, but Tellstick is not supported"))
	} else if lazyProc == nil {
		// First time proc is is used. Load it using NewProc.
		lazyProc = library.NewProc(name)
		// Save symbol pointer so it can be reused later
		lazyProcs[name] = lazyProc
	}

	return lazyProc
}

func tdReleaseString(cString uintptr) {
	lazy("tdReleaseString").Call(cString)
}

func tdGetNumberOfDevices() int {
	ret, _, _ := lazy("tdGetNumberOfDevices").Call()
	return int(ret)
}

func tdGetDeviceId(index int) int {
	ret, _, _ := lazy("tdGetDeviceId").Call(uintptr(index))
	return int(ret)
}

// tdGetName will automatically free the c string using tdReleaseString
// before it is converted to a Go string
func tdGetName(id int) string {
	ret, _, _ := lazy("tdGetName").Call(uintptr(id))
	defer tdReleaseString(ret)
	return uintptrToString(ret)
}

func tdSetName(id int, name string) bool {
	ret, _, _ := lazy("tdSetName").Call(uintptr(id), stringToUintptr(name))
	if ret == 0 {
		return false
	}
	return true
}

func tdAddDevice() int {
	ret, _, _ := lazy("tdAddDevice").Call()
	return int(ret)
}

func tdRemoveDevice(id int) bool {
	ret, _, _ := lazy("tdRemoveDevice").Call(uintptr(id))
	if ret == 0 {
		return false
	}
	return true
}

// tdGetErrorString will automatically free the c string using tdReleaseString
// before it is converted to a Go string
func tdGetErrorString() string {
	ret, _, _ := lazy("tdGetErrorString").Call()
	defer tdReleaseString(ret)
	return uintptrToString(ret)
}

func tdMethods(id int, methodsSupported int) int {
	ret, _, _ := lazy("tdMethods").Call(uintptr(id), uintptr(methodsSupported))
	return int(ret)
}
func tdGetProtocol(id int) string {
	ret, _, _ := lazy("tdGetProtocol").Call(uintptr(id))
	defer tdReleaseString(ret)
	return uintptrToString(ret)
}

func tdSetProtocol(id int, protocol string) bool {
	ret, _, _ := lazy("tdSetProtocol").Call(uintptr(id), stringToUintptr(protocol))
	if ret == 0 {
		return false
	}
	return true
}

func tdGetModel(id int) string {
	ret, _, _ := lazy("tdGetModel").Call(uintptr(id))
	defer tdReleaseString(ret)
	return uintptrToString(ret)
}

func tdSetModel(id int, model string) bool {
	ret, _, _ := lazy("tdSetModel").Call(uintptr(id), stringToUintptr(model))
	if ret == 0 {
		return false
	}
	return true
}

// tdGetDeviceParameter will automatically free the c string using tdReleaseString
// before it is converted to a Go string
func tdGetDeviceParameter(id int, name string, defaultValue string) string {
	ret, _, _ := lazy("tdGetDeviceParameter").Call(uintptr(id),
		stringToUintptr(name), stringToUintptr(defaultValue))
	defer tdReleaseString(ret)
	return uintptrToString(ret)
}

func tdSetDeviceParameter(id int, name string, value string) bool {
	ret, _, _ := lazy("tdSetDeviceParameter").Call(uintptr(id),
		stringToUintptr(name), stringToUintptr(value))
	if ret == 0 {
		return false
	}
	return true
}

/*
func main() {
	if tellstickSupported {
		fmt.Printf("Tellstick is supported\n")
	} else {
		fmt.Printf("Tellstick not supported. %s\n", tellstickErrorReason)
	}
	fmt.Printf("Number of devices: %d\n", tdGetNumberOfDevices())
	fmt.Printf("Name ID 1: %s\n", tdGetName(1))
}
*/
