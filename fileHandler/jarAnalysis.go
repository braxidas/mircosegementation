package fileHandler

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
func ListJarFile(folder string) ([]*mstype.Application,[]string, error) {
	var (
		applicationList []*mstype.Application
		pathList []string
	)

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".jar") {
			application, err := getJarYamlFile(path)
			if err != nil {
				return err
			}

			// fmt.Println("\\\\\\\\\\\\\\\\"+getParentDirectory(path))
			pathList = append(pathList, path)
			applicationList = append(applicationList, application)

		}
		return nil
	})
	if err != nil{
		fmt.Println(err)
	}
	return applicationList,pathList,nil
}

// 获得指定jar包中的a配置yaml文件
func getJarYamlFile(jarFile string) (*mstype.Application, error) {
	application := new(mstype.Application)

	r, err := zip.OpenReader(jarFile)
	if err != nil {
		fmt.Errorf("fail to open", jarFile, err)
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
			if err != nil{
				return application, err
			}

			_, err = application.GetApplicationName()
			if err == nil{
				return application, nil
			}
		}
	}
	// fmt.Println("fail to find application name when scan yaml" , jarFile)
	return application, nil
}

//获得jar包上级目录
func getParentDirectory(path string) string {
	separator := string(os.PathSeparator)
	return path[0:strings.LastIndex(path, separator)]
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
