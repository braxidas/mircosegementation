package fileHandler

import (
	"fmt"
	"microsegement/mstype"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

//根据nacos的文件路径和service的app名查找相关配置文件
func ListNacosYamlFile(folder string, k8sService *mstype.K8sService) (error) {
	if folder == ""{
		return nil
	}

	// folder = getParentDirectory(folder) //获得上级路径

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") && strings.Contains(info.Name(), k8sService.ApplicationName){
			// conf, serviceName, err := parser.ParseYaml(path)
			application, err := getNacosConfig(path)
			if err != nil {
				fmt.Printf("not find ipblock information  in %s ,because %v \n", path, err)
			}
			// fmt.Println(path)
			k8sService.ApplicationList = append(k8sService.ApplicationList, application)
		}
		return nil
	})
	return err
}

//获得配置yaml文件中的ipblock信息
/*
Egress
	nacos注册中心地址
	nacos服务器上export的配置文件中的ip地址
Ingress
	jar包中的server.port表示对外暴露的端口

*/
func getNacosConfig(configPath string)(*mstype.Application, error){
	application := new(mstype.Application)
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return application, fmt.Errorf("fail to read ", configPath, err)
	}
	err = yaml.Unmarshal(yamlFile, application)
	if err != nil {
		return application, fmt.Errorf("fail to unmarshal ", configPath, err)
	}
	return application, nil
}
