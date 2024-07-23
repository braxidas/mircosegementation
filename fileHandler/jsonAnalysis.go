package fileHandler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

//根据svc.json获得svc到pod的映射
func GetSvc2Pod(myfolder string)map[string]string {
	res := make(map[string]string)
	jsonPath := filepath.Join(myfolder,"svc.json")
	file, err := os.ReadFile(jsonPath)
	if err != nil{
		fmt.Println("读取文件失败：", err)
		return res
	}
	err = json.Unmarshal(file, &res)
	if err != nil {
        fmt.Println("解析JSON失败:", err)
    }
	return res	
}
//获得环境变量
func getEnv(myfolder string)map[string]string {
	res := make(map[string]string)
	jsonPath := filepath.Join(myfolder,"env.json")
	file, err := os.ReadFile(jsonPath)
	if err != nil{
		fmt.Println("读取文件失败：", err)
		return res
	}
	err = json.Unmarshal(file, &res)
	if err != nil {
        fmt.Println("解析JSON失败:", err)
    }
	return res	
}