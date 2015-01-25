package main

import "fmt"

func main() {
	handler := &MemHandler{
		Processes: make([]Process, 256),
	}
	handler.getProcesses()
	for _, v := range handler.Processes{
		if v.getPid() != 0{
			fmt.Println(v.getPid())
		}
	}
}
