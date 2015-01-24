package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

func hddspace() {
	h := syscall.MustLoadDLL("kernel32.dll")
	c := h.MustFindProc("GetDiskFreeSpaceExW")
	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)
	_, _, err := c.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("C:"))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(lpFreeBytesAvailable)
}

func main() {
	var modKernel32 = syscall.NewLazyDLL("kernel32.dll")

	var readProcessMemory = modKernel32.NewProc("ReadProcessMemory")
}
