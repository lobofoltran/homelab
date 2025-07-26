package types

import "net/http"

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

type Router interface {
	Routes() []Route
}

func NewGetRoute(path string, handler http.HandlerFunc) Route {
	return Route{
		Method:  http.MethodGet,
		Path:    path,
		Handler: handler,
	}
}

func NewPostRoute(path string, handler http.HandlerFunc) Route {
	return Route{
		Method:  http.MethodPost,
		Path:    path,
		Handler: handler,
	}
}
