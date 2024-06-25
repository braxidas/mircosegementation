package main

import (
	
	// fileHandler "microsegement/pkg/fileHandler"
	// "microsegement/fileHandler"
	"microsegement/cmd"
	// "microsegement/serviceHandler"
)


func main() {

	err := cmd.Execute()
	if err != nil {
		return
	}

}
