package main

import (
	"fmt"
	// fileHandler "microsegement/pkg/fileHandler"
	"microsegement/fileHandler"
)


func main() {

	myfolder := `target`

	// fileHandler.ListJarFile(myfolder)
	// fileHandler.TestYaml()
	list, _ := fileHandler.ListDeploymentFlie(myfolder)
	for _, v := range list{
		fmt.Println(*v)
	}

	fmt.Println("hello world")

}
