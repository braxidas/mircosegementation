package main

import (
	"fmt"
	// fileHandler "microsegement/pkg/fileHandler"
	"microsegement/fileHandler"
	"microsegement/serviceHandler"
)


func main() {

	myfolder := `C:\Users\li_sh\Desktop\msexample\example\0ruoyiExample`

	// fileHandler.ListJarFile(myfolder)
	// fileHandler.TestYaml()
	TestDiscovery()
	fmt.Println(myfolder)

}

func TestScan(myfolder string){
	list1, _ := fileHandler.ListDeploymentFile(myfolder)
	list2,list3, _ := fileHandler.ListJarFile(myfolder)
	// list, _ := serviceHandler.RegisterService(myfolder)
	for i, _ := range list1{
		fmt.Println(list1[i].Metadata.Name)
	}
	for _, v := range list2{
		fmt.Println(*v)
	}
	for _, v := range list3{
		fmt.Println(v)
	}
}

func TestRegister(myfolder string){
	list, _ := serviceHandler.RegisterService(myfolder)
	for _, v := range list{
		fmt.Println(*v)
	}
}

func TestDiscovery(){
	serviceHandler.DiscoverService(nil)
}