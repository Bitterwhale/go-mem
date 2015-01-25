package main


import
(
	"fmt"
	"os"
	"log"
	"syscall"
	"unsafe"
)

var
(
	psapi = syscall.NewLazyDLL("psapi.dll")
	procEnumProcesses = psapi.NewProc("EnumProcesses")
	kernel32 = syscall.NewLazyDLL("Kernel32.dll")
	getModuleBaseName = kernel32.NewProc("K32GetModuleBaseNameA")
)

func EnumProcesses(processIds []uint32, cb uint32, bytesReturned *uint32) bool {
	ret, _, _ := procEnumProcesses.Call(
		uintptr(unsafe.Pointer(&processIds[0])),
		uintptr(cb),
		uintptr(unsafe.Pointer(bytesReturned)))
	 
	return ret != 0

}

func GetProcessHandle(pid int) uintptr {
	handle, _ := syscall.OpenProcess(0x0010, false, uint32(pid))
	return uintptr(handle)
}



func GetProcessName(p *os.Process) string {
	handle := GetProcessHandle(p.Pid)
	ret := ""
	_, _, e := getModuleBaseName.Call(
		handle,
		0,
		uintptr(unsafe.Pointer(syscall.StringBytePtr(ret))),
		32)
	if e != nil {
		log.Panic(e)
	}
	return ret
}




func main() {
	processIds := make([]uint32, 256)
	bytesReturned := uint32(256)
	EnumProcesses(processIds, 256, &bytesReturned)
	//	fmt.Println(processIds)

	for i := 0; i < len(processIds); i++ {
		process, err := os.FindProcess(int(processIds[i]))
		if err == nil {
			fmt.Println(processIds[i], "\t: ", GetProcessName(process))
		}
	}



}