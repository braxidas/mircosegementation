package fileHandler

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"

	// "log"
	"microsegement/sql"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	tempDependency string //"tempDependency"
	tempJarfile    string //"tempJar"
)

/*
将所有需要分析的jar包解压到tempJarfile下
并对每个解压的jar包中的依赖包进行筛选
将其中的非第三方库拷贝到tempDependency以供后续分析
*/

// 解压jar包到指定文件夹下,并获得所有非第三方库
func GetNotThirdJar(folder string) {
	sql.InitDB()
	tempDependency = filepath.Join(folder, "tempDependency")
	tempJarfile = "tempJar"
	os.RemoveAll(tempDependency)
	os.RemoveAll(tempJarfile)
	os.Mkdir(tempDependency, 0755)
	os.Mkdir(tempJarfile, 0755)
	// defer os.RemoveAll(tempJarfile)

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".jar") {
			if err = cpNotThirdJar(path); err != nil {
				fmt.Println(err)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

// 把jar包解压到指定文件夹下
func cpNotThirdJar(jarFile string) error {
	targetJarFile := filepath.Join(tempJarfile, getJarName(jarFile))
	// exec.Command("ls", targetJarFile)
	cmd := exec.Command("unzip", jarFile, "-d", targetJarFile)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to unzip file %s because %v", jarFile, err)
	}
	filterThirdArtifact(filepath.Join(targetJarFile, "BOOT-INF", "lib"))
	return nil
}

// 根据路径获得扫描该路径下的所有jar包并且将非第三方库jar包拷贝的到可供分析的文件夹
func filterThirdArtifact(folder string) {
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".jar") {
			if isThird := isThirdArtifact(path); isThird == false {
				return copyDependency(path, tempDependency)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

// 是否是第三方库
func isThirdArtifact(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	hash := sha1.New()
	_, _ = io.Copy(hash, file)
	sha1 := hex.EncodeToString(hash.Sum(nil))
	sha1b, _ := hex.DecodeString(sha1)
	res, _ := sql.QueryBySha1(sha1b)
	// res, err := sql.QueryBySha1(sha1b)
	// if err == nil{
	// 	log.Println(path,"is third lib")
	// }
	return res
}

// 拷贝目标文件到指定文件夹下
func copyDependency(path string, targetPath string) error {
	cmd := exec.Command("cp", path, targetPath)
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to cp file %s : %v", path, err)
	}
	// fmt.Printf("succeed to cp file %s\n", path)
	return nil
}

// 获得一个jar文件的名称
func getJarName(path string) string {
	separator := string(os.PathSeparator)
	return path[strings.LastIndex(path, separator):strings.LastIndex(path, ".")]
}
