package wiki

import (
	"html/template" // TODO: switch to updated template package
	"io"
	"net/http"
	"regexp"

	// local
	"server/core"
)

type Wiki struct {
	core.Server
	pages *pageStore
}

func New(prefix string) *Wiki {
	wi := &Wiki{Server: *core.NewServer(prefix)}
	for path, handler := range handlers {
		wi.AddHandler(path, wi.makeHandler(path, handler))
	}
	wi.addTemplates() // TODO: template redo
	wi.pages = newPageStore(wi.Data)

	return wi
}

// TODO: Upgrade to new template package.
// // inserts the template that matches the given field's name
// func tmplFormatter(w io.Writer, s string, data ...interface{}) {
//
// }
// // inserts braces surrounding the name of the given field
// func braceFormatter(w io.Writer, s string, data ...interface{}) {
// 	w.Write([]byte("{"+s+"}"))
// }

// Add the templates for the wiki, adding the header as appropriate.
func (wi *Wiki) addTemplates() { // TODO: template redo
	formatters := map[string]func(io.Writer, string, ...interface{}){
		"":      template.StringFormatter, // TODO: template redo
		"html":  template.HTMLFormatter,   // TODO: template redo
		"token": tokenFormatter,
		"wiki":  wikiFormatter}

	prefix := "tmpl/"
	list, _ := wi.Data.ListData(prefix)
	tmpls := make(map[string]string)
	for _, v := range list {
		data, _ := wi.Data.Load(prefix + v)
		tmpls[v] = string(data)
	}
	h := "{header}"
	lenh := len(h)
	// Pass through data a second time so the header will have been already loaded.
	for i, v := range tmpls {
		if len(v) > lenh && v[:lenh] == h {
			v = tmpls["header"] + v[lenh:]
		}
		tmpl, _ := template.Parse(v, formatters) // TODO: template redo
		wi.Templates[i] = tmpl
	}
}

func (wi *Wiki) renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := wi.Templates[tmpl].Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (wi *Wiki) redirect(w http.ResponseWriter, r *http.Request, to string) {
	redir := "/" + wi.Prefix + "/" + to
	http.Redirect(w, r, redir, http.StatusFound)
}

func (wi *Wiki) authenticate(user, pass string) bool {
	fPass, err := wi.Data.Load("users/" + user)
	if err != nil || string(fPass) != pass {
		return false
	}
	return true
}
func (wi *Wiki) getUser(r *http.Request) string {
	var user, pass string
	for _, v := range r.Cookies() {
		switch v.Name {
		case "user":
			user = v.Value
		case "pass":
			pass = v.Value
		}
	}
	if wi.authenticate(user, pass) {
		return user
	}
	return r.RemoteAddr
}

var userCheck = regexp.MustCompile("^[a-zA-Z0-9]+$")

// TODO: replace with core.User and core.UserStore
func (wi *Wiki) makeAccount(user, pass string) bool {
	if !userCheck.Match([]byte(user)) {
		return false
	}
	if data, _ := wi.Data.Load("users/" + user); data != nil { // if user already exists
		return false
	}
	if err := wi.Data.Save("users/"+user, []byte(pass)); err != nil {
		return false
	}
	return true
}
