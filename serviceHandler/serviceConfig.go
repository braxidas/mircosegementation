package serviceHandler

import (
	"fmt"
	"microsegement/fileHandler"
	"microsegement/mstype"
	"net/url"
	"strconv"
	"strings"
)

//根据k8sList生成networkPolicy类
func GenerateIpPolicy(k8sServiceList []*mstype.K8sService) []*mstype.NetworkPolicy{
	var result []*mstype.NetworkPolicy
	finalK8sServiceList := handleConfig(k8sServiceList)//获得最终部署pod的列表
	for _, v := range finalK8sServiceList {
		networkPolicy := new(mstype.NetworkPolicy)
		networkPolicy.ApiVerson = "networking.k8s.io/v1"
		networkPolicy.Kind = "NetworkPolicy"
		networkPolicy.Metadata.Name = v.PodName + "-policy"
		networkPolicy.Metadata.Namespace = "default"
		networkPolicy.Spec.Egress = v.EgressOut
		networkPolicy.Spec.Ingress = v.IngressOut
		networkPolicy.Spec.PodSelector.MatchLabels.App = v.PodName
		networkPolicy.Spec.PolicyTypes = []string{"Egress", "Igress"}
		result = append(result, networkPolicy)
		fileHandler.WriteToYaml(networkPolicy)//写成yaml文件
	}
	return result
}

//根据k8sService的applicaitonList中存储的配置，获得相应的进出policy
//获得配置yaml文件中的ipblock信息
/*
Egress
	nacos注册中心地址
	nacos服务器上export的配置文件中的ip地址
Ingress
	jar包中的server.port表示对外暴露的端口

*/
func handleConfig(k8sServiceList []*mstype.K8sService)[]*mstype.K8sService{
	var finalK8sServiceList []*mstype.K8sService
	for _, v := range k8sServiceList{
		if v.PodName == ""{//如果没有podname 则说明不会被部署，因此不加入最终的config分析
			continue
		}
		for _, va := range v.ApplicationList{
			egress := handleEgress(va)
			v.EgressOut = append(v.EgressOut, egress...)
			ingress := handleIngress(va)
			v.IngressOut = append(v.IngressOut, ingress...)
		}
		finalK8sServiceList = append(finalK8sServiceList, v)
	}
	return finalK8sServiceList
}

func handleEgress(application *mstype.Application)[]*mstype.Policy{
	var egress []*mstype.Policy
	//Nacos
	if application.Spring.Cloud.Nacos.Discovery.ServerAddr != ""{
		addr := strings.Split(application.Spring.Cloud.Nacos.Discovery.ServerAddr,":")[0]
		egress = append(egress, []*mstype.Policy{mstype.NewEgress(8848,addr), mstype.NewEgress(9848,addr)}...)
	}else if application.Spring.Cloud.Nacos.Config.ServerAddr != ""{
		addr := strings.Split(application.Spring.Cloud.Nacos.Config.ServerAddr,":")[0]
		egress = append(egress, []*mstype.Policy{mstype.NewEgress(8848,addr), mstype.NewEgress(9848,addr)}...)
	}
	//Redis
	if application.Spring.Redis.Host != ""{
		port, _ := strconv.Atoi(application.Spring.Redis.Port)
		egress = append(egress,mstype.NewEgress(port,application.Spring.Redis.Host))
	}
	//Database
	if application.Spring.DataSource.Url != ""{
		host, portStr := getHostPort(application.Spring.DataSource.Url)
		port, _ := strconv.Atoi(portStr)
		egress = append(egress, mstype.NewEgress(port, host))
	}
	if application.Spring.DataSource.Dynamic.DataSource.Master.Url != ""{
		host, portStr := getHostPort(application.Spring.DataSource.Dynamic.DataSource.Master.Url)
		port, _ := strconv.Atoi(portStr)
		egress = append(egress, mstype.NewEgress(port, host))
	}
	//
	return egress
}
func handleIngress(application *mstype.Application)[]*mstype.Policy{
	var ingress []*mstype.Policy
	if application.Server.Port != ""{
		port, _ := strconv.Atoi(application.Server.Port)
		ingress = append(ingress, mstype.NewIngress(port))
	}
	return ingress
}

func getHostPort(urlstr string)(string,string){
	urlstr = strings.ReplaceAll(urlstr,"jdbc:","")
	u, err := url.Parse(urlstr)
	if err != nil{
		fmt.Println("wrong parser", urlstr, err.Error())
	}
	return u.Hostname(), u.Port()
}
