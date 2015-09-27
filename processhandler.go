// +build windows

package main

import (
	"syscall"
	"unsafe"
)

// Windows API functions
var (
	modKernel32   = syscall.NewLazyDLL("kernel32.dll")
	Module32First = modKernel32.NewProc("Module32FirstW")
	Module32Next  = modKernel32.NewProc("Module32NextW")
)

// Some constants from the Windows API
const (
	ERROR_NO_MORE_FILES = 0x12
	MAX_PATH            = 260
	MAX_MODULE_NAME32   = 255
)

type MemHandler struct {
	Processes []Process
}

func (m *MemHandler) getProcesses() error {

	snapshot, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return err
	}
	defer syscall.CloseHandle(snapshot)

	var procEntry syscall.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))
	if err = syscall.Process32First(snapshot, &procEntry); err != nil {
		return err
	}

	results := make([]syscall.ProcessEntry32, 0, 50)
	for {
		results = append(results, procEntry)
		if err := syscall.Process32Next(snapshot, &procEntry); err != nil {
			break
		}
	}
	//fmt.Println(results2)
	var p Process
	for _, v := range results {
		if v.ProcessID != 0 {
			p.PROCESSENTRY = v
			m.Processes = append(m.Processes, p)
		}
	}
	return nil
}
