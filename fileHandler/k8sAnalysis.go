package fileHandler

import (
	"fmt"
	"microsegement/mstype"
	"os"
	"path/filepath"
	"strings"
	"gopkg.in/yaml.v3"
)

//获得folder文件夹下所有部署yaml文件
func ListDeploymentFile(folder string)([]*mstype.Yaml2Go, error){

	var (
		deploymentList []*mstype.Yaml2Go
	)
	folder = getParentDirectory(folder)//获得上级路径

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() &&  (strings.HasSuffix(info.Name(), ".yml") || strings.HasSuffix(info.Name(), ".yaml")) {
			// conf, serviceName, err := parser.ParseYaml(path)
			deployment, err := getK8sYamlFile(path)
			if err != nil {
				fmt.Printf("no valid deployment ymal in %s ,because %v \n", path, err)
			}
			// fmt.Println(path)
			deploymentList = append(deploymentList, deployment)
		}
		return nil
	})
	return deploymentList, err
}



//读取deployment文件
func getK8sYamlFile(yamlFilePath string)(*mstype.Yaml2Go, error){

	deployment := new(mstype.Yaml2Go)
	yamlFile, err := os.ReadFile(yamlFilePath)
	if err != nil {
		return deployment, fmt.Errorf("fail to read %s, because %v ", yamlFilePath, err)
	}

	err = yaml.Unmarshal(yamlFile, deployment)
	if err != nil {
		return deployment, fmt.Errorf("fail to unmarshal %s, because %v ", yamlFilePath, err)
	}

	if deployment.ApiVersion == "" && deployment.Kind == ""  {
		return deployment,fmt.Errorf("missing required fields")
	}
	return deployment, nil
}
