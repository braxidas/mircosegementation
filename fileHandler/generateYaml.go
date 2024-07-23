package fileHandler

import (
	"fmt"
	"microsegement/mstype"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

//将networkPolicy写成文件
func WriteToYaml(networkPolicy *mstype.NetworkPolicy) error {

	yamlData, err := yaml.Marshal(networkPolicy)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("fail to marshal " + networkPolicy.Metadata.Name)
	}

	filename := filepath.Join("output", networkPolicy.Metadata.Name+".yaml")
	err = os.WriteFile(filename, yamlData, 0777)

	if err != nil {
		return fmt.Errorf("failed to write networkpolicy to file '%s': %v", filename, err)
	}
	return nil
	
}