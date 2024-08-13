package mstype

/*
此类用于生成目标文件中的微隔离策略
*/

// 策略
type NetworkPolicy struct {
	ApiVerson string    `yaml:"apiVersion"`
	Kind      string    `yaml:"kind"`
	Metadata  PMetadata `yaml:"metadata"`
	Spec      PSpec     `yaml:"spec"`
}

// -----------------------
type PMetadata struct {
	Name      string `yaml:"name,omitempty"`
	Namespace string `yaml:"namespace,omitempty"`
}
type PSpec struct {
	Egress      []*Policy   `yaml:"egress,omitempty"`
	Ingress     []*Policy   `yaml:"ingress,omitempty"`
	PodSelector PodSelector `yaml:"podSelector,omitempty"`
	PolicyTypes []string    `yaml:"policyTypes,omitempty"`
}

//-----------------------

type PodSelector struct {
	MatchLabels map[string]string `yaml:"matchLabels,omitempty"`
}

//-------------------------

type MatchLabels struct {
	App string `yaml:"app,omitempty"`
}

// 描述具体规则
type Policy struct {
	Ports []TargetPort `yaml:"ports,omitempty"`
	To    []TargetTo   `yaml:"to,omitempty,omitempty"`
}
type TargetPort struct {
	Port     int    `yaml:"port,omitempty"`
	Protocol string `yaml:"protocol,omitempty"`
}
type TargetTo struct {
	Ipblock     IpBlock     `yaml:"ipBlock,omitempty"`
	NamespaceSelector NamespaceSelector `yaml:"namespaceSelector,omitempty"`
	PodSelector PodSelector `yaml:"podSelector,omitempty"`
}
type IpBlock struct {
	Cidr string `yaml:"cidr,omitempty"`
}

type NamespaceSelector struct{
	
}

// 根据外部ip生成egress
func NewEgress(port int, url string) *Policy {
	policy := new(Policy)
	policy.Ports = []TargetPort{{Port: port, Protocol: "TCP"}}
	policy.To = []TargetTo{{Ipblock: IpBlock{Cidr: url + `/32`}}}
	return policy
}

// 根据port参数生成ingress
func NewIngress(port int) *Policy {
	policy := new(Policy)
	policy.Ports = []TargetPort{{Port: port}}
	return policy
}

// 根据pod名生成策略
func NewPodPolicy(labels map[string]string) *Policy {
	policy := new(Policy)
	policy.To = []TargetTo{{PodSelector: PodSelector{MatchLabels: labels}}}
	return policy
}

// func (podSelector PodSelector) MatchSelector(podSelector2 PodSelector)bool{
// 	if podSelector.MatchLabels.App == podSelector2.MatchLabels.App{
// 		return true
// 	}
// 	return false
// }
