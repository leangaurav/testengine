package main

import (
	"fmt"
	"github.com/leangaurav/testengine/pkg/test"
)

func main() {
	fmt.Println("Starting server")
	id := "abcd1234"
	code := "import time; print('sleeping'); time.sleep(1); print('Woke up')"
	
	test := test.NewTest(id, code, test.PYTHON, "/tmp", []test.TestCase{})
	if err := test.Run(); err != nil {
		fmt.Printf("Error  running test : %v", err)
	} else {
		fmt.Println("Test ran")
	}
}
