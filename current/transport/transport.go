package transport
// blatantly stolen from https://github.com/morganhein/backend-takehome-telegraph

import (
	"encoding/json"
	_ "fmt"
	"github.com/antiharmonic/current_backend/current"
	"net/http"
	"strconv"
	_ "strings"
	"log"
	"github.com/gorilla/mux"
)

// what is this and why is it here? I haven't updated it with new methods and everything continues to work as intended...
type Transport interface {
	ListMedia(w http.ResponseWriter, r *http.Request)
}

func NewHTTPTransport(srv current.Service) *httpEndpoint {
	return &httpEndpoint {
		srv: srv,
	}
}

type httpEndpoint struct {
	srv current.Service
}

func (t httpEndpoint) ListMedia(w http.ResponseWriter, r *http.Request) {
	media_type := r.URL.Query().Get("type")
	limit := r.URL.Query().Get("limit")
	genre := r.URL.Query().Get("genre")
	log.Printf("media_type: %s, limit: %s\n", media_type, limit)
	m, err := t.srv.ListMedia(media_type, limit, genre)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//marshal
	b, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		//broken pipe?
		log.Println(err)
	}
}

func (t httpEndpoint) ListRecentMedia(w http.ResponseWriter, r *http.Request) {
	media_type := r.URL.Query().Get("type")
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "10"
	}
	m, err := t.srv.ListRecentMedia(media_type, limit)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//marshal
	b, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		log.Println(err)
	}
}

func (t httpEndpoint) StartMedia(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["id"] == "" {
		http.Error(w, "This endpoint requires an ID", 500)
	}

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	m, err := t.srv.StartMedia(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	b, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Println(err)
	}
}