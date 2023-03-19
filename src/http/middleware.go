package http

import (
	"fmt"
	"net/http"
	"sort"
)

type Middleware interface {
	Mount(http.Handler) http.Handler
	Priority() int
}

type Stack struct {
	stack []Middleware
}

func NewMiddlewareStack() *Stack {
	return &Stack{}
}

func (s *Stack) Add(middleware Middleware) *Stack {
	s.stack = append(s.stack, middleware)
	sort.SliceStable(s.stack, func(i, j int) bool {
		return s.stack[i].Priority() < s.stack[j].Priority()
	})
	return s
}

func (s *Stack) Mount(next http.Handler) http.Handler {
	output := next
	for _, middleware := range s.stack {
		fmt.Println(middleware.Priority())
		output = middleware.Mount(output)
	}
	return output
}

func (s *Stack) Priority() int {
	max := 0
	for _, middleware := range s.stack {
		if middleware.Priority() > max {
			max = middleware.Priority()
		}
	}
	return max
}
