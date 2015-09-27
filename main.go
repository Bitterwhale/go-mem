package main

import (
	"fmt"
)

func main() {
	handler := &MemHandler{
		Processes: make([]Process, 256),
	}
	fmt.Println(handler.getProcesses())
	/*
	 * Dumping processes and their modules for the fun of it.
	 */
	for _, proc := range handler.Processes {
		/*if v.getPid() != 0{
			fmt.Println(v.getPid(), ":", v.getName())
		}*/

		proc.getModules()
		for _, module := range proc.ModuleList {
			if module.ProcessID != 0 {
				if module.getName() == "notepad.exe" {
					proc.OpenProcess(module.ProcessID)
					s, err := ReadString(&proc, 0x00000256, 24)
					fmt.Println(s, err)
				}
				//fmt.Println("\t\t", module.getName())
			}
		}

		//fmt.Println("\t\t", buffer)

	}

}
