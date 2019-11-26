package main

import (
	"flag"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/rustyeddy/store"
	log "github.com/sirupsen/logrus"
)

var (
	addr *string
	name *string
	base *string

	storage *store.FileStore
)

func init() {
	name = flag.String("name", "default", "The name of the storage")
	base = flag.String("base", "/srv/storge", "The base directory for the stores")
	addr = flag.String("addr", "0.0.0.0:3333", "Address and port to run on")
}

func main() {
	var err error
	flag.Parse()

	// TODO change to os.join
	if storage, err = store.UseFileStore(filepath.Join(*base, *name)); err != nil {
		log.Fatalln(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/health", handleHealth)
	router.HandleFunc("/api/store/{name}", handleStore)

	spa := spaHandler{staticPath: "pub", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	addr := "0.0.0.0:8444"
	srv := &http.Server{
		Handler: router,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infof("Store spa listing for connections: %s ", addr)
	log.Fatal(srv.ListenAndServe())
}
