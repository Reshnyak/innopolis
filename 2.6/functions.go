package main

func Filter[T any](slice []T, f func(T) bool) []T {
	var res []T
	for _, v := range slice {
		if f(v) {
			res = append(res, v)
		}
	}
	return res
}

func Map[T1, T2 any](slice []T1, f func(T1) T2) []T2 {
	var res []T2
	for _, v := range slice {
		res = append(res, f(v))
	}
	return res
}

func Reduce[T1, T2 any](slice []T1, init T2, f func(T1, T2) T2) T2 {
	res := init
	for _, value := range slice {
		res = f(value, res)
	}
	return res
}
