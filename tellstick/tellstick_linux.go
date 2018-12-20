package tellstick

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdlib.h>

#define bool char

int call_ri(void *f)
{
	int (*function)() = f;
	return function();
}

char * call_rs(void *f)
{
	char * (*function)() = f;
	return function();
}

int call_ri_pi(void *f, int i)
{
	int (*function)(int) = f;
	return function(i);
}

int call_ri_pii(void *f, int i1, int i2)
{
	int (*function)(int, int) = f;
	return function(i1, i2);
}

int call_ri_ps(void *f, const char *s)
{
	int (*function)(const char *) = f;
	return function(s);
}

void call_rv_ps(void *f, const char *s)
{
	void (*function)(const char *) = f;
	function(s);
}

char * call_rs_pi(void *f, int i)
{
	char * (*function)(int) = f;
	return function(i);
}

bool call_rb_pi(void *f, int i)
{
	bool (*function)(int) = f;
	return function(i);
}

bool call_rb_piss(void *f, int i, const char * s1, const char * s2)
{
	bool (*function)(int, const char *, const char *) = f;
	return function(i, s1, s2);
}

char * call_rs_piss(void *f, int i, const char * s1, const char * s2)
{
	char * (*function)(int, const char *, const char *) = f;
	return function(i, s1, s2);
}


bool call_rb_pis(void *f, int i, const char *s)
{
	bool (*function)(int, const char *s) = f;
	return function(i, s);
}

*/
import "C"

import (
	"fmt"
	"unsafe"
)

var (
	tellstickSupported   = false
	tellstickErrorReason = ""
	handle               unsafe.Pointer
	symbolPointers       = make(map[string]unsafe.Pointer)
)

func init() {
	libraryName := C.CString("libtelldus-core.so")
	defer C.free(unsafe.Pointer(libraryName))
	handle = C.dlopen(libraryName, C.RTLD_LAZY)
	if handle == nil {
		errorReason := C.GoString(C.dlerror())
		tellstickErrorReason = fmt.Sprintf("Unable to load library. Reason: %s\n", errorReason)
		return
	}
	tellstickSupported = true
}

func getSymbolPointer(name string) unsafe.Pointer {
	symbolPointer := symbolPointers[name]

	if symbolPointer == nil && tellstickSupported == false {
		panic(fmt.Sprintf("Tried to access Tellstick function, but Tellstick is not supported"))
	} else if symbolPointer == nil {
		// First time symbol is used. Load it using dlsym.
		sym := C.CString(name)
		defer C.free(unsafe.Pointer(sym))
		C.dlerror()
		symbolPointer = C.dlsym(handle, sym)
		e := C.dlerror()
		if symbolPointer == nil || e != nil {
			panic(fmt.Sprintf("Error resolving Tellstick symbol %q: %v", name, C.GoString(e)))
		}
		// Save symbol pointer so it can be reused later
		symbolPointers[name] = symbolPointer
	}

	return symbolPointer
}

func callRi(name string) int {
	return int(C.call_ri(getSymbolPointer(name)))
}

// callRs frees the returned char * after it has been copied to string
func callRs(name string) string {
	cString := C.call_rs(getSymbolPointer(name))
	defer tdReleaseString(cString)
	return C.GoString(cString)
}

func callRiPi(name string, i int) int {
	return int(C.call_ri_pi(getSymbolPointer(name), C.int(i)))
}

func callRiPii(name string, i1 int, i2 int) int {
	return int(C.call_ri_pii(getSymbolPointer(name), C.int(i1), C.int(i2)))
}

func callRiPs(name string, s string) int {
	sParam := C.CString(s)
	defer C.free(unsafe.Pointer(sParam))
	return int(C.call_ri_ps(getSymbolPointer(name), sParam))
}

func callRvPs(name string, s string) {
	sParam := C.CString(s)
	defer C.free(unsafe.Pointer(sParam))
	C.call_rv_ps(getSymbolPointer(name), sParam)
}

// callRsPi frees the returned char * after it has been copied to string
func callRsPi(name string, i int) string {
	cString := C.call_rs_pi(getSymbolPointer(name), C.int(i))
	defer tdReleaseString(cString)
	return C.GoString(cString)
}

func callRbPi(name string, i int) bool {
	if C.call_rb_pi(getSymbolPointer(name), C.int(i)) == 0 {
		return false
	}
	return true
}

func callRbPiss(name string, i int, s1 string, s2 string) bool {
	s1Param := C.CString(s1)
	defer C.free(unsafe.Pointer(s1Param))
	s2Param := C.CString(s2)
	defer C.free(unsafe.Pointer(s2Param))
	if C.call_rb_piss(getSymbolPointer(name), C.int(i), s1Param, s2Param) == 0 {
		return false
	}
	return true
}

// callRsPiss frees the returned char * after it has been copied to string
func callRsPiss(name string, i int, s1 string, s2 string) string {
	s1Param := C.CString(s1)
	defer C.free(unsafe.Pointer(s1Param))
	s2Param := C.CString(s2)
	defer C.free(unsafe.Pointer(s2Param))
	cString := C.call_rs_piss(getSymbolPointer(name), C.int(i), s1Param, s2Param)
	defer tdReleaseString(cString)
	return C.GoString(cString)
}

func callRbPis(name string, i int, s string) bool {
	sParam := C.CString(s)
	defer C.free(unsafe.Pointer(sParam))
	if C.call_rb_pis(getSymbolPointer(name), C.int(i), sParam) == 0 {
		return false
	}
	return true
}

func tdReleaseString(cString *C.char) {
	C.call_rv_ps(getSymbolPointer("tdReleaseString"), cString)
}

func tdGetNumberOfDevices() int {
	return callRi("tdGetNumberOfDevices")
}

func tdGetDeviceId(index int) int {
	return callRiPi("tdGetDeviceId", index)
}

// tdGetName will automatically free the c string using tdReleaseString
// before it is converted to a Go string
func tdGetName(id int) string {
	return callRsPi("tdGetName", id)
}

func tdSetName(id int, name string) bool {
	return callRbPis("tdSetName", id, name)
}

func tdAddDevice() int {
	return callRi("tdAddDevice")
}

func tdRemoveDevice(id int) bool {
	return callRbPi("tdRemoveDevice")
}

// tdGetErrorString will automatically free the c string using tdReleaseString
// before it is converted to a Go string
func tdGetErrorString() string {
	return callRs("tdGetErrorString")
}

func tdMethods(id int, methodsSupported int) int {
	return callRiPii("tdMethods", id, methodsSupported)
}

// tdGetDeviceParameter will automatically free the c string using tdReleaseString
// before it is converted to a Go string
func tdGetDeviceParameter(id int, name string, defaultValue string) string {
	return callRsPiss("tdGetDeviceParameter", id, name, defaultValue)
}

func tdSetDeviceParameter(id int, name string, value string) bool {
	return callRbPiss("tdSetDeviceParameter", id, name, value)
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
