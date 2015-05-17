package main

import (
	"database/sql"
	"net/http"
)

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		httpError, ok := err.(HTTPError)
		if ok {
			http.Error(w, httpError.Message, httpError.Code)
			return
		}
		// Default to 500 Internal Server Error
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (d *Dating) indexHandler(w http.ResponseWriter, r *http.Request) error {
	sc := d.config.getSiteConfig(d.getToken(r))
	s := NewSession(sc, d.tp.Get(sc.Language))
	s.AddPath("", s.Lang("Home"))
	return s.render(w, r, sc.templatePath("layout.html"), sc.templatePath("index.html"))
}
