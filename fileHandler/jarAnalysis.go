package fileHandler

import (
	// jar "pkg.re/essentialkaos/go-jar.v1"
	"archive/zip"
	"fmt"
	"io"
	"microsegement/mstype"
	"os"
	"path/filepath"
	"regexp"
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

// 获得指定jar包中的application配置yaml文件
func getJarYamlFile(jarFile string) (*mstype.Application, error) {
	application := new(mstype.Application)

	env := getEnv(getParentDirectory(jarFile))

	r, err := zip.OpenReader(jarFile)
	if err != nil {
		return application, fmt.Errorf("fail to open %s, because %v", jarFile, err)
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".yml") || strings.HasSuffix(f.Name, ".yaml") {
			rc, err := f.Open()
			if err != nil {
				return application, err
			}
			defer rc.Close()

			yamlFile, err := io.ReadAll(rc)
			if err != nil {
				return application, err
			}
			yamlFile, err = handleEnv(yamlFile, env)
			if err != nil {
				return application, err
			}

			err = yaml.Unmarshal(yamlFile, application)
			if err != nil{
				return application, err
			}

			if err != nil{
				return application, nil
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

func handleEnv(yamlFile []byte, env map[string]string) ([]byte, error){
	yamlContent := string(yamlFile)
	re := regexp.MustCompile(`\$\{[\w]+\}`)
	envs := re.FindAllString(yamlContent, -1)
	for _, v := range envs{
		temp := v[2: len(v) - 1]
		if val, ok := env[temp]; ok{//如果环境变量中含有该值，则用环境变量对应的值替换
			strings.ReplaceAll(yamlContent, v, val)
		}else if strings.Contains(temp, ":"){//如果环境变量中不含有该值，则用：后的默认值替换
			strings.ReplaceAll(yamlContent, v, temp[strings.LastIndex(temp, ":") + 1:])
		}
	}

	return []byte(yamlContent),nil
}

