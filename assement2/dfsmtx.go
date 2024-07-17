package main

//import "slices"

type GraphMtx struct {
	adjMatrix [][]int
}

func (g *GraphMtx) CalcMaxGrade() int {
	maxGrade := 0
	for i := 0; i < len(g.adjMatrix); i++ {
		_, grade := g.DFS(i)
		if maxGrade < grade {
			maxGrade = grade
		}

	}
	return maxGrade
}
func (g *GraphMtx) DFS(startVertex int) ([]int, int) {
	visited := make([]bool, len(g.adjMatrix))
	path := make([]int, 0, len(g.adjMatrix))
	grade := 0
	path, grade = g.DFSUtil(startVertex, visited, path, grade)

	return path, grade
}
func (g *GraphMtx) DFSUtil(vertex int, visited []bool, path []int, grade int) ([]int, int) {
	visited[vertex] = true
	path = append(path, vertex)
	pathMax := make([]int, len(path))
	pa := make([]int, len(path))
	copy(pathMax, path)
	copy(pa, path)
	maxGrade := grade

	for i := 0; i < len(g.adjMatrix); i++ {
		if !visited[i] && g.adjMatrix[vertex][i] != 0 {
			vis := make([]bool, len(g.adjMatrix))
			copy(vis, visited)
			p, g := g.DFSUtil(i, vis, pa, grade+g.adjMatrix[vertex][i])
			if g > maxGrade {
				pathMax = p
				maxGrade = g
			}
		}
	}
	return pathMax, maxGrade
}

type GraphLists struct {
	adjList map[int][]int
}

type emptyVal struct{}

func (g *GraphLists) DFS(startVertex int) []int {
	visited := make(map[int]emptyVal, len(g.adjList))
	path := g.dFSUtil(startVertex, []int{}, visited)
	return path
}
func (g *GraphLists) dFSUtil(vertex int, path []int, visited map[int]emptyVal) []int {
	visited[vertex] = emptyVal{}
	path = append(path, vertex)
	for _, v := range g.adjList[vertex] {
		if _, ok := visited[v]; !ok {
			path = g.dFSUtil(v, path, visited)
		}
	}
	return path
}
