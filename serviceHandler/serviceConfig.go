package serviceHandler

import (
	"fmt"
	"microsegement/fileHandler"
	"microsegement/mstype"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// 根据k8sList生成networkPolicy类
func GenerateIpPolicy(k8sServiceList []*mstype.K8sService) []*mstype.NetworkPolicy {
	var result []*mstype.NetworkPolicy
	finalK8sServiceList := handleConfig(k8sServiceList) //获得最终部署pod的列表
	os.RemoveAll("output")
	os.Mkdir("output",0766)
	for _, v := range finalK8sServiceList {
		networkPolicy := new(mstype.NetworkPolicy)
		networkPolicy.ApiVerson = "networking.k8s.io/v1"
		networkPolicy.Kind = "NetworkPolicy"
		networkPolicy.Metadata.Name = v.PodName + "-policy"
		networkPolicy.Metadata.Namespace = v.Namespace
		networkPolicy.Spec.Egress = v.EgressOut
		networkPolicy.Spec.Ingress = v.IngressOut
		networkPolicy.Spec.PodSelector.MatchLabels = v.Labels
		networkPolicy.Spec.PolicyTypes = []string{"Egress", "Ingress"}
		result = append(result, networkPolicy)
		fileHandler.WriteToYaml(networkPolicy) //写成yaml文件
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
func handleConfig(k8sServiceList []*mstype.K8sService) []*mstype.K8sService {
	var finalK8sServiceList []*mstype.K8sService
	for _, v := range k8sServiceList {
		if v.PodName == "" { //如果没有podname 则说明不会被部署，因此不加入最终的config分析
			continue
		}
		for _, va := range v.ApplicationList {
			egress,egress2 := handleEgress(va)
			v.EgressOut = append(v.EgressOut, egress...)
			v.MergeIgress(egress2)
			ingress := handleIngress(va)
			v.IngressOut = append(v.IngressOut, ingress...)
		}
		v.Egress2EgressOut()
		finalK8sServiceList = append(finalK8sServiceList, v)
	}
	return finalK8sServiceList
}

func handleEgress(application *mstype.Application) ([]*mstype.Policy,map[*mstype.K8sService]struct{}) {
	var egress []*mstype.Policy
	egress2 := make(map[*mstype.K8sService]struct{})
	//Nacos
	if application.Spring.Cloud.Nacos.Discovery.ServerAddr != "" && application.Spring.Cloud.Nacos.Discovery.ServerAddr != "localhost" {
		addr := strings.Split(application.Spring.Cloud.Nacos.Discovery.ServerAddr, ":")[0]
		if v, ok := svc2Pod[addr]; ok {
			egress = append(egress, mstype.NewPodPolicy(getLabel(v)))
		} else {
			if (len(strings.Split(application.Spring.Cloud.Nacos.Discovery.ServerAddr, ":")) > 1){
				port, _ := strconv.Atoi((strings.Split(application.Spring.Cloud.Nacos.Discovery.ServerAddr, ":"))[1])
				egress = append(egress, []*mstype.Policy{mstype.NewEgress(port, addr),mstype.NewEgress(port+1000, addr),
					mstype.NewEgress(port+1001, addr),mstype.NewEgress(port-1000, addr)}...)
			}else{
				egress = append(egress, []*mstype.Policy{mstype.NewEgress(8848, addr), mstype.NewEgress(9848, addr)}...)
			}
		}
	} 
	if application.Spring.Cloud.Nacos.Config.ServerAddr != "" && application.Spring.Cloud.Nacos.Config.ServerAddr != "localhost" {
		addr := strings.Split(application.Spring.Cloud.Nacos.Config.ServerAddr, ":")[0]
		if v, ok := svc2Pod[addr]; ok {
			egress = append(egress, mstype.NewPodPolicy(getLabel(v)))
		} else {
			if (len(strings.Split(application.Spring.Cloud.Nacos.Config.ServerAddr, ":")) > 1){
				port, _ := strconv.Atoi((strings.Split(application.Spring.Cloud.Nacos.Config.ServerAddr, ":"))[1])
				egress = append(egress, []*mstype.Policy{mstype.NewEgress(port, addr),mstype.NewEgress(port+1000, addr),
					mstype.NewEgress(port+1001, addr),mstype.NewEgress(port-1000, addr)}...)
			}else{
				egress = append(egress, []*mstype.Policy{mstype.NewEgress(8848, addr), mstype.NewEgress(9848, addr)}...)
			}
		}
	}
	//Redis
	if application.Spring.Redis.Host != "" && application.Spring.Redis.Host != "localhost"{
		port, _ := strconv.Atoi(application.Spring.Redis.Port)
		if v, ok := svc2Pod[application.Spring.Redis.Host]; ok {
			egress = append(egress, mstype.NewPodPolicy(getLabel(v)))
		} else {
			egress = append(egress, mstype.NewEgress(port, application.Spring.Redis.Host))
		}
	}
	//Database
	if application.Spring.DataSource.Url != "" {
		host, portStr := getHostPort(strings.ReplaceAll(application.Spring.DataSource.Url, "jdbc:", ""))
		if host != "localhost"&& host != ""{
			if v, ok := svc2Pod[host]; ok {
				egress = append(egress, mstype.NewPodPolicy(getLabel(v)))
			} else {
				port, _ := strconv.Atoi(portStr)
				egress = append(egress, mstype.NewEgress(port, host))
			}
		}
	}
	if application.Spring.DataSource.Dynamic.DataSource.Master.Url != "" {
		host, portStr := getHostPort(strings.ReplaceAll(application.Spring.DataSource.Dynamic.DataSource.Master.Url, "jdbc:", ""))
		if host != "localhost"&& host != ""{
			if v, ok := svc2Pod[host]; ok {
				egress = append(egress, mstype.NewPodPolicy(getLabel(v)))
			} else {
				port, _ := strconv.Atoi(portStr)
				egress = append(egress, mstype.NewEgress(port, host))
			}
		}
	}
	// Minio
	if application.Minio.Url != "" {
		host, portStr := getHostPort(application.Minio.Url)
		if host != "localhost"&& host != ""{
			if v, ok := svc2Pod[host]; ok {
				egress = append(egress, mstype.NewPodPolicy(getLabel(v)))
			} else {
				port, _ := strconv.Atoi(portStr)
				egress = append(egress, mstype.NewEgress(port, host))
			}
		}
	}
	//Fdfs
	if application.Fdfs.Domain != "" {
		host, _ := getHostPort(application.Fdfs.Domain)
		if host != "localhost" && host != ""{
			if v, ok := svc2Pod[host]; ok {
				egress = append(egress, mstype.NewPodPolicy(getLabel(v)))
			} else {
				portStr := strings.Split(application.Fdfs.TrackerList, ":")[1]
				port, _ := strconv.Atoi(portStr)
				egress = append(egress, mstype.NewEgress(port, host))
			}
		}
	}
	//spring.cloud.gateway.routes
	if len(application.Spring.Cloud.Gateway.Routes) > 0 {
		for _, v := range application.Spring.Cloud.Gateway.Routes {
			if vs, ok := name2K8sService[getUriName(v.Uri)]; ok {
				egress2[vs] = struct{}{}// append(egress, mstype.NewPodPolicy(vs.Labels))
			}
		}
	}
	return egress, egress2
}
func handleIngress(application *mstype.Application) []*mstype.Policy {
	var ingress []*mstype.Policy
	if application.Server.Port != "" {
		port, _ := strconv.Atoi(application.Server.Port)
		ingress = append(ingress, mstype.NewIngress(port))
	}
	return ingress
}

// 从url中获得端口和ip
func getHostPort(urlstr string) (string, string) {
	// urlstr = strings.ReplaceAll(urlstr,"jdbc:","")
	u, err := url.Parse(urlstr)
	if err != nil {
		fmt.Println("wrong parser", urlstr, err.Error())
		return "", ""
	}
	// if v, ok := svc2Pod[u.Hostname()]; ok == true {
	// 	return v, u.Port()
	// }
	return u.Hostname(), u.Port()
}

// 从uri中获得服务名称
func getUriName(uri string) string {
	return strings.ReplaceAll(uri, "lb://", "")
}

// 从地址 比如"alertmanager=main, app=alertmanager"中解析出labels
func getLabel(addr string) map[string]string {
	mp := make(map[string]string)
	str := strings.Split(addr, ",")[0]
	mp[strings.Split(str, "=")[0]] = strings.Split(str, "=")[1]
	return mp
}
