package mstype

import "fmt"

/*
此类用于描述一个微服务
*/

// 描述K8s服务
type K8sService struct {
	FilePath             string //jar包所在路径
	DeploymentPath       string //部署文件所在路径
	NacosApplicationFile string //注册中心配置文件所在路径

	ApplicationName string            //向注册中心注册服务名
	PodName         string            //最终部署到k8s中的微服务名称
	Namespace       string            //Namespace
	ApiVersion      string            //apiVersion
	Labels          map[string]string //匹配选择器
	ServiceName     string

	Ingress map[*K8sService]struct{} //networkpolicy的ingress集合 针对nacos中的微服务
	Egress  map[*K8sService]struct{} //networkpolicy的egress列表	针对nacos中的微服务

	// Ingress map[string]struct{} //networkpolicy的ingress集合 针对nacos中的微服务
	// Egress  map[string]struct{} //networkpolicy的egress列表	针对nacos中的微服务

	IngressOut []*Policy //networkpolicy的ingress集合	包括外部组件的ip
	EgressOut  []*Policy //networkpolicy的egress列表	包括外部组件的ip

	JavaClassList    []*JavaClass //soot扫描到的该服务中直接定义的类
	JavaClassAllList []string     //图分析后该服务中实际应用的全部类
	//用string是因为实际应用的类不一定是自定义类，如redisService，为了之后可能有非自定义类需要分析，所以保留

	ApplicationList []*Application //用于表示k8sService的配置文件列表
}

// 描述java中的类接口
type JavaClass struct {
	K8sService     *K8sService //指向该类所属的模块
	ClassName      string      `json:"className"`      //类名称
	Consume        []string    `json:"consume"`        //restTemplate或者openFeign显示调用的微服务
	Field          []string    `json:"field"`          //使用自动注入调用其他java的interface的集合
	DubboReference []string    `json:"dubboReference"` //提供的dubbo的service
	DubboService   []string    `json:"dubboService"`   //消费的dubbo的service
	DefineAspect   []string    `json:"defineAspect"`   //表示定义的aspect的集合
	UseAspect      []string    `json:"useAspect"`      //表示类方法使用的注解集合，可能与aspect对应
}

func (k8sService *K8sService) AppendIngress(ingress *K8sService) {
	k8sService.Ingress[ingress] = struct{}{}
}



func (k8sService *K8sService) AppendEgress(egress *K8sService) {
	k8sService.Egress[egress] = struct{}{}
}
//合并map
func (k8sService *K8sService) MergeIgress(egress map[*K8sService]struct{}){
	for k, _ := range egress{
		k8sService.Egress[k] = struct{}{}
	}
} 

// 判断一个类是否调用了该微服务
func (k8sService *K8sService) ProvideService(dubboReference string) bool {
	for _, v := range k8sService.JavaClassList {
		for _, vc := range v.DubboService {
			if vc == dubboReference {
				return true
			}
		}
	}
	return false
}


func (k8sService *K8sService) Egress2EgressOut(){
	for k := range k8sService.Egress {
		k8sService.EgressOut = append(k8sService.EgressOut, NewPodPolicy(k.Labels))
	}
}

// 判断一个类是否调用了该类的dubbo的rpc
func (javaClass *JavaClass) ProvideDubbo(dubboReference string) bool {
	for _, vc := range javaClass.DubboService {
		if vc == dubboReference {
			return true
		}
	}
	return false
}

func (javaClass *JavaClass) PrintInfo() {
	fmt.Println("------ClassName:", javaClass.ClassName, "------")
	fmt.Println("Consume:", javaClass.Consume)
	fmt.Println("Field:", javaClass.Field)
	fmt.Println("DubboReference:", javaClass.DubboReference)
	fmt.Println("DubboService :", javaClass.DubboService)
	fmt.Println("DefineAspect:", javaClass.DefineAspect)
	fmt.Println("UseAspect :", javaClass.UseAspect)
}
