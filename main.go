package main

import (
	
	// fileHandler "microsegement/pkg/fileHandler"
	// "microsegement/fileHandler"
	"microsegement/cmd"
	// "microsegement/soot"
	// "microsegement/serviceHandler"
)


func main() {
	err := cmd.Execute()
	if err != nil {
		return
	}

}


