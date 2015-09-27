package main

import (
	"encoding/binary"
	"fmt"
	"syscall"
	"unsafe"
)

/**
 * Reads the memory from 'process', starting at 'address' up to 'address + size'.
 * Outputs into 'buffer'.
 */

func ReadProcessMemory(process syscall.Handle, address uintptr, buffer []byte, size uint32, bytesRead uint32) (uint32, error) {
	ret, _, err := modKernel32.NewProc("ReadProcessMemory").Call(uintptr(process), address, uintptr(unsafe.Pointer(&buffer[0])), uintptr(size), uintptr(unsafe.Pointer(&bytesRead)))
	if ret == 0 {
		return 0, err
	}

	return uint32(ret), nil
}

func ReadString(process *Process, offset uint32, length uint32) (string, error) {
	buffer := make([]byte, length)
	fmt.Println(process.BaseAddress())
	_, err := ReadProcessMemory(process.Handle, process.BaseAddress(), buffer, 24, 0)
	if err != nil {
		panic(err)
	}
	return UTF16BytesToString(buffer, binary.LittleEndian), nil
}
