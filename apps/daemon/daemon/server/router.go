package server

import (
	"net/http"

	"github.com/lobofoltran/homelab/apps/daemon/daemon/server/router/grpc"
	"github.com/lobofoltran/homelab/apps/daemon/daemon/server/router/rest"
	"github.com/lobofoltran/homelab/apps/daemon/daemon/server/router/types"
)

func NewPostRoute(path string, fn http.HandlerFunc) types.Route {
	return types.Route{Method: http.MethodPost, Path: path, Handler: fn}
}

func NewGetRoute(path string, fn http.HandlerFunc) types.Route {
	return types.Route{Method: http.MethodGet, Path: path, Handler: fn}
}

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// REST handlers
	for _, r := range rest.NewRouter().Routes() {
		mux.HandleFunc(r.Path, methodHandler(r.Method, r.Handler))
	}

	// gRPC via HTTP/2 (h2c)
	for _, r := range grpc.NewRouter().Routes() {
		mux.HandleFunc(r.Path, methodHandler(r.Method, r.Handler))
	}

	return mux
}

func methodHandler(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}
