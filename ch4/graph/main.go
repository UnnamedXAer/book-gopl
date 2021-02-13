package main

import "fmt"

var graph = make(map[string]map[string]bool)

func main() {

	addEdge("a", "b")
	addEdge("a", "c")
	addEdge("a", "d")
	addEdge("r", "1")
	addEdge("w", "1")
	fmt.Printf("%v", graph)
	fmt.Printf("%v", hasEdge("r", "1"))
	fmt.Printf("%v", hasEdge("r", "2"))
	fmt.Printf("%v", hasEdge("rX", "21"))
}

func addEdge(from, to string) {
	edges := graph[from]

	if edges == nil {
		edges = make(map[string]bool)
		graph[from] = edges
	}

	edges[to] = true
}

func hasEdge(from, to string) bool {
	return graph[from][to]
}
