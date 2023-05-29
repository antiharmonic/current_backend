package main

// architecture blatantly stolen from @morganhein https://github.com/morganhein/backend-takehome-telegraph

import (
  "fmt"
  "log"
  "net/http"
  "net/url"
  "os"

  "github.com/antiharmonic/current_backend/current"
  "github.com/antiharmonic/current_backend/current/store"
  "github.com/antiharmonic/current_backend/current/transport"


  "github.com/gookit/config/v2"
  "github.com/gookit/config/v2/ini"
  "github.com/gorilla/mux"
  "github.com/gorilla/handlers"
)

func main() {
  // load config
  config.WithOptions(config.ParseEnv)
  config.AddDriver(ini.Driver)

  err := config.LoadFiles("config.ini")
  if err != nil {
    log.Fatal(err)
  }
  connstr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
                          config.String("database.user"),
                          url.QueryEscape(config.String("database.pass")),
                          config.String("database.host"),
                          config.String("database.port"),
                          config.String("database.db_name"),
                        )
  db, err := store.CreatePostgresStore(connstr)
  if err != nil {
    log.Fatal(err)
  }

  srv := transport.NewHTTPTransport(current.New(db))
  // /healthcheck
  router := mux.NewRouter()
  router.HandleFunc("/health-check", HealthCheck).Methods("Get")

  // /push

  // /pop

  // /update

  // /upgrade/id (priority)
  // /downgrade/id (priority)

  // /random

  // /delete

  // /start

  // /search

  // /list
  router.HandleFunc("/list", srv.ListMedia).Methods("Get")
  // /list/recent

  // /count

  // /top
  loggedRouter := handlers.LoggingHandler(os.Stdout, router)
  log.Println("Starting server on " + config.String("server.port"))
  if err := http.ListenAndServe("0.0.0.0:" + config.String("server.port"), loggedRouter); err != nil {
    log.Fatal(err)
  }
  //fmt.Printf("config data:\n%#v\n", config.Data())
}

func SetJSON(w http.ResponseWriter) {
  w.Header().Set("Content-Type", "application/json")
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  fmt.Fprintf(w, "API is up and running")
}

func ListMedia(w http.ResponseWriter, r *http.Request) {
  SetJSON(w)
  w.WriteHeader(http.StatusOK)
  
}