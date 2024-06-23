package main

import "fmt"

type Result struct {
	URL      string
	fileName string
	Status   int
	err      error
}

func NewResult(args ...string) *Result {
	res := &Result{}
	switch {
	case len(args) > 0:
		res = res.AddUrl(args[0])
	case len(args) > 1:
		res = res.AddFileName(args[1])
	}
	return res
}

// Добавляет статус к результату и возвращат указатель на него.
func (r *Result) AddStatus(status int) *Result {
	r.Status = status
	return r
}

// Добавляет URL к результату и возвращат указатель на него.
func (r *Result) AddUrl(url string) *Result {
	r.URL = url
	return r
}

// Добавляет имя файла к результату и возвращат указатель на него.
func (r *Result) AddFileName(name string) *Result {
	r.fileName = name
	return r
}

// Добавляет ошибку к результату и возвращат указатель на него.
func (r *Result) AddError(err error) *Result {
	r.err = err
	return r
}

// Распечатка структуры результ
func (r Result) String() string {
	return fmt.Sprintf("StatusCode:%d Download: %s from URL: %s", r.Status, r.fileName, r.URL)
}
