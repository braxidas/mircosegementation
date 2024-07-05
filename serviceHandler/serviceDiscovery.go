package serviceHandler

import (
	"fmt"
	"microsegement/fileHandler"
	"microsegement/mstype"
	"microsegement/soot"
	"sort"
)

var (
	// class2service map[string][]*mstype.K8sService //通过类取得其调用的服务
	// class2Service	map[string]*mstype.K8sService //通过类名取得使用该类的微服务
	aspect2Class map[string]*mstype.JavaClass    //通过切面获得定义该切面的类
	name2JavaClass map[string]*mstype.JavaClass	//通过类名获得其类信息//如果类名存在于key中则表示是自定义类
	graph         mstype.Graph                    //每个类之间的调用关系图
)

func DiscoverService(k8sServiceList []*mstype.K8sService) ([]*mstype.K8sService,error) {

	// class2service = make(map[string][]*mstype.K8sService)
	// class2Service =	make(map[string]*mstype.K8sService)

	for i, _ := range k8sServiceList {
		javaClassList, err := soot.ScanDiscoverService(k8sServiceList[i].FilePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(javaClassList) == 1 {
			continue
		}
		k8sServiceList[i].JavaClassList = javaClassList
		for _, v := range k8sServiceList[i].JavaClassList{//建议类名和使用该类的微服务的索引
			// class2Service[v.ClassName] = k8sServiceList[i]
			v.K8sService = k8sServiceList[i]
		}
	}

	analysisCallGraph(k8sServiceList)
	analysisDirectCall(k8sServiceList)
	analysisDubboCall(k8sServiceList)

	for _, v := range k8sServiceList {
		//生成json文件
		err := fileHandler.WriteToJson(v)
		if err != nil {
			fmt.Println(err)
		}
		//根据pod名称 向最终的策略文件中添加策略
		for k, _ := range v.Egress{
			v.EgressOut = append(v.EgressOut, mstype.NewPodPolicy(k.PodName))
		}
		// 
		
	}

	return k8sServiceList,nil
}

// 扫描直接调用微服务
func analysisDirectCall(k8sServiceList []*mstype.K8sService) {
	for i, _ := range k8sServiceList { //每个模块
		for _, v := range k8sServiceList[i].JavaClassAllList { //每个类名
			_, ok1 := name2JavaClass[v] //如果该类为自定义的类
			if ok1{
				for _, vc := range name2JavaClass[v].Consume { //每个类声明的调用
					ks, ok := name2K8sService[vc]
					if ok {
						k8sServiceList[i].AppendEgress(ks)
						ks.AppendIngress(k8sServiceList[i])
					}
					// updateClass2Service(name2JavaClass[v].ClassName, ks)
				}
			}
		}
	}
}

// 扫描dubbo调用微服务
func analysisDubboCall(k8sServiceList []*mstype.K8sService) {
	for i, _ := range k8sServiceList { //每个模块
		for _, v := range k8sServiceList[i].JavaClassAllList { //每个类
			_, ok1 := name2JavaClass[v]
			if ok1{
				for _, vd := range name2JavaClass[v].DubboReference { //每个类声明的dubbo远程调用
					for j, _ := range k8sServiceList {
						if k8sServiceList[j].ProvideService(vd) {
							k8sServiceList[i].AppendEgress(k8sServiceList[j])
							k8sServiceList[j].AppendIngress(k8sServiceList[i])
							// updateClass2Service(name2JavaClass[v].ClassName, k8sServiceList[j])
							break
						}
					}
				}
			}
		}
	}
}

// 扫描间接调用未服务，需要class间的成员变量调用和切面调用的信息
func analysisCallGraph(k8sServiceList []*mstype.K8sService) {
	buildGraph(k8sServiceList)
	for _, vn := range graph.Nodes{
		vc, ok := name2JavaClass[vn.Value]//若该节点的类被模块使用，则有必要检查其可达节点
		if ok{
			graph.Reachable(vn, func(node *mstype.Node) {
				vc.K8sService.JavaClassAllList = append(vc.K8sService.JavaClassAllList, node.Value)
			})
		}
	}
}

/* 
构造图 标记每个类之间的调用关系
若一个类A的成员变量中使用了自动注入的另一个类B的变量
若一个类A的成员方法中使用了另一个类B定义的Aspect注解
则视作A调用了B，B的调用微服务A也会调用,
则建立A到B的边，则通过遍历A的可达节点可知其实际可能使用的全部类
*/
func buildGraph(k8sServiceList []*mstype.K8sService) {

	//先建立切面和定义切面的类的联系
	aspect2Class = make(map[string]*mstype.JavaClass)
	for _, v := range k8sServiceList { //每个模块
		for _, vc := range v.JavaClassList { //每个类
			for _, vca := range vc.DefineAspect{//该类定义的切面
				aspect2Class[vca] = vc
			}
		}
	}

	//初始化图
	graph = mstype.Graph{}
	//增加点
	name2JavaClass = make(map[string]*mstype.JavaClass)
	for _, v := range k8sServiceList { //每个模块
		for _, vc := range v.JavaClassList { //每个类
			name2JavaClass[vc.ClassName] = vc //建立类到类名的索引
			graph.AddNode(&mstype.Node{Value:vc.ClassName})
		}
	}
	//增加边
	for _, v := range k8sServiceList { //每个模块
		for _, vc := range v.JavaClassList { //每个类
			for _, vcf := range vc.Field{ //若A的成员变量里有B
				graph.AddEdge(&mstype.Node{Value:vc.ClassName}, &mstype.Node{Value:vcf})
			}
			for _, vca := range vc.UseAspect{//若A使用的注解为B定义的切面
				vcac, ok := aspect2Class[vca]
				if ok{
					graph.AddEdge(&mstype.Node{Value:vc.ClassName}, &mstype.Node{Value:vcac.ClassName})
				}
			}
		}
	}
	
}

// 删除数组中重复内容
func removeDuplicates(elements []string) []string {
	if len(elements) < 2 {
		return elements
	}
	sort.Strings(elements) // 先对字符串数组进行排序
	j := 0
	for i := 1; i < len(elements); i++ {
		if elements[i] != elements[j] {
			j++
			elements[j] = elements[i]
		}
	}
	return elements[:j+1]
}

// // 根据每个类新增的调用更新其调用列表
// func updateClass2Service(className string, k8sService *mstype.K8sService) {
// 	_, ok := class2service[className]
// 	if !ok {
// 		class2service[className] = make([]*mstype.K8sService, 0, 2)
// 	}
// 	class2service[className] = append(class2service[className], k8sService)
// }
