package createhttp

import (
	"github.com/cr00z/gocqrs/internal/controller"
	"github.com/cr00z/gocqrs/internal/domain"
	"github.com/cr00z/gocqrs/internal/repository"
	"github.com/cr00z/gocqrs/pkg/util"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
	"html/template"
	"log"
	"net/http"
	"time"
)

type httpServer struct {
	router *mux.Router
	repo   repository.RepositoryLister
	ctrl   controller.EventStore
}

func NewHttpServer(repo repository.RepositoryLister, ctrl controller.EventStore) *httpServer {
	server := &httpServer{
		router: mux.NewRouter(),
		repo:   repo,
		ctrl:   ctrl,
	}

	server.router.HandleFunc("/messages", server.createMessageHandler).
		Methods("POST").
		Queries("body", "{body}")

	return server
}

func (s *httpServer) Start() error {
	return http.ListenAndServe(":8080", s.router)
}

func (s *httpServer) createMessageHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID string `json:"id"`
	}

	body := template.HTMLEscapeString(r.FormValue("body"))
	if len(body) < 1 || len(body) > 140 {
		util.ResponseError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	createdAt := time.Now().UTC()
	id, err := ksuid.NewRandomWithTime(createdAt)
	if err != nil {
		util.ResponseError(w, http.StatusInternalServerError, "Failed")
		return
	}

	ctx := r.Context()

	msg := domain.Message{
		ID:        id.String(),
		Body:      body,
		CreatedAt: createdAt,
	}

	err = s.repo.Insert(ctx, msg)
	if err != nil {
		log.Print(err)
		util.ResponseError(w, http.StatusInternalServerError, "Failed")
		return
	}

	err = s.ctrl.Publish(msg)
	if err != nil {
		log.Print(err)
	}

	util.ResponseOk(w, response{ID: msg.ID})
}
