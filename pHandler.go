// +build windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"
	"log"
)

// Windows API functions
var (
	modKernel32                  = syscall.NewLazyDLL("kernel32.dll")
	procCloseHandle              = modKernel32.NewProc("CloseHandle")
	procCreateToolhelp32Snapshot = modKernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First           = modKernel32.NewProc("Process32FirstW")
	procProcess32Next            = modKernel32.NewProc("Process32NextW")
	Module32First 				 = modKernel32.NewProc("Module32FirstW")
	Module32Next 				 = modKernel32.NewProc("Module32NextW")
)

// Some constants from the Windows API
const (
	ERROR_NO_MORE_FILES = 0x12
	MAX_PATH            = 260
	MAX_MODULE_NAME32 	= 255
)

// PROCESSENTRY32 is the Windows API structure that contains a process's
// information.
type PROCESSENTRY32 struct {
	Size              uint32
	CntUsage          uint32
	ProcessID         uint32
	DefaultHeapID     uintptr
	ModuleID          uint32
	CntThreads        uint32
	ParentProcessID   uint32
	PriorityClassBase int32
	Flags             uint32
	ExeFile           [MAX_PATH]uint16
}

type MODULEENTRY32 struct {
	Size 			uint32
	ModuleID		uint32
	ProcessID 		uint32
	GlobalUsage 	uint32
	ProcessUsage 	uint32
	modBaseAddr 	*uint8
	modBaseSize 	uint32
	hModule 		uintptr
	szModule 		[MAX_MODULE_NAME32+1]uint16
	ExeFile         [MAX_PATH]uint16

}

func getProcessBase (pid uintptr) uint8 {
	var baseAddress uint8
	handle, _, _ := procCreateToolhelp32Snapshot.Call(
		0x00000008,
		pid)
	defer procCloseHandle.Call(handle)
	var entry MODULEENTRY32
	
	entry.Size = uint32(unsafe.Sizeof(entry))

	ret, _, _ := Module32First.Call(handle, uintptr(unsafe.Pointer(&entry)))
	if ret == 0 {
		log.Panic("?!")
	}
	fmt.Println(entry.modBaseAddr)
	return baseAddress

}





/**
func (p *MODULEENTRY32) BaseAddress() string {
	return string(*p.modBaseAddr)
}
*/

// WindowsProcess is an implementation of Process for Windows.
type WindowsProcess struct {
	pid  int
	ppid int
	exe  string
}

func (p *WindowsProcess) Pid() int {
	return p.pid
}

func (p *WindowsProcess) PPid() int {
	return p.ppid
}

func (p *WindowsProcess) Executable() string {
	return p.exe
}

func newWindowsProcess(e *PROCESSENTRY32) *WindowsProcess {
	// Find when the string ends for decoding
	end := 0
	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}

	return &WindowsProcess{
		pid:  int(e.ProcessID),
		ppid: int(e.ParentProcessID),
		exe:  syscall.UTF16ToString(e.ExeFile[:end]),
	}
}

/**func findProcess(pid int) (Process, error) {
	ps, err := processes()
	if err != nil {
		return nil, err
	}

	for _, p := range ps {
		if p.Pid() == pid {
			return p, nil
		}
	}

	return nil, nil
}*/



func getProcessModules(pid uintptr) []MODULEENTRY32 {
	results := make([]MODULEENTRY32, 128)
	handle, _, _ := procCreateToolhelp32Snapshot.Call(
		0x00000008,
		pid)

	defer procCloseHandle.Call(handle)

	var entry MODULEENTRY32
	
	entry.Size = uint32(unsafe.Sizeof(entry))



	ret, _, _ := Module32First.Call(handle, uintptr(unsafe.Pointer(&entry)))
	if ret == 0 {
		log.Panic("NO MODULES!?!¤(&/¤&3452))")
	}


	for {
		results = append(results, entry)
		ret, _, _ := Module32Next.Call(handle, uintptr(unsafe.Pointer(&entry)))
		if ret == 0 {
			break
		}
		fmt.Println(entry.modBaseAddr)

	}
	return results
	
}




func processes() ([]PROCESSENTRY32, error) {


	handle, _, _ := procCreateToolhelp32Snapshot.Call(
		0x00000002,
		2372)
	if handle < 0 {
		return nil, syscall.GetLastError()
	}
	defer procCloseHandle.Call(handle)

	var entry PROCESSENTRY32
	entry.Size = uint32(unsafe.Sizeof(entry))


	ret, _, _ := procProcess32First.Call(handle, uintptr(unsafe.Pointer(&entry)))
	if ret == 0 {
		return nil, fmt.Errorf("Error retrieving process info.")
	}

	results := make([]PROCESSENTRY32, 0, 50)
	for {
		results = append(results, entry)
		ret, _, _ := procProcess32Next.Call(handle, uintptr(unsafe.Pointer(&entry)))
		if ret == 0 {
			break
		}
	}
	//fmt.Println(results2)
	return results, nil
}

