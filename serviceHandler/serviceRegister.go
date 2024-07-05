package serviceHandler

import (
	"fmt"
	"microsegement/fileHandler"
	"microsegement/mstype"
)

var(
	name2K8sService map[string]*mstype.K8sService//通过服务名称取得该微服务的信息	
)


//扫描jar包和deployment文件 返回扫描得到的服务列表
func RegisterService(folder string, nacosFolder string) ([]*mstype.K8sService, error) {
	var k8sServiceList []*mstype.K8sService
	name2K8sService = make(map[string]*mstype.K8sService)

	applicationList, pathList, err := fileHandler.ListJarFile(folder)
	if err != nil {
		fmt.Println(err)
	}

	for i := range pathList {
		k8sService := new(mstype.K8sService)
		k8sService.Egress = make(map[*mstype.K8sService]struct{})
		k8sService.Ingress = make(map[*mstype.K8sService]struct{})
		k8sService.FilePath = pathList[i]
		k8sService.ApplicationName, err = applicationList[i].GetApplicationName()
		k8sService.ApplicationList = append(k8sService.ApplicationList, applicationList[i])
		if err != nil {
			fmt.Println(err, pathList[i])
		}
		deploymentList, err := fileHandler.ListDeploymentFile(pathList[i])
		if err != nil {
			fmt.Println(err, pathList[i])
		}
		if len(deploymentList) > 0 {//有deployment文件的情况
			// k8sService.PodName = deploymentList[0].Spec.Template.Metadata.Labels.App
			k8sService.PodName = deploymentList[0].Spec.Template.Metadata.Labels.App
			name2K8sService[k8sService.ApplicationName] = k8sService
			
		} else {
			fmt.Println("no deployment file in ", pathList[i])
		}
		k8sServiceList = append(k8sServiceList, k8sService)
	}

	for _, v := range k8sServiceList{
		if v.PodName != ""{
			err = fileHandler.ListNacosYamlFile(nacosFolder, v)
			if err != nil{
				fmt.Println(err)
			}
		}
	}

	return k8sServiceList, nil
}
