package mstype

//描述K8s服务
type K8sService struct {
	FilePath        string
	ApplicationName string
	PodName         string
	ServiceName     string
	Ingress         []*K8sService
	Egress          []*K8sService
}

func (k8sService *K8sService) AppendIngress(ingress *K8sService){
	k8sService.Ingress = append(k8sService.Ingress, ingress)
}

func (k8sService *K8sService) AppendEgress(egress *K8sService){
	k8sService.Ingress = append(k8sService.Egress, egress)
}