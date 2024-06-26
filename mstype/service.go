package mstype

import "fmt"

//描述K8s服务
type K8sService struct {
	FilePath             string //jar包所在路径
	DeploymentPath       string //部署文件所在路径
	NacosApplicationFile string //注册中心配置文件所在路径

	ApplicationName string // 向注册中心注册服务名
	PodName         string //最终部署到k8s中的微服务名称
	ServiceName     string

	Ingress map[*K8sService]struct{} //networkpolicy的ingress集合
	Egress  map[*K8sService]struct{} //networkpolicy的egress列表

	JavaClassList    []*JavaClass //soot扫描到的该服务中直接定义的类
	JavaClassAllList []string     //图分析后该服务中实际应用的全部类

	// //spring cloud逻辑
	// Consume       []string //消费的服务
	// JavaInterface []string //使用的interface(interface使用的服务也要记录)
	// //spring dubbo逻辑
	// DubboReference []string //提供的dubbo的service
	// DubboService   []string //消费的dubbo的service
}

//描述java中的类接口
type JavaClass struct {
	ClassName      string   `json:"className"`      //类名称
	Consume        []string `json:"consume"`        //restTemplate或者openFeign显示调用的微服务
	Field          []string `json:"field"`          //使用自动注入调用其他java的interface的集合
	DubboReference []string `json:"dubboReference"` //提供的dubbo的service
	DubboService   []string `json:"dubboService"`   //消费的dubbo的service
	DefineAspect   []string `json:"defineAspect"`   //表示定义的aspect的集合
	UseAspect      []string `json:"useAspect"`      //表示类方法使用的注解集合，可能与aspect对应
}

func (k8sService *K8sService) AppendIngress(ingress *K8sService) {
	k8sService.Ingress[ingress] = struct{}{}
}

func (k8sService *K8sService) AppendEgress(egress *K8sService) {
	k8sService.Egress[egress] = struct{}{}
}

//判断一个类是否调用了该微服务
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

func (javaClass *JavaClass) PrintInfo() {
	fmt.Println("------ClassName:", javaClass.ClassName, "------")
	fmt.Println("Consume:", javaClass.Consume)
	fmt.Println("Field:", javaClass.Field)
	fmt.Println("DubboReference:", javaClass.DubboReference)
	fmt.Println("DubboService :", javaClass.DubboService)
	fmt.Println("DefineAspect:", javaClass.DefineAspect)
	fmt.Println("UseAspect :", javaClass.UseAspect)
}
