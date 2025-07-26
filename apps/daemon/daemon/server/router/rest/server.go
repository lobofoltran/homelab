package rest

import (
	"net/http"

	"github.com/lobofoltran/homelab/apps/daemon/daemon/server/router/types"
	"github.com/lobofoltran/homelab/apps/daemon/internal/backend"
	"github.com/lobofoltran/homelab/libs/logger"
)

type restRouter struct {
	routes []types.Route
}

func NewRouter() types.Router {
	r := &restRouter{}
	r.initRoutes()
	return r
}

func (r *restRouter) Routes() []types.Route {
	return r.routes
}

func (r *restRouter) initRoutes() {
	r.routes = []types.Route{
		types.NewGetRoute("/ping", pingHandler),
	}
}

func pingHandler(w http.ResponseWriter, req *http.Request) {
	logger.Debug("Requisição em /ping")

	w.Write([]byte(backend.Ping()))
}
