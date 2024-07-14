package main

import "fmt"

func main() {
	mtx := GraphMtx{
		adjMatrix: [][]int{
			{0, 1, 1, 0, 0},
			{1, 0, 0, 1, 1},
			{1, 0, 0, 0, 0},
			{0, 1, 0, 0, 0},
			{0, 1, 0, 0, 0},
		},
	}

	fmt.Printf("%v\n", mtx.BFS(0))
}
