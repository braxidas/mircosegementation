package mstype

import "fmt"


/*
此类用于分析服务间的类调用关系
*/
type Node struct {
	Value string//节点 通过类名表示  原因是并非所有类都是模块中自定义的，比如redisService之类的，但这些类可能有用，但没有指针

}

type Graph struct {
	Nodes []*Node          // 节点集
	Edges map[Node][]*Node // 邻接表表示的无向图
}

// 增加节点
func (g *Graph) AddNode(n *Node) {
	g.Nodes = append(g.Nodes, n)
}

// 增加边
func (g *Graph) AddEdge(u, v *Node) {
	// 首次建立图
	if g.Edges == nil {
		g.Edges = make(map[Node][]*Node)
	}
	g.Edges[*u] = append(g.Edges[*u], v) // 建立 u->v 的边
	// g.edges[*v] = append(g.edges[*v], u) // 由于是无向图，同时存在 v->u 的边
}

// 输出图
func (g *Graph) String() {
	str := ""
	for _, iNode := range g.Nodes {
		str += iNode.String() + " -> "
		nexts := g.Edges[*iNode]
		for _, next := range nexts {
			str += next.String() + " "
		}
		str += "\n"
	}
	fmt.Println(str)
}

// 输出节点
func (n *Node) String() string {
	return fmt.Sprintf("%v", n.Value)
}


type NodeQueue struct {
	nodes []Node
}

// 实现 BFS 遍历
func (g *Graph) BFS(f func(node *Node)) {

	// 初始化队列
	q := NewNodeQueue()
	// 取图的第一个节点入队列
	head := g.Nodes[0]
	q.Enqueue(*head)
	// 标识节点是否已经被访问过
	visited := make(map[*Node]bool)
	visited[head] = true
	// 遍历所有节点直到队列为空
	for {
		if q.IsEmpty() {
			break
		}
		node := q.Dequeue()
		visited[node] = true
		nexts := g.Edges[*node]
		// 将所有未访问过的邻接节点入队列
		for _, next := range nexts {
			// 如果节点已被访问过
			if visited[next] {
				continue
			}
			q.Enqueue(*next)
			visited[next] = true
		}
		// 对每个正在遍历的节点执行回调
		if f != nil {
			f(node)
		}
	}
}


// 实现一个节点的所有可达节点
func (g *Graph) Reachable(head *Node, f func(node *Node)) {

	// 初始化队列
	q := NewNodeQueue()
	// 取图的第一个节点入队列
	q.Enqueue(*head)
	// 标识节点是否已经被访问过
	visited := make(map[*Node]bool)
	visited[head] = true
	// 遍历所有节点直到队列为空
	for {
		if q.IsEmpty() {
			break
		}
		node := q.Dequeue()
		visited[node] = true
		nexts := g.Edges[*node]
		// 将所有未访问过的邻接节点入队列
		for _, next := range nexts {
			// 如果节点已被访问过
			if visited[next] {
				continue
			}
			q.Enqueue(*next)
			visited[next] = true
		}
		// 对每个正在遍历的节点执行回调
		if f != nil {
			f(node)
		}
	}
}

// 生成节点队列
func NewNodeQueue() *NodeQueue {
	q := NodeQueue{}
	q.nodes = []Node{}
	return &q
}

// 入队列
func (q *NodeQueue) Enqueue(n Node) {
	q.nodes = append(q.nodes, n)
}

// 出队列
func (q *NodeQueue) Dequeue() *Node {
	node := q.nodes[0]
	q.nodes = q.nodes[1:]
	return &node
}

// 判空
func (q *NodeQueue) IsEmpty() bool {
	return len(q.nodes) == 0
}


