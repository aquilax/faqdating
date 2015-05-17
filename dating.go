package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Dating struct {
	config *Config
	model  Model
	tp     *TransPool
}

type ValidationErrors []string

func NewDating(c *Config) *Dating {
	return &Dating{
		config: c,
	}
}

func (d *Dating) Run() {
	d.model = NewSQLiteModel()

	if err := d.model.Init(d.config); err != nil {
		panic(err)
	}

	d.tp = NewTransPool(d.config.Translations)

	r := mux.NewRouter()
	r.HandleFunc("/", appHandler(d.indexHandler).ServeHTTP).Methods("GET")

	// Static assets
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public_html")))

	http.Handle("/", r)

	log.Printf("Starting server at %s", d.config.Server)
	if err := http.ListenAndServe(d.config.Server, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func (d *Dating) getToken(r *http.Request) string {
	return r.Header.Get(d.config.Token)
}
