package feeds

import (
	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/gorilla/mux"
)

// HandleSubrouter uses the subrouter for a specific calls and creates a tree of sorts
// handling each route with a different subrouter
func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {

	s = respond.PrepAppRoutes(s, confhandler, appRoutesV2)

}

var appRoutesV2 = []respond.AppRoutes{
	{Name: "feeds.topo.update", Verb: "PUT", Path: "/topology", SubrouterHandler: UpdateTopo},
	{Name: "feeds.topo.get", Verb: "GET", Path: "/topology", SubrouterHandler: ListTopo},
	{Name: "feeds.topo.options", Verb: "OPTIONS", Path: "/topology", SubrouterHandler: Options},
}
