package middleware

import "net/http"

type Middleware interface {
	Mount(http.Handler) http.Handler
}

type Stack struct {
	stack []Middleware
}

func NewStack() *Stack {
	return &Stack{}
}

func (s *Stack) Add(middleware Middleware) *Stack {
	s.stack = append(s.stack, middleware)
	return s
}

func (s *Stack) Mount(next http.Handler) http.Handler {
	output := next
	for _, middleware := range s.stack {
		output = middleware.Mount(output)
	}
	return output
}
