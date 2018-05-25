package main

import (
	"fmt"
	"os"
)

func main() {
	switch os.Args[1] {
	case "read":
		read()
	case "write":
		write()
	default:
		panic(fmt.Errorf("unsupported operation %s", os.Args[1]))
	}
}
