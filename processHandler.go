package main


import
(
	"fmt"
	"syscall"
	"unsafe"
)

var
(
	psapi = syscall.NewLazyDLL("psapi.dll")
	procEnumProcesses = psapi.NewProc("EnumProcesses")
)

func EnumProcesses(processIds []uint32, cb uint32, bytesReturned *uint32) bool {
	ret, _, _ := procEnumProcesses.Call(
		uintptr(unsafe.Pointer(&processIds[0])),
		uintptr(cb),
		uintptr(unsafe.Pointer(bytesReturned)))
	 
	return ret != 0

}

func main() {
	processIds := make([]uint32, 256)
	bytesReturned := uint32(256)
	EnumProcesses(processIds, 256, &bytesReturned)
	fmt.Println(processIds)
}
	