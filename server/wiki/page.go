package wiki

import (
	"time"

	// local
	"server/core"
)

type page struct {
	Title string
	Body  []byte
}

type pageStore struct {
	data  core.DataStore
	pages core.DataStore
}

func newPageStore(store core.DataStore) *pageStore {
	return &pageStore{data: store, pages: core.NewSubDataStore("pages/", store)}
}

func (ps *pageStore) list() ([]string, error) {
	return ps.pages.ListData("")
}

func (ps *pageStore) load(title string) (*page, error) {
	body, err := ps.pages.Load(title)
	if err != nil {
		return nil, err
	}
	return &page{Title: title, Body: body}, nil
}

func (ps *pageStore) save(p *page) error {
	return ps.pages.Save(p.Title, p.Body)
}

func (ps *pageStore) delete(p *page) error {
	return ps.pages.Remove(p.Title)
}

// logs current page contents, such as before deletion/editing
func (ps *pageStore) log(user, action string, p *page) error {
	toAdd := "==" + user + " " + action + " and eliminated old \"" + p.Title + "\" on " + time.Now().UTC().Format(time.RFC1123) + "==\n" + string(p.Body) + "\n\n"
	return ps.data.Append("log/"+p.Title, []byte(toAdd))
}
