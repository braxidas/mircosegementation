package filehandler

import(
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"gopkg.in/yaml.v3"
	// k8syaml "sigs.k8s.io/yaml"
	"log"
)


//获得folder文件夹下所有yaml文件
func ListYamlFlie(folder string)error{
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") {
			// conf, serviceName, err := parser.ParseYaml(path)
			if err != nil {
				fmt.Printf("Error parsing YAML file %s: %v\n", path, err)
				return nil // Continue processing other files even if this one fails
			}
			fmt.Println(path)

			// Assuming you want to track YAML files that successfully parsed and contained the app label
			// if serviceName != "" {
			// 	parsedYamls[serviceName] = conf
			// 	applicationFolders[serviceName] = filepath.Dir(path)
			// 	validYamlFiles = append(validYamlFiles, path)
			// }
		}
		return nil
	})
	if err != nil{
		return err
	}
	return nil
}

//读取yaml文件
func readJavaYamlFile(fileFullPath string)(error){
	result := make(map[string]interface{})
	yamlFile, err := os.ReadFile(fileFullPath)
	if err != nil {
		log.Println("fail to read ", fileFullPath, err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, &result)
	
	return nil
}

func readK8sYamlFile(fileFullPath string)(error){
	result := make(map[string]interface{})
	yamlFile, err := os.ReadFile(fileFullPath)
	if err != nil {
		log.Println("fail to read ", fileFullPath, err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, &result)
	
	return nil
}
