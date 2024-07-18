package main

type GraphMtx struct {
	adjMatrix [][]int
}

func (g *GraphMtx) BFS(startVertex int) []int {
	visited := make([]bool, len(g.adjMatrix))
	path := make([]int, 0, len(g.adjMatrix))
	var queue []int
	visited[startVertex] = true
	queue = append(queue, startVertex)
	for len(queue) != 0 {
		currentVertex := queue[0]
		path = append(path, currentVertex)
		for i := 0; i < len(g.adjMatrix); i++ {
			if !visited[i] && g.adjMatrix[currentVertex][i] == 1 {
				visited[i] = true
				queue = append(queue, i)
			}
		}
		queue = queue[1:]
	}
	return path
}
