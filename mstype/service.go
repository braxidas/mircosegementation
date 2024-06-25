package mstype

//描述K8s服务
type K8sService struct {
	FilePath             string //jar包所在路径
	DeploymentPath       string //部署文件所在路径
	NacosApplicationFile string //注册中心配置文件所在路径

	ApplicationName string // 向注册中心注册服务名
	PodName         string //最终部署到k8s中的微服务名称
	ServiceName     string

	Ingress []*K8sService //networkpolicy的ingress列表
	Egress  []*K8sService //networkpolicy的egress列表

	JavaClassList []JavaClass //soot扫描到的有效信息

	//spring cloud逻辑
	Consume       []string //消费的服务
	JavaInterface []string //使用的interface(interface使用的服务也要记录)
	//spring dubbo逻辑
	DubboReference []string //提供的dubbo的service
	DubboService   []string //消费的dubbo的service
}

//描述java中的service接口
type JavaClass struct {
	ClassName      string   `json:"className"`      //类名称
	Consume        []string `json:"consume"`        //restTemplate或者openFeign显示调用的微服务
	Field          []string `json:"field"`          //使用自动注入调用其他java的interface的集合
	DubboReference []string `json:"dubboReference"` //提供的dubbo的service
	DubboService   []string `json:"dubboService"`   //消费的dubbo的service
	DefineAspect   []string `json:"defineAspect"`   //表示定义的aspect的集合  //注意aspect获得类使用"."表示层次 而一般的类都使用"\"表示层次
	UseAspect      []string `json:"useAspect"`      //表示类方法使用的注解集合，可能与aspect对应
}

func (k8sService *K8sService) AppendIngress(ingress *K8sService) {
	k8sService.Ingress = append(k8sService.Ingress, ingress)
}

func (k8sService *K8sService) AppendEgress(egress *K8sService) {
	k8sService.Egress = append(k8sService.Egress, egress)
}

func (k8sService *K8sService) ProvideService(dubboReference string) bool {
	for _, v := range k8sService.DubboService {
		if dubboReference == v {
			return true
		}
	}
	return false
}
