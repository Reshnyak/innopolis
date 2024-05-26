package main

import (
	"fmt"
	"sort"
)

type Candidate struct {
	Name  string
	Votes int
}

// Функция подсчета каждого голоса за кандидата. Входной аргумент - массив с именами кандидатов.
// Результативный - массив структуры Candidate, отсортированный по убыванию количества голосов.
func CalculateVotes(names []string) []Candidate {
	//Сложили в мапу
	candidatesMap := make(map[string]int, len(names))
	for _, name := range names {
		candidatesMap[name]++
	}

	candidates := make([]Candidate, 0, len(candidatesMap))
	//Сложили в слайс
	for k, v := range candidatesMap {
		candidates = append(candidates, Candidate{Name: k, Votes: v})
	}
	//Отсорторовали в нужную сторону
	sort.Slice(candidates, func(i int, j int) bool {
		return candidates[i].Votes > candidates[j].Votes
	})
	return candidates
}
func main() {
	names := []string{"Ann", "Kate", "Peter", "Kate", "Ann", "Ann", "Helen"}
	fmt.Println(CalculateVotes(names))
}
