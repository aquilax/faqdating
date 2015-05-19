package main

import (
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
)

const sessionName = "faqd"

type PathLink struct {
	URL   string
	Label string
}

type Session struct {
	r      *http.Request
	td     TemplateData
	ln     *Language
	store  sessions.Store
	path   []*PathLink
	userID int
}

type TemplateData map[string]interface{}

func NewSession(r *http.Request, sc *SiteConfig, ln *Language) *Session {
	return &Session{
		r:     r,
		td:    NewTemplateData(sc),
		ln:    ln,
		path:  []*PathLink{},
		store: sessions.NewCookieStore([]byte(sc.SessionSecret)),
	}
}

func NewTemplateData(sc *SiteConfig) TemplateData {
	td := make(TemplateData)
	td.Set("Title", sc.Title)
	td.Set("LanguageCode", sc.LanguageCode)
	td.Set("Description", sc.Description)
	td.Set("ShowVote", false)
	td.Set("Css", sc.CSS)
	td.Set("FormTitle", "")
	td.Set("Analytics", sc.Analytics)
	td.Set("Domain", sc.Domain)
	td.Set("PostHeader", template.HTML(sc.PostHeader))
	td.Set("PreFooter", template.HTML(sc.PreFooter))
	return td
}

func (s *Session) getHelpers() template.FuncMap {
	return template.FuncMap{
		"lang": s.Lang,
		// "time":     hfTime,
		// "slug":     hfSlug,
		// "mod":      hfMod,
		// "gravatar": hfGravatar,
	}
}

func (s *Session) Lang(text string) string {
	return s.ln.Lang(text)
}

func (s *Session) render(w http.ResponseWriter, filenames ...string) error {
	t := template.New("layout.html")
	// Add helper functions
	t.Funcs(s.getHelpers())
	// Add pad
	s.td.Set("Path", s.path)
	return template.Must(t.ParseFiles(filenames...)).Execute(w, s.td)
}

func (td TemplateData) Set(name string, value interface{}) {
	td[name] = value
}

func (s *Session) Set(name string, value interface{}) {
	s.td.Set(name, value)
}

func (s *Session) AddPath(url, label string) {
	s.path = append(s.path, &PathLink{url, label})
}

func (s *Session) Logged() bool {
	session, _ := s.store.Get(s.r, sessionName)
	userID, found := session.Values["userId"]
	if found {
		s.userID, _ = userID.(int)
	}
	return found
}

func (s *Session) logInUser(userID int, w http.ResponseWriter) {
	session, _ := s.store.Get(s.r, sessionName)
	session.Values["userId"] = userID
	session.Save(s.r, w)
}

func (s *Session) logOutUser(w http.ResponseWriter) {
	session, _ := s.store.Get(s.r, sessionName)
	delete(session.Values, "userId")
	session.Save(s.r, w)
}
