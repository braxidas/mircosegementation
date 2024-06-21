package mstype

//描述K8s服务
type K8sService struct {
	FilePath        string  //jar包所在路径
	DeploymentPath  string//部署文件所在路径
	ApplicationName string// 向注册中心注册服务名
	PodName         string//最终部署到k8s中的微服务名称
	ServiceName     string

	Ingress []*K8sService//networkpolicy的ingress列表
	Egress  []*K8sService//networkpolicy的egress列表

	Consume        []string//消费的服务
	JavaInterface  []string//使用的interface(interface使用的服务也要记录)
	DubboReference []string//提供的dubbo的service
	DubboService   []string//消费的dubbo的service
}

//描述java中的service接口
// type JavaService struct {
// 	Name   string
// 	Caller string
// }

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
