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
	s := NewSession(r, sc, d.tp.Get(sc.Language))
	s.AddPath("", s.Lang("Home"))
	return s.render(w, sc.templatePath("layout.html"), sc.templatePath("index.html"))
}

func (d *Dating) authRegisterHandler(w http.ResponseWriter, r *http.Request) error {
	sc := d.config.getSiteConfig(d.getToken(r))
	l := d.tp.Get(sc.Language)
	s := NewSession(r, sc, l)

	if s.Logged() {
		http.Redirect(w, r, "/user/profile.html", http.StatusTemporaryRedirect)
	}
	if r.Method == "POST" {
		v := NewValidation(l)
		email, password, ve := v.validateUserRegister(r)
		if len(ve) == 0 {
			userID, err := d.model.RegisterUser(email, password)
			if err != nil {
				return err
			}
			s.logInUser(userID, w)
			// Redirect to the profile page
			http.Redirect(w, r, "/user/profile.html", http.StatusFound)
			return nil
		}
	}
	return s.render(w, sc.templatePath("layout.html"), sc.templatePath("auth/register.html"))
}

func (d *Dating) authLoginHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (d *Dating) authLogoutHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (d *Dating) userProfileHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}
