package file

// TODO: move to refola/storage

import (
	"fmt"
	"net/http"

	// local
	"server/core"
)

type FileServer struct {
	core.Server
}

func New(prefix string) *FileServer {
	fs := &FileServer{*core.NewServer(prefix)}
	fs.AddHandler("", getHandler(fs))
	err := fs.AddTemplates("tmpl/") // TODO: template redo
	if err != nil {
		panic(err)
	}
	fs.Data = core.NewSubDataStore("files/", fs.Data)
	return fs
}

func getHandler(fs *FileServer) func(http.ResponseWriter, *http.Request) {
	pre := "/" + fs.Prefix + "/"
	lenPre := len(pre)
	get := "get/"
	lenGet := len(get)
	upload := "upload/"
	lenUpload := len(upload)
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[lenPre:]
		lenTitle := len(title)
		switch {
		case lenTitle >= lenGet && title[:lenGet] == get:
			fs.get(w, r, title[lenGet:])
		case lenTitle >= lenUpload && title[:lenUpload] == upload:
			fs.upload(w, r, title[lenUpload:])
		default:
			http.Redirect(w, r, "/"+fs.Prefix+"/"+get, http.StatusFound)
		}
	}
}

func (fs *FileServer) get(w http.ResponseWriter, r *http.Request, where string) {
	if len(where) == 0 || where[len(where)-1] == '/' {
		fs.getList(w, r, where)
	} else {
		fs.getFile(w, r, where)
	}
}
func (fs *FileServer) getList(w http.ResponseWriter, r *http.Request, where string) {
	locs, err := fs.Data.ListLocs(where)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "File listing error: could not get list of subdirectories", http.StatusInternalServerError)
	}
	data, err := fs.Data.ListData(where)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "File listing error: could not get list of files", http.StatusInternalServerError)
	}
	for i, v := range locs {
		locs[i] = v + "/"
	}
	list := append(locs, data...)
	fmt.Println(list)
	err = fs.Templates["list"].Execute(w, list) // TODO: template redo
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (fs *FileServer) getFile(w http.ResponseWriter, r *http.Request, where string) {
	// TODO: implement
	http.Error(w, "Downloading not implemented", http.StatusNotImplemented)
}

func (fs *FileServer) upload(w http.ResponseWriter, r *http.Request, where string) {
	// TODO: implement
	http.Error(w, "Uploading not implemented", http.StatusNotImplemented)
}
