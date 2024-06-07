package filehandler

import (
	// jar "pkg.re/essentialkaos/go-jar.v1"
	"archive/zip"
	"fmt"
	"io"
	"microsegement/mstype"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// 获得folder文件夹下所有jar文件
func ListJarFile(folder string) error {

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".jar") {
			fmt.Println(path)
			application, err := getJarYamlFile(path)
			if err != nil {
				fmt.Printf("Error read Jar YAML file %s: %v\n", path, err)
				return nil
			}
			fmt.Println(*application)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// 获得指定jar包中的application.yaml文件
func getJarYamlFile(jarFile string) (*mstype.Application, error) {
	application := new(mstype.Application)

	r, err := zip.OpenReader(jarFile)
	if err != nil {
		return application, err
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".yml") {
			rc, err := f.Open()
			if err != nil {
				return application, err
			}
			defer rc.Close()

			yamlFile, err := io.ReadAll(rc)
			if err != nil {
				return application, err
			}

			err = yaml.Unmarshal(yamlFile, application)
			return application, err
		}
	}
	return application, nil
}

// func TestYaml() {
// 	application := new(mstype.Application)
// 	mp := make(map[string]interface{})
// 	yamlFile, _ := os.ReadFile(`target\application.yml`)
// 	err := yaml.Unmarshal(yamlFile, &mp)
// 	fmt.Println(mp)
// 	err = yaml.Unmarshal(yamlFile, application)
// 	fmt.Println(err)
// 	fmt.Println(application)
// }
