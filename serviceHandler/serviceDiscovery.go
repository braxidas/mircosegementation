package serviceHandler

import (
	"fmt"
	"microsegement/fileHandler"
	"microsegement/mstype"
	"microsegement/soot"
	"sort"
	"strings"
)

var interface2service map[string]string

func DiscoverService(k8sServiceList []*mstype.K8sService) error {

	interface2service = make(map[string]string)

	for i, _ := range k8sServiceList {
		strList, err := soot.ScanDiscoverService(k8sServiceList[i].FilePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(strList) == 1 {
			continue
		}
		analysisStdOut(strList, k8sServiceList[i])
	}

	analysisCall(k8sServiceList)

	for _, v := range k8sServiceList {
		err := fileHandler.WriteToJson(v)
		if err != nil {
			return err
		}
	}

	return nil
}

// maybe optimize
func analysisStdOut(strList []string, k8sService *mstype.K8sService) {
	for _, v := range strList {
		if strings.HasPrefix(v, "consume:") {
			k8sService.Consume = append(k8sService.Consume, (strings.Split(v, ":"))[1])
		} else if strings.HasPrefix(v, "consume-rpc:") {
			k8sService.DubboReference = append(k8sService.DubboReference, (strings.Split(v, ":"))[1])
		} else if strings.HasPrefix(v, "consume-rpc:") {
			k8sService.DubboService = append(k8sService.DubboService, (strings.Split(v, ":"))[1])
		} else if strings.HasPrefix(v, "interface") {
			k8sService.JavaInterface = append(k8sService.JavaInterface, (strings.Split(v, ":"))[1])
		} else if strings.HasPrefix(v, "map") {
			interface2service[(strings.Split(v, ":"))[1]] = (strings.Split(v, ":"))[2]
		}
	}
	k8sService.Consume = removeDuplicates(k8sService.Consume)
	k8sService.DubboReference = removeDuplicates(k8sService.DubboReference)
	k8sService.DubboService = removeDuplicates(k8sService.DubboService)
	k8sService.JavaInterface = removeDuplicates(k8sService.JavaInterface)
}

// maybe optimize
func analysisCall(k8sServiceList []*mstype.K8sService) {
	for i, _ := range k8sServiceList {
		//处理openfeign和restTemplate
		for _, vc := range k8sServiceList[i].Consume {
			for j, _ := range k8sServiceList {
				if vc == k8sServiceList[j].ApplicationName {
					k8sServiceList[i].AppendEgress(k8sServiceList[j])
					k8sServiceList[j].AppendIngress(k8sServiceList[i])
					break
				}
			}
		}
		//处理Dubbo
		for _, vdc := range k8sServiceList[i].DubboReference {
			for j, _ := range k8sServiceList {
				if k8sServiceList[j].ProvideService(vdc) {
					k8sServiceList[i].AppendEgress(k8sServiceList[j])
					k8sServiceList[j].AppendIngress(k8sServiceList[i])
					break
				}
			}
		}
		//处理interface
		for _, vi := range k8sServiceList[i].JavaInterface {
			value, ok := interface2service[vi]
			if ok {
				for j, _ := range k8sServiceList {
					if value == k8sServiceList[j].ApplicationName {
						k8sServiceList[i].AppendEgress(k8sServiceList[j])
						k8sServiceList[j].AppendIngress(k8sServiceList[i])
						break
					}				
				}
			}
		}
	}
}

func removeDuplicates(elements []string) []string {
	if len(elements) < 2 {
		return elements
	}
	sort.Strings(elements) // 先对字符串数组进行排序
	j := 0
	for i := 1; i < len(elements); i++ {
		if elements[i] != elements[j] {
			j++
			elements[j] = elements[i]
		}
	}
	return elements[:j+1]
}
