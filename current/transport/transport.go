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

func SetJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func (t httpEndpoint) ListMedia(w http.ResponseWriter, r *http.Request) {
	p_media_type := r.URL.Query().Get("type")
	p_limit := r.URL.Query().Get("limit")
	genre := r.URL.Query().Get("genre")

	media_type, err := strconv.Atoi(p_media_type)
	if err != nil {
		media_type = 0
	}

	limit, err := strconv.Atoi(p_limit)
	if err != nil {
		limit = 0
	}
	
	log.Printf("media_type: %d, limit: %d\n", media_type, limit)
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

	SetJSON(w)
	_, err = w.Write(b)
	if err != nil {
		//broken pipe?
		log.Println(err)
	}
}

func (t httpEndpoint) ListRecentMedia(w http.ResponseWriter, r *http.Request) {
	p_media_type := r.URL.Query().Get("type")
	p_limit := r.URL.Query().Get("limit")
	limit := 10
	if p_limit != "" {
		n, err := strconv.Atoi(p_limit)
		if err == nil {
			limit = n
		}
	}
	media_type, err := strconv.Atoi(p_media_type)
	if err != nil {
		media_type = 0
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

	SetJSON(w)
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

	SetJSON(w)
	_, err = w.Write(b)
	if err != nil {
		log.Println(err)
	}
}

func (t httpEndpoint) prioritizeMedia(w http.ResponseWriter, r *http.Request, priority bool) {
	vars := mux.Vars(r)
	if vars["id"] == "" {
		http.Error(w, "This endpoint requires an ID", 500)
	}

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var m *current.Media
	if priority {
		m, err = t.srv.UpgradeMedia(id)
	} else {
		m, err = t.srv.DowngradeMedia(id)
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	b, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	SetJSON(w)
	_, err = w.Write(b)
	if err != nil {
		log.Println(err)
	}
}

func (t httpEndpoint) UpgradeMedia(w http.ResponseWriter, r *http.Request) {
	t.prioritizeMedia(w, r, true)
}

func (t httpEndpoint) DowngradeMedia(w http.ResponseWriter, r *http.Request) {
	t.prioritizeMedia(w, r, false)
}

func (t httpEndpoint) SearchMedia(w http.ResponseWriter, r *http.Request) {
	var vars map[string]string
	vars = make(map[string]string)
	vars["type"] = r.URL.Query().Get("type")
	vars["id"] = r.URL.Query().Get("id")
	vars["title"] = r.URL.Query().Get("title")
	log.Println(vars)
	if vars["id"] == "" && vars["title"] == "" && vars["type"] == "" {
		http.Error(w, "This endpoint requires at least one of: id, title, type.", 500)
		return
	}
	var id int
	var err error
	// if the id is defined, return a single item and ignore other params.
	if vars["id"] != "" {
		id, err = strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		m, err := t.srv.GetMediaByID(id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		b, err := json.Marshal(m)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		SetJSON(w)
		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
		}
		return
	}

	my_type, err := strconv.Atoi(vars["type"])
	if err != nil {
		my_type = 0
	}

	//mq := current.MediaQuery{Media: current.Media{Title: title, MediaType: media_type}}
	// i'm not sure if this is the right way to do this. should it be a struct so the values that aren't 
	// being used can be nil? or does it not matter since if they're not defined or 0 I can just ignore them in the model side?
	// should I make a struct like SearchableMedia?
	m, err := t.srv.SearchMedia(vars["title"], my_type)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	b, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	SetJSON(w)
	_, err = w.Write(b)
	if err != nil {
		log.Println(err)
	}
}

func (t httpEndpoint) TopMedia(w http.ResponseWriter, r *http.Request) {
	p_media_type := r.URL.Query().Get("type")
	media_type, err := strconv.Atoi(p_media_type)
	if err != nil {
		media_type = 0
	}
	m, err := t.srv.TopMedia(media_type)
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

	SetJSON(w)
	_, err = w.Write(b)
	if err != nil {
		log.Println(err)
	}
}