package main

import "fmt"

func IsEqualArrays[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	mapA := make(map[T]int, len(a))
	mapB := make(map[T]int, len(b))
	for i, v := range a {
		mapA[v]++
		mapB[b[i]]++
	}
	for k, va := range mapA {
		if vb, ok := mapB[k]; !ok || va != vb {
			return false
		}
	}
	return true
}
func main() {
	arrA := [...]int{3, 4, 2, 9, 1, 5, 8, 7, 6}
	arrB := [...]int{9, 8, 7, 6, 5, 3, 4, 2, 1}
	fmt.Println(IsEqualArrays(arrA[:], arrB[:]))
}
