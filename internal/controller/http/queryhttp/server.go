package queryhttp

import (
	"github.com/cr00z/gocqrs/internal/repository"
	"github.com/cr00z/gocqrs/pkg/util"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type httpServer struct {
	router *mux.Router
	pgRepo repository.RepositoryLister
	elRepo repository.RepositorySearcher
}

func NewHttpServer(pgRepo repository.RepositoryLister, elRepo repository.RepositorySearcher) *httpServer {
	server := &httpServer{
		router: mux.NewRouter(),
		pgRepo: pgRepo,
		elRepo: elRepo,
	}

	server.router.HandleFunc("/messages", server.listMessagesHandler).
		Methods("GET")
	server.router.HandleFunc("/search", server.searchMessagesHandler).
		Methods("GET")

	return server
}

func (s *httpServer) Start() error {
	return http.ListenAndServe(":8080", s.router)
}

func (s *httpServer) listMessagesHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	skip, take := uint64(0), uint64(100)

	skipStr := r.FormValue("skip")
	if skipStr != "" {
		skip, err = strconv.ParseUint(skipStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid skip: "+err.Error())
			return
		}
	}

	takeStr := r.FormValue("take")
	if takeStr != "" {
		take, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid take: "+err.Error())
			return
		}
	}

	ctx := r.Context()
	messages, err := s.pgRepo.List(ctx, skip, take)
	if err != nil {
		util.ResponseError(w, http.StatusInternalServerError, "Couldn't list messages")
		return
	}

	util.ResponseOk(w, messages)
}

func (s *httpServer) searchMessagesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")
	if query == "" {
		util.ResponseError(w, http.StatusBadRequest, "Invalid query")
		return
	}

	var err error
	skip, take := uint64(0), uint64(100)

	skipStr := r.FormValue("skip")
	if skipStr != "" {
		skip, err = strconv.ParseUint(skipStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid skip: "+err.Error())
			return
		}
	}

	takeStr := r.FormValue("take")
	if takeStr != "" {
		take, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid take: "+err.Error())
			return
		}
	}

	ctx := r.Context()
	messages, err := s.elRepo.Search(ctx, query, skip, take)
	if err != nil {
		util.ResponseError(w, http.StatusInternalServerError, "Search: "+err.Error())
		return
	}

	util.ResponseOk(w, messages)
}
