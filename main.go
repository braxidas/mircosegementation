package main

import (
	
	// fileHandler "microsegement/pkg/fileHandler"
	// "microsegement/fileHandler"
	"microsegement/cmd"
	// "microsegement/soot"
	// "microsegement/serviceHandler"
)


func main() {
	if err := cmd.Execute();err != nil {
		return
	}
}


