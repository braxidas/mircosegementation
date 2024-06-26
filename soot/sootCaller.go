package soot

import (
	"fmt"
	"microsegement/mstype"
	"os/exec"
	// "strings"
	// "path/filepath"
	"encoding/json"
)

func ScanDiscoverService(applicationName string)([]*mstype.JavaClass, error){
	// applicationName = `C:\Users\li_sh\Desktop\msexample\example\0ruoyiExample\rouyi-system-api\ruoyi-api-system-3.6.4.jar`
	//  applicationName = `C:\Users\li_sh\Desktop\msexample\example\0ruoyiExample\ruoyi-common\ruoyi-common-log\ruoyi-common-log-3.6.4.jar`
	// var res []string
	cmd := exec.Command("java", "-jar", "soot-analysis-1.0-SNAPSHOT.jar", applicationName)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to call soot: %v\n", err)
	}


	var javaClassList []*mstype.JavaClass	
	// fmt.Println("ssssssssssssssssssssss",out)
	err = json.Unmarshal(out, &javaClassList)
	
	// for _, v := range(javaClassList){
	// 	v.PrintInfo()
	// }
	return javaClassList, nil
}


