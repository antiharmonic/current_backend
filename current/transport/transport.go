package transport
// blatantly stolen from https://github.com/morganhein/backend-takehome-telegraph

import (
	"encoding/json"
	"fmt"
	"github.com/antiharmonic/current_backend/current"
	"net/http"
	_ "strconv"
	_ "strings"
)

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
	eq, err := t.srv.ListMedia()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//marshal
	b, err := json.Marshal(eq)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		//broken pipe?
		fmt.Println(err)
	}
}