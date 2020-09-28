package main

import (
	"fmt"
	"github.com/leangaurav/testengine/pkg/process"
)

func main() {
	fmt.Println("Starting server")
	proc := process.NewProcess("import time; print(time.sleep(10))")
	fmt.Println(proc)

	if err := proc.StartProcess(); err != nil {
		fmt.Printf("Error  starting  process : %v", err)
	} else {
		fmt.Println("Command success")
	}
}
