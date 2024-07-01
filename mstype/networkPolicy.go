package mstype


/*
此类用于生成目标文件中的微隔离策略
*/



//策略
type NetworkPolicy struct{
	ApiVerson string `yaml:"apiVersion"`
	Kind string `yaml:"kind"`
	Metadata PMetadata `yaml:"metadata"`
	Spec PSpec `yaml:"spec"`
}
//-----------------------
type PMetadata struct{
	Name string `yaml:"name,omitempty"`
	Namespace string `yaml:"namespace,omitempty"`
}
type PSpec struct{
	Egress []*Policy	`yaml:"egress,omitempty"`
	Ingress []*Policy	`yaml:"ingress,omitempty"`
	PodSelector PodSelector `yaml:"podSelector,omitempty"`
	PolicyTypes []string `yaml:"policyTypes,omitempty"`
}
//-----------------------

type PodSelector struct{
	MatchLabels MatchLabels `yaml:"matchLabels,omitempty"`
}
//-------------------------

type MatchLabels struct{
	App string `yaml:"app,omitempty"`
}
//描述具体规则
type Policy struct{
	Ports []TargetPort `yaml:"ports,omitempty"`
	To []TargetTo `yaml:"to,omitempty,omitempty"`
}
type TargetPort struct{
	Port int `yaml:"port,omitempty"`
	Protocol string `yaml:"protocol,omitempty"`
}
type TargetTo struct{
	Ipblock IpBlock `yaml:"ipBlock,omitempty"`
	PodSelector PodSelector `yaml:"podSelector,omitempty"`
}
type IpBlock struct{
	Cidr string `yaml:"cidr,omitempty"`
}

func NewEgress(port int, url string )*Policy{
	policy := new(Policy)
	policy.Ports = []TargetPort{TargetPort{Port:port, Protocol:"TCP"}}
	policy.To = []TargetTo{TargetTo{Ipblock:IpBlock{Cidr:url+`/32`}}}
	return policy
}
func NewIngress(port int)*Policy{
	policy := new(Policy)
	policy.Ports = []TargetPort{TargetPort{Port:port}}
	return policy
}