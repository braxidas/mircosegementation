package fileHandler

import (
	"fmt"
	// "microsegement/mstype"
	"os"
	"path/filepath"
	"strings"

	// "gopkg.in/yaml.v3"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/yaml"
	// "k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/rest"
)

//获得folder文件夹下所有部署yaml文件
func ListDeploymentFile(folder string)([]*v1.Deployment, error){

	var (
		deploymentList []*v1.Deployment
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
				fmt.Printf("no valid deployment yaml in %s ,because %v \n", path, err)
			}
			// fmt.Println(path)
			deploymentList = append(deploymentList, deployment)
		}
		return nil
	})
	return deploymentList, err
}



//读取deployment文件
func getK8sYamlFile(yamlFilePath string)(*v1.Deployment, error){
	
	deployment := new(v1.Deployment)
	// deployment := new(mstype.Yaml2Go)
	yamlFile, err := os.ReadFile(yamlFilePath)
	if err != nil {
		return deployment, fmt.Errorf("fail to read %s, because %v ", yamlFilePath, err)
	}

	err = yaml.Unmarshal(yamlFile, deployment)
	if err != nil {
		return deployment, fmt.Errorf("fail to unmarshal %s, because %v ", yamlFilePath, err)
	}

	// if deployment.Spec.Template.Spec == "" && deployment.Kind == ""  {
	// 	return deployment,fmt.Errorf("missing required fields")
	// }
	return deployment, nil
}
