package main

import (
	"fmt"
)

func debug(format string, a ...interface{}) {
	// fmt.Fprintf(os.Stderr, format, a...)
}

type Edge struct {
	x int
	y int
}

type Graph map[int][]int

func (g Graph) root() int {
	var queue []int
	for n, _ := range g {
		queue = append(queue, n)
	}
	root := queue[0]
	for len(queue) > 0 {
		var n int
		n, queue = queue[0], queue[1:]
		for _, edge := range g[n] {
			if edge == root {
				root = n
			}
		}
	}
	return root
}
func (g Graph) Length() int {
	length := -1
	for n, _ := range g {
		l := Length(n, g)
		if l > length {
			length = l
		}
	}
	return length
}

func Length(root int, g Graph) int {
	debug("Longest(%v, %v)\n", root, g[root])
	if len(g[root]) == 0 {
		debug("return 1\n")
		return 1
	}
	longest := -1
	for _, edge := range g[root] {
		length := 1 + Length(edge, g)
		debug("longer? %v\n", longest < length)
		if longest <= length {
			longest = length
		}
	}
	debug("return (%v)\n", longest)
	return longest
}

func main() {
	// n: the number of relationships of influence
	var n int
	fmt.Scan(&n)
	graph := make(Graph)

	for i := 0; i < n; i++ {
		// x: a relationship of influence between two people (x influences y)
		var x, y int
		fmt.Scan(&x, &y)
		if _, ok := graph[x]; !ok {
			graph[x] = make([]int, 0)
		}
		graph[x] = append(graph[x], y)
		if _, ok := graph[y]; !ok {
			graph[y] = make([]int, 0)
		}
	}

	// fmt.Fprintln(os.Stderr, "Debug messages...")
	debug("graph: %#v\n", graph)
	debug("root: %#v\n", graph.root())
	// The number of people involved in the longest succession of influences
	fmt.Println(graph.Length())
}
