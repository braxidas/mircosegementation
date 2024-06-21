package main

import (
	"fmt"
	"os"
	// fileHandler "microsegement/pkg/fileHandler"
	// "microsegement/fileHandler"
	"microsegement/serviceHandler"
)


func main() {

	// myfolder := `C:\Users\li_sh\Desktop\msexample\example\0ruoyiExample`
	if len(os.Args) < 2{
		fmt.Println("please input path to scan")
		return
	}
	myfolder := os.Args[1]
	fmt.Println("start to scan",myfolder)

	k8sServiceList, err := serviceHandler.RegisterService(myfolder)
	if err != nil{
		fmt.Println("fail to service register%v\n", err)
	}
	err = serviceHandler.DiscoverService(k8sServiceList)
	if err != nil{
		fmt.Println("fail to service discoverr%v\n", err)
	}
	fmt.Println("finish scan",myfolder)

}

// func TestScan(myfolder string){
// 	list1, _ := fileHandler.ListDeploymentFile(myfolder)
// 	list2,list3, _ := fileHandler.ListJarFile(myfolder)
// 	// list, _ := serviceHandler.RegisterService(myfolder)
// 	for i, _ := range list1{
// 		fmt.Println(list1[i].Metadata.Name)
// 	}
// 	for _, v := range list2{
// 		fmt.Println(*v)
// 	}
// 	for _, v := range list3{
// 		fmt.Println(v)
// 	}
// }

// func TestRegister(myfolder string){
// 	list, _ := serviceHandler.RegisterService(myfolder)
// 	for _, v := range list{
// 		fmt.Println(*v)
// 	}
// }

// func TestDiscovery(){
// 	serviceHandler.DiscoverService(nil)
// }