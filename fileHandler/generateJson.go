package fileHandler

import (
	"encoding/json"
	"fmt"
	"microsegement/mstype"
	"os"
	"path/filepath"
)

//生成json文件
func WriteToJson(k8sService *mstype.K8sService) error{
	
	manifest := ServicetoManifest(k8sService)

	jsonData, err := json.MarshalIndent(manifest, "", " ")

	if k8sService.PodName == ""{
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to marshal TCPManifest for service '%s': %w", k8sService.PodName, err)
	}

	filename := filepath.Join("output" , manifest.Service + ".json")
	err = os.WriteFile(filename, jsonData, 0777)

	if err != nil {
		return fmt.Errorf("failed to write TCPManifest to file '%s': %w", filename, err)
	}

	return nil
}

func ServicetoManifest(k8sService *mstype.K8sService) mstype.TCPManifest{
	var manifest mstype.TCPManifest
	manifest.Service = k8sService.PodName
	manifest.Version = "v1"
	for k, _ := range(k8sService.Egress){
		var request mstype.TCPRequest
		request.Name = k.PodName
		request.Type = "Tcp"
		manifest.Requests = append(manifest.Requests, request)
	}
	return manifest
}