package grpc

import (
	"net/http"

	"github.com/lobofoltran/homelab/apps/daemon/daemon/server/router/types"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

type grpcRouter struct {
	routes     []types.Route
	grpcServer *grpc.Server
}

func NewRouter() types.Router {
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(10<<20),
		grpc.MaxSendMsgSize(10<<20),
	)

	r := &grpcRouter{grpcServer: grpcServer}
	r.initRoutes()

	return r
}

func (r *grpcRouter) Routes() []types.Route {
	return r.routes
}

func (r *grpcRouter) initRoutes() {
	r.routes = []types.Route{
		types.NewPostRoute("/grpc", r.serveGRPC),
	}
}

func (r *grpcRouter) serveGRPC(w http.ResponseWriter, req *http.Request) {
	h2Handler := h2c.NewHandler(r.grpcServer, &http2.Server{})
	h2Handler.ServeHTTP(w, req)
}
