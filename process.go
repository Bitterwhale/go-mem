package main

import (
	"syscall"
	"unsafe"
)

type Process struct {
	PROCESSENTRY syscall.ProcessEntry32
	ModuleList   []ModuleEntry32
	Handle       syscall.Handle
}

type ProcessEntry32 syscall.ProcessEntry32

type ModuleEntry32 struct {
	Size         uint32
	ModuleID     uint32
	ProcessID    uint32
	GlobalUsage  uint32
	ProcessUsage uint32
	modBaseAddr  *uint8
	modBaseSize  uint32
	hModule      uintptr
	szModule     [MAX_MODULE_NAME32 + 1]uint16
	ExeFile      [MAX_PATH]uint16
}

func (p *Process) getPid() uint32 {
	return p.PROCESSENTRY.ProcessID
}

func (p *Process) getName() string {
	file := p.PROCESSENTRY.ExeFile
	var str string
	for _, value := range file {
		if int(value) == 0 {
			break
		}
		str += string(int(value))
	}
	return str
}

func (p *Process) OpenProcess(pid uint32) error {
	handle, err := syscall.OpenProcess(0x30, false, pid) // Opens with all rights
	if err != nil {
		return err
	}
	p.Handle = handle
	return nil
}

func (m *ModuleEntry32) getName() string {
	var str string
	for _, value := range m.szModule {
		if int(value) == 0 {
			break
		}
		str += string(int(value))
	}
	return str
}

func (m *ModuleEntry32) getFullPath() string {
	var str string
	for _, value := range m.ExeFile {
		if int(value) == 0 {
			break
		}
		str += string(int(value))
	}
	return str
}

func (p *Process) BaseAddress() uintptr {
	handle, _ := syscall.CreateToolhelp32Snapshot(
		syscall.TH32CS_SNAPMODULE,
		p.PROCESSENTRY.ProcessID)
	defer syscall.Close(handle)
	var entry ModuleEntry32

	entry.Size = uint32(unsafe.Sizeof(entry))

	ret, _, _ := Module32First.Call(uintptr(handle), uintptr(unsafe.Pointer(&entry)))
	if ret == 0 {
		panic(syscall.GetLastError())
	}
	return uintptr(*entry.modBaseAddr)

}

func (p *Process) getModules() error {

	handle, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPMODULE, p.getPid())
	if err != nil {
		return err
	}
	defer syscall.Close(handle)

	var entry ModuleEntry32

	entry.Size = uint32(unsafe.Sizeof(entry))

	ret, _, _ := Module32First.Call(uintptr(handle), uintptr(unsafe.Pointer(&entry)))
	if ret == 0 {
		return syscall.GetLastError() //log.Panic("NO MODULES!?!¤(&/¤&3452))")
	}
	results := make([]ModuleEntry32, 128)
	for {
		results = append(results, entry)
		ret, _, _ := Module32Next.Call(uintptr(handle), uintptr(unsafe.Pointer(&entry)))
		if ret == 0 {
			break
		}
	}
	p.ModuleList = results
	return nil
}
