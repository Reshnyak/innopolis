package main

import (
	"fmt"
	"sort"
)

type Candidate struct {
	Name  string
	Votes int
}

func CalculateVotes(names []string) []Candidate {
	candidatesMap := make(map[string]int, len(names))
	for _, name := range names {
		candidatesMap[name]++
	}
	candidates := make([]Candidate, 0, len(candidatesMap))
	for k, v := range candidatesMap {
		candidates = append(candidates, Candidate{Name: k, Votes: v})
	}
	sort.Slice(candidates, func(i int, j int) bool {
		return candidates[i].Votes > candidates[j].Votes
	})
	return candidates
}
func main() {
	names := []string{"Ann", "Kate", "Peter", "Kate", "Ann", "Ann", "Helen"}
	fmt.Println(CalculateVotes(names))
}
