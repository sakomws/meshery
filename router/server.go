package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/layer5io/meshery/models"
)

type Router struct {
	s    *http.ServeMux
	port int
}

// New returns a new ServeMux with app routes.
func NewRouter(ctx context.Context, h models.HandlerInterface, port int) *Router {
	mux := http.NewServeMux()

	mux.Handle("/api/user", h.AuthMiddleware(http.HandlerFunc(h.UserHandler)))

	mux.Handle("/api/k8sconfig", h.AuthMiddleware(http.HandlerFunc(h.K8SConfigHandler)))
	mux.Handle("/api/load-test", h.AuthMiddleware(http.HandlerFunc(h.LoadTestHandler)))
	mux.Handle("/api/results", h.AuthMiddleware(http.HandlerFunc(h.FetchResultsHandler)))

	mux.Handle("/api/mesh/manage", h.AuthMiddleware(http.HandlerFunc(h.MeshAdapterConfigHandler)))
	mux.Handle("/api/mesh/ops", h.AuthMiddleware(http.HandlerFunc(h.MeshOpsHandler)))
	mux.Handle("/api/mesh/adapters", h.AuthMiddleware(http.HandlerFunc(h.GetAllAdaptersHandler)))
	mux.Handle("/api/events", h.AuthMiddleware(http.HandlerFunc(h.EventStreamHandler)))

	mux.Handle("/api/grafana/config", h.AuthMiddleware(http.HandlerFunc(h.GrafanaConfigHandler)))
	mux.Handle("/api/grafana/boards", h.AuthMiddleware(http.HandlerFunc(h.GrafanaBoardsHandler)))
	mux.Handle("/api/grafana/query", h.AuthMiddleware(http.HandlerFunc(h.GrafanaQueryHandler)))

	mux.HandleFunc("/logout", h.LogoutHandler)
	mux.HandleFunc("/login", h.LoginHandler)

	// TODO: have to change this too
	mux.Handle("/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hr
		http.ServeFile(w, r, "../ui/out/static/img/meshery-logo.png")
	}))
	mux.Handle("/", h.AuthMiddleware(http.FileServer(http.Dir("../ui/out/"))))

	return &Router{
		s:    mux,
		port: port,
	}
}

func (r *Router) Run() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", r.port), r.s)
}
