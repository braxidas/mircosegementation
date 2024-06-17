package serviceHandler

import (
	"fmt"
	"microsegement/fileHandler"
	"microsegement/mstype"
)

func RegisterService(folder string)([]*mstype.K8sService, error){
	var k8sServiceList []*mstype.K8sService

	applicationList,pathList,err := fileHandler.ListJarFile(folder)
	if(err != nil){
		fmt.Println(err)
	}

	deploymentList, err := fileHandler.ListDeploymentFile(folder)
	if(err != nil){
		return k8sServiceList,err
	}
	
	for i, _:= range  pathList{
		k8sService := new(mstype.K8sService)
		k8sService.PodName = deploymentList[i].Metadata.Name
		k8sService.FilePath = pathList[i]
		k8sService.ApplicationName, err = applicationList[i].GetApplicationName()
		if(err != nil){
			fmt.Println(err, pathList)
		}
		k8sServiceList = append(k8sServiceList, k8sService)
	}
	return k8sServiceList, nil
}

