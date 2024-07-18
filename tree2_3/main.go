package main

import (
	"fmt"

	"github.com/Reshnyak/innopolis/tree2_3/tree23"
)

func main() {
	tree := tree23.NewNode(10)
	tree.Insert(20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 5, 15, 25, 8)
	tree.PrintTree()
	fmt.Println("\n__________________________________________")
	tree.Remove(5, 8, 10, 30, 15)
	tree.PrintTree()
}
