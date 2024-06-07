package main

import (
	"fmt"
	// fileHandler "microsegement/pkg/fileHandler"
	fileHandler "microsegement/fileHandler"
)


func main() {

	myfolder := `target`

	// fileHandler.ListJarFile(myfolder)
	fileHandler.ListJarFile(myfolder)

	fmt.Println("hello world")

}
