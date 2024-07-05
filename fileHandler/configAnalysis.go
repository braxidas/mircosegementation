package fileHandler

import (
	"fmt"
	"microsegement/mstype"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

//根据nacos的文件路径和service的app名查找相关配置文件,并修改k8sService的applicationList
func ListNacosYamlFile(folder string, k8sService *mstype.K8sService) (error) {
	if folder == ""{
		return nil
	}

	// folder = getParentDirectory(folder) //获得上级路径

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".yml") || strings.HasSuffix(info.Name(), ".yaml")) && strings.Contains(info.Name(), k8sService.ApplicationName){
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


func getNacosConfig(configPath string)(*mstype.Application, error){
	application := new(mstype.Application)
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return application, fmt.Errorf("fail to read %s,because %v", configPath, err)
	}
	err = yaml.Unmarshal(yamlFile, application)
	if err != nil {
		return application, fmt.Errorf("fail to unmarshal %s, because %v ", configPath, err)
	}
	return application, nil
}
