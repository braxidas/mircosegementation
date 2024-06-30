package policygenerate

import (
	"microsegement/mstype"
	
)

func GenerateIpPolicy(k8sServiceList []*mstype.K8sService) []*mstype.NetworkPolicy{
	var result []*mstype.NetworkPolicy
	for _, v := range k8sServiceList {
		networkPolicy := new(mstype.NetworkPolicy)
		networkPolicy.ApiVerson = "networking.k8s.io/v1"
		networkPolicy.Kind = "NetworkPolicy"
		networkPolicy.Metadata.Name = v.PodName + "-ipBlock"
		networkPolicy.Metadata.Namespace = "default"
		networkPolicy.Spec.Egress = v.EgressOut
		networkPolicy.Spec.Ingress = v.IngressOut
		networkPolicy.Spec.PodSelector.MatchLabels.App = v.PodName
		networkPolicy.Spec.PolicyTypes = []string{"Egress", "Igress"}
		result = append(result, networkPolicy)
	}
	return result
}


