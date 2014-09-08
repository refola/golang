package core

import (
	"html/template" // TODO: switch to updated template package
	"net/http"
	"os"
	"strconv"
	"time"
)

type Server struct {
	Prefix    string                        // should not be changed, but needs to be easily accessible
	Templates map[string]*template.Template // TODO: template redo
	Data      DataStore
	Users     *UserStore
}

func NewServer(name string) *Server {
	prefix := Sanitize(name)
	tmpl := make(map[string]*template.Template) // TODO: template redo
	data := NewDataStore(prefix)
	s := &Server{Prefix: prefix, Templates: tmpl, Data: data} // TODO: template redo
	return s
}

// convenience method to add all templates found in a storage place
func (s *Server) AddTemplates(where string) error {
	list, err := s.Data.ListData(where)
	if err != nil {
		return err
	}
	for _, name := range list {
		if name[len(name)-1] == '~' {
			continue
		}
		data, err := s.Data.Load(where + name)
		if err != nil {
			return err
		}
		tmpl, err := template.New(name).Parse(string(data))
		if err != nil {
			return err
		}
		s.Templates[name] = tmpl
	}
	return nil
}

func (s *Server) AddHandler(path string, handler http.HandlerFunc) {
	url := s.Prefix + "/" + path
	if s.Prefix != "" { // if server isn't operating at site's root
		url = "/" + url
	}
	http.HandleFunc(url, handler)
}

func (s *Server) AddHandlers(handlers map[string]http.HandlerFunc) {
	for path, handler := range handlers {
		s.AddHandler(path, handler)
	}
}

// Convenience function to make a new default server with some common default handlers
func NewDefaultServer(defaultServer string, includeExit bool) {
	s := NewServer("")
	handlers := map[string]http.HandlerFunc{"": func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, defaultServer, http.StatusFound) }, "favicon.ico": func(w http.ResponseWriter, r *http.Request) { http.NotFound(w, r) }}

	if includeExit {
		handlers["exit"] = func(w http.ResponseWriter, r *http.Request) {
			s.Data.Append("shutdowns", []byte(r.RemoteAddr+" on "+time.Now().UTC().Format(time.RFC1123)+"\n"))
			os.Exit(0)
		}
	}
	s.AddHandlers(handlers)
}

func StartServers(port uint) error {
	return http.ListenAndServe(":"+strconv.FormatUint(uint64(port), 10), nil)
}
