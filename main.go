package main

import "fmt"

func main() {
	handler := &MemHandler{
		Processes: make([]Process, 256),
	}
	handler.getProcesses()
	fmt.Println(handler.Processes[16].getName())
}
