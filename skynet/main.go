package main

import (
	"container/heap"
	"fmt"
	"os"
	"sort"
)

func debug(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", format), a...)
}

type PItem struct {
	value    interface{}
	priority int
	index    int
}
type PQueue []*PItem

func (q PQueue) String() string {
	o := "("
	for i, item := range q {
		if i != 0 {
			o += ", "
		}
		o += fmt.Sprintf("(%d, %v)", item.priority, item.value)

	}
	return o
}
func (q PQueue) Len() int           { return len(q) }
func (q PQueue) Less(i, j int) bool { return q[i].priority > q[j].priority }
func (q PQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}
func (q *PQueue) Push(x interface{}) {
	n := len(*q)
	item := x.(*PItem)
	item.index = n
	*q = append(*q, item)
}
func (q *PQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	item.index = -1
	*q = old[0 : n-1]
	return item
}
func (q *PQueue) update(value interface{}, priority int) {
	for _, item := range *q {
		if item.value == value {
			item.priority = priority
			break
		}
	}
	//	debug("pq update: v: %v, p: %v", value, priority)
	sort.Sort(q)
}

type Edge [2]int

func NewEdge(a, b int) *Edge {
	return &Edge{a, b}
}
func (e Edge) Between(a, b int) bool {
	if (e[0] == a && e[1] == b) ||
		(e[0] == b && e[1] == a) {
		return true
	}
	return false
}

type Graph struct {
	size  int
	edges []Edge
	exits []int
}

func NewGraph(n int) *Graph {
	return &Graph{n, make([]Edge, 0), make([]int, 0)}
}

func (g Graph) neighbours(n int) []int {
	neighbours := make([]int, 0)
	for _, edge := range g.edges {
		if n == edge[0] {
			neighbours = append(neighbours, edge[1])
		} else if n == edge[1] {
			neighbours = append(neighbours, edge[0])
		}
	}
	return neighbours
}

func (g Graph) djikstra(start int) (dist, prev []int) {
	dist = make([]int, g.size)
	prev = make([]int, g.size)
	pq := make(PQueue, g.size)
	dist[start] = 0
	for v := 0; v < g.size; v++ {
		if v != start {
			dist[v] = 90000
			prev[v] = -1
		}
		pq[v] = &PItem{v, dist[v], v}
	}
	heap.Init(&pq)
	sort.Sort(&pq)
	for len(pq) > 0 {
		//		debug("step ->\n\tdist: %v\n\tprev: %v\n\tpq: %s", dist, prev, pq)
		u := pq.Pop().(*PItem)
		for _, v := range g.neighbours(u.value.(int)) {
			old := dist[u.value.(int)] + 1
			//			debug("before -> u: %v, v: %v, old: %v, dist: %v", u.value, v, old, dist[v])
			if old < dist[v] {
				dist[v] = old
				prev[v] = u.value.(int)
				pq.update(v, old)
				//				debug("update ->\n\tdist: %v\n\tprev: %v\n\tpq: %s", dist, prev, pq)
			}
		}
	}
	prev[start] = -1
	return dist, prev
}

func (g *Graph) Cut(a, b int) {
	debug("new edges: %v", g.edges)
	for idx, edge := range g.edges {
		if edge.Between(a, b) {
			g.edges = append(g.edges[:idx], g.edges[idx+1:]...)
		}
	}
	debug("new edges: %v", g.edges)
}

func PathTo(prev []int, target int) []int {
	path := make([]int, 0)
	for prev[target] != -1 {
		path = append([]int{target}, path...)
		target = prev[target]
	}
	debug("p: %v", path)
	return path
}

type PathList [][]int

func (l PathList) Len() int           { return len(l) }
func (l PathList) Less(i, j int) bool { return len(l[i]) < len(l[j]) }
func (l PathList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

func (g *Graph) Turn(SI int) string {
	var cut [2]int

	dist, prev := g.djikstra(SI)
	debug("%#v, %#v", dist, prev)
	list := make(PathList, len(g.exits))
	for idx, exit := range g.exits {
		list[idx] = PathTo(prev, exit)
	}
	sort.Sort(list)
	for _, shortest := range list {
		if len(shortest) < 1 {
			continue
		}
		debug("shortest: %v", shortest)
		cut[0] = shortest[len(shortest)-1]
		cut[1] = SI
		if len(shortest) > 1 {
			cut[1] = shortest[len(shortest)-2]
		}
		break
	}
	// fmt.Fprintln(os.Stderr, "Debug messages...")

	// Example: 3 4 are the indices of the nodes you wish to sever the link between
	debug("cut: %v", cut)
	g.Cut(cut[0], cut[1])
	return fmt.Sprintf("%d %d", cut[0], cut[1])
}

func main() {
	// N: the total number of nodes in the level, including the gateways
	// L: the number of links
	// E: the number of exit gateways
	var N, L, E int
	fmt.Scan(&N, &L, &E)
	graph := NewGraph(N)
	for i := 0; i < L; i++ {
		// N1: N1 and N2 defines a link between these nodes
		var N1, N2 int
		fmt.Scan(&N1, &N2)
		graph.edges = append(graph.edges, *NewEdge(N1, N2))
	}
	for i := 0; i < E; i++ {
		// EI: the index of a gateway node
		var EI int
		fmt.Scan(&EI)
		graph.exits = append(graph.exits, EI)
	}

	for {
		// SI: The index of the node on which the Skynet agent is positioned this turn
		var SI int
		fmt.Scan(&SI)
		fmt.Println(graph.Turn(SI))
	}
}
