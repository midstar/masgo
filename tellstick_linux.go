package main

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdlib.h>


int call_ri(void *f)
{
	int (*function)() = f;
	return function();
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

*/
import "C"

import (
	"fmt"
	"unsafe"
)

var (
	tellstickSupported = false
	tellstickErrorReason = ""
	handle unsafe.Pointer
	symbolPointers = make(map[string]unsafe.Pointer) 
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
	}

	return symbolPointer
}

func callRi(name string) int {	
	return int(C.call_ri(getSymbolPointer(name)))
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

func tdGetNumberOfDevices() int {
	return callRi("tdGetNumberOfDevices")
}

func tdReleaseString(theString string) {
	callRvPs("tdReleaseString", theString) 
}

func main() {
	if tellstickSupported {
		fmt.Printf("Tellstick is supported\n")
	} else {
		fmt.Printf("Tellstick not supported. %s\n", tellstickErrorReason) 
	}
	fmt.Printf("Number of devices: %d\n", tdGetNumberOfDevices())
}
	
	
