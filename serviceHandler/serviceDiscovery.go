package serviceHandler

import (
	"fmt"
	"microsegement/mstype"
	"microsegement/soot"
	"strings"
)

func DiscoverService(k8sServiceList []*mstype.K8sService)error{
	for i, _:= range(k8sServiceList){		
		strList, err := soot.DiscoverService(k8sServiceList[i].ApplicationName)
		if err != nil{
			fmt.Println(err)
			continue
		}
		if(len(strList) == 1){
			continue
		}
		analysisStdOut(strList, k8sServiceList[i])
	}

	analysisCall(k8sServiceList)

	return nil
}

func analysisStdOut(strList []string, k8sService *mstype.K8sService){
	for _, v := range(strList){
		if strings.HasPrefix(v,"consume:"){
			k8sService.Consume = append(k8sService.Consume, (strings.Split(v, ":"))[1])
		}else if strings.HasPrefix(v,"consume-rpc:"){
			k8sService.DubboReference = append(k8sService.DubboReference, (strings.Split(v, ":"))[1])
		}else if strings.HasPrefix(v,"consume-rpc:"){
			k8sService.DubboService = append(k8sService.DubboService, (strings.Split(v, ":"))[1])
		}
	}
}

func analysisCall(k8sServiceList []*mstype.K8sService){
	for i, _ := range(k8sServiceList){
		//处理openfeign和restTemplate
		for _, vc := range(k8sServiceList[i].Consume){
			for j, _ := range(k8sServiceList){
				if vc == k8sServiceList[j].ApplicationName{
					k8sServiceList[i].AppendEgress(k8sServiceList[j])
					k8sServiceList[j].AppendIngress(k8sServiceList[i])
					break
				}
			}
		}
		//处理Dubbo
		for _, vdc := range(k8sServiceList[i].DubboReference){
			for j, _ := range(k8sServiceList){
				if k8sServiceList[j].ProvideService(vdc){
					k8sServiceList[i].AppendEgress(k8sServiceList[j])
					k8sServiceList[j].AppendIngress(k8sServiceList[i])
					break
				}
			}
		}
	}
}

