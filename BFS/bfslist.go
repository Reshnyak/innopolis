package main

type GraphList struct {
	adjList map[int][]int
}
type emptyVal struct{}

func (g *GraphList) BFS(startVertex int) []int {
	path := make([]int, 0, len(g.adjList))
	visited := make(map[int]emptyVal)
	var queue []int
	visited[startVertex] = emptyVal{}
	queue = append(queue, startVertex)
	for len(queue) != 0 {
		currentVertex := queue[0]
		path = append(path, currentVertex)
		for _, adjVertex := range g.adjList[currentVertex] {
			if _, ok := visited[adjVertex]; !ok {
				visited[adjVertex] = emptyVal{}
				queue = append(queue, adjVertex)
			}
		}
		queue = queue[1:]
	}
	return path
}
