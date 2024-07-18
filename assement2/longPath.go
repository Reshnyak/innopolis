package main

type GraphMtx2 struct {
	adjMatrix [][]int
}

func (g *GraphMtx2) findMaxGradeByStartNode(input <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for startVertex := range input {
			visited := make([]bool, len(g.adjMatrix))
			grade := 0
			out <- g.DFSUtil(startVertex, visited, grade)
		}
	}()
	return out
}

func (g *GraphMtx2) DFSUtil(vertex int, visited []bool, grade int) int {
	visited[vertex] = true
	maxGrade := grade
	for i := 0; i < len(g.adjMatrix); i++ {
		if !visited[i] && g.adjMatrix[vertex][i] != 0 {
			g := g.DFSUtil(i, visited, grade+g.adjMatrix[vertex][i])
			if g > maxGrade {
				maxGrade = g
			}
		}
	}
	visited[vertex] = false
	return maxGrade
}
