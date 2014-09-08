package wiki

import (
	"net/http"

	// local
	"server/core"
)

type handler func(http.ResponseWriter, *http.Request, string, *Wiki)

var handlers = map[string]handler{"view/": viewHandler, "edit/": editHandler, "save/": saveHandler, "delete/": deleteHandler, "index": indexHandler, "login": loginHandler, "newacc": newaccHandler, "token/": tokenHandler, "": defaultHandler}

func tokenHandler(w http.ResponseWriter, r *http.Request, title string, wi *Wiki) {
	if title == "" {
		defaultHandler(w, r, title, wi)
		return
	}
	p, err := wi.pages.load(title)
	if err != nil {
		wi.redirect(w, r, "edit/"+title)
		return
	}
	wi.renderTemplate(w, "token", p)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string, wi *Wiki) {
	if title == "" {
		defaultHandler(w, r, title, wi)
		return
	}
	p, err := wi.pages.load(title)
	if err != nil {
		wi.redirect(w, r, "edit/"+title)
		return
	}
	wi.renderTemplate(w, "view", p)
}
func editHandler(w http.ResponseWriter, r *http.Request, title string, wi *Wiki) {
	p, err := wi.pages.load(title)
	if err != nil {
		p = &page{Title: title}
	}
	wi.renderTemplate(w, "edit", p)
}
func saveHandler(w http.ResponseWriter, r *http.Request, title string, wi *Wiki) {
	e := func(s string) {
		http.Error(w, "Error: "+s+". Could not save \""+title+"\".", http.StatusInternalServerError)
	}
	body := r.FormValue("body")
	p := &page{Title: title, Body: []byte(body)}
	if body == "" {
		wi.renderTemplate(w, "save", p)
		return
	}
	old, err := wi.pages.load(title)
	if err == nil { // assumes err is from the old page not existing
		err = wi.pages.log(wi.getUser(r), "overwrote", old)
		if err != nil { // do not attempt overwriting old page without logging
			e(err.Error())
			return
		}
	}
	err = wi.pages.save(p)
	if err != nil {
		e(err.Error())
		return
	}
	wi.redirect(w, r, "view/"+title)
}
func deleteHandler(w http.ResponseWriter, r *http.Request, title string, wi *Wiki) {
	p, err := wi.pages.load(title)
	if err != nil {
		defaultHandler(w, r, title, wi)
		return
	}
	wi.renderTemplate(w, "delete", p)
	err = wi.pages.log(wi.getUser(r), "deleted", p)
	if err != nil {
		return
	}
	wi.pages.delete(p)
}

func indexHandler(w http.ResponseWriter, r *http.Request, title string, wi *Wiki) {
	v, err := wi.pages.list()
	if err != nil {
		http.Error(w, "Could not get page list", http.StatusInternalServerError)
		return
	}
	wi.renderTemplate(w, "index", v)
}

func loginHandler(w http.ResponseWriter, r *http.Request, title string, wi *Wiki) {
	if r.Method != "POST" {
		wi.renderTemplate(w, "login", nil)
		return
	}
	user, pass := r.FormValue("user"), r.FormValue("pass")
	if !wi.authenticate(user, pass) {
		// TODO: add invalid login logging
		wi.renderTemplate(w, "loginErr", nil)
		return
	}
	// set login cookies
	header := w.Header() // TODO: make cookies more permanent in browser
	header.Add("Set-Cookie", "user="+core.Sanitize(user))
	header.Add("Set-Cookie", "pass="+core.Sanitize(pass))
	wi.renderTemplate(w, "loginGood", nil)
}
func newaccHandler(w http.ResponseWriter, r *http.Request, title string, wi *Wiki) {
	if r.Method != "POST" {
		wi.renderTemplate(w, "newAcc", nil)
		return
	}
	if !wi.makeAccount(r.FormValue("user"), r.FormValue("pass")) {
		wi.renderTemplate(w, "newAccErr", nil)
		return
	}
	wi.renderTemplate(w, "newAccGood", nil)
}

func defaultHandler(w http.ResponseWriter, r *http.Request, title string, wi *Wiki) {
	wi.redirect(w, r, "view/Home")
}

func (wi *Wiki) makeHandler(path string, fn handler) http.HandlerFunc {
	lenPath := len("/" + wi.Prefix + "/" + path)
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[lenPath:]
		fn(w, r, title, wi)
	}
}
