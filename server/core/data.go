// TODO: Move all functionality to refola/file?

package core

import (
	"io/ioutil"
	"os"
	"sort"
	"sync"
)

// // TODO: Move mutexes to the data stores. Global state is bad in this case.
// // Commented out to cause compiler errors where this isn't yet fixed.
// var dataLockers map[string]*sync.RWMutex
// var dataLockersLocker *sync.RWMutex
//
// func init() {
// 	dataLockers = make(map[string]*sync.RWMutex)
// 	dataLockersLocker = new(sync.RWMutex)
// }

func lockData(where ...string) {
	dataLockersLocker.RLock() // prevent concurrent changes to dataLockers
	defer dataLockersLocker.RUnlock()
	for _, s := range where {
		if nil == dataLockers[s] {
			dataLockers[s] = new(sync.RWMutex)
		}
		dataLockers[s].Lock()
		// TODO: handle case where dataLockers[s] is already locked;
		// is-as, if it is, dataLockersLocker being locked will prevent
		// unlockData() from unlocking dataLockers[s]
	}
}

func unlockData(where ...string) {
	dataLockersLocker.RLock()
	defer dataLockersLocker.RUnlock()
	for _, s := range where {
		dataLockers[s].Unlock()
	}
}

var dataPrefix = "/dev/null/fake/path/" // causes errors if the prefix is not set before use
func SetDataPrefix(prefix string) {
	dataPrefix = endingSlash(prefix)
}

func endingSlash(s string) string {
	switch l := len(s); true {
	case l == 0:
		return "/"
	case s[l-1] != '/':
		return s + "/"
	default:
		return s
	}
	panic("unreachable")
}

// "Loc" is short for "Location"
// "Data" and "Loc" correspond to file and directory in the only implementation
// so far, but it's quite possible for a database or something to implement it in
// a way that doesn't translate directly to files/folders.
type DataStore interface {
	Save(where string, what []byte) error
	Append(where string, what []byte) error
	Load(where string) ([]byte, error)
	Remove(where string) error
	Rename(from, to string) error
	MakeLoc(where string) error
	RemLoc(where string) error // REMove LOCation
	ListLocs(where string) ([]string, error)
	ListData(where string) ([]string, error)
}

// The first (and so far only) implementation of DataStore.
// It's insecure from evil servers giving other servers' names,
// but it centralizes data access and safely handles special characters.
type fileStore struct {
	path  string
	mutex *sync.RWMutex // TODO: Use *this* mutex.
}

// Get a DataStore backed by the local filesystem. "where" is the path to a directory where the DataStore can store, read, modify, and generally do stuff to data. "name" is the unique name for the DataStore.
func NewFileDataStore(where, name string) DataStore {
	return &fileStore{path: endingSlash(where) + endingSlash(Sanitize(name))}
}

func (fs *fileStore) lock(where ...string) {
	lockData(fs.pathsFor(where)...)
}
func (fs *fileStore) unlock(where ...string) {
	unlockData(fs.pathsFor(where)...)
}
func (fs *fileStore) pathFor(where string) string {
	return fs.path + Sanitize(where)
}
func (fs *fileStore) pathsFor(where []string) []string {
	ret := make([]string, len(where))
	for i, v := range where {
		ret[i] = fs.pathFor(v)
	}
	return ret
}

// save the given data to the centralized place
func (fs *fileStore) Save(where string, what []byte) error {
	fs.lock(where)
	defer fs.unlock(where)
	return ioutil.WriteFile(fs.pathFor(where), what, 0600)
}
func (fs *fileStore) Append(where string, what []byte) error {
	fs.lock(where)
	defer fs.unlock(where)
	path := fs.pathFor(where)
	old, _ := ioutil.ReadFile(path) // ignore error; probably non-existent file
	return ioutil.WriteFile(path, append(old, what...), 0600)
}

// load the requested data from the centralized place
func (fs *fileStore) Load(where string) ([]byte, error) {
	fs.lock(where)
	defer fs.unlock(where)
	return ioutil.ReadFile(fs.pathFor(where))
}
func (fs *fileStore) Remove(where string) error {
	fs.lock(where)
	defer fs.unlock(where)
	return os.Remove(fs.pathFor(where))
}
func (fs *fileStore) Rename(from, to string) error {
	fs.lock(from, to)
	defer fs.unlock(from, to)
	return os.Rename(fs.pathFor(from), fs.pathFor(to))
}
func (fs *fileStore) MakeLoc(where string) error {
	fs.lock(where)
	defer fs.unlock(where)
	return os.MkdirAll(fs.pathFor(where), 0600)
}
func (fs *fileStore) RemLoc(where string) error {
	fs.lock(where)
	defer fs.unlock(where)
	return os.Remove(fs.pathFor(where))
}

// lists available data in the centralized place
func (fs *fileStore) list(where string, keep func(os.FileInfo) bool) ([]string, error) {
	fs.lock(where)
	defer fs.unlock(where)
	dirName := endingSlash(fs.pathFor(where))
	contents, err := ioutil.ReadDir(dirName)
	if err != nil {
		return nil, err
	}
	trash := 0
	for i, v := range contents {
		if !keep(v) {
			contents[i] = nil
			trash++
		}
	}
	names := make([]string, len(contents)-trash)
	trashPassed := 0
	for i, v := range contents {
		if nil == v {
			trashPassed++
		} else {
			names[i-trashPassed] = Unsanitize(v.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

// list data items in location
func (fs *fileStore) ListData(where string) ([]string, error) {
	return fs.list(where, func(f os.FileInfo) bool {
		return !f.IsDir()
	})
}

// list location items in the location
func (fs *fileStore) ListLocs(where string) ([]string, error) {
	return fs.list(where, func(f os.FileInfo) bool {
		return f.IsDir()
	})
}

// Prepends a string to all accesses of a given DataStore.
// Useful when several functions deal with the same location.
type subStore struct {
	where string
	data  DataStore
}

func NewSubDataStore(where string, data DataStore) DataStore {
	return &subStore{where: endingSlash(Sanitize(where)), data: data}
}

// minimally call respective methods of given DataStore
func (ss *subStore) Save(where string, what []byte) error {
	return ss.data.Save(ss.where+where, what)
}
func (ss *subStore) Append(where string, what []byte) error {
	return ss.data.Append(ss.where+where, what)
}
func (ss *subStore) Load(where string) ([]byte, error) {
	return ss.data.Load(ss.where + where)
}
func (ss *subStore) Remove(where string) error {
	return ss.data.Remove(ss.where + where)
}
func (ss *subStore) Rename(from, to string) error {
	return ss.data.Rename(ss.where+from, ss.where+to)
}
func (ss *subStore) MakeLoc(where string) error {
	return ss.data.MakeLoc(ss.where + where)
}
func (ss *subStore) RemLoc(where string) error {
	return ss.data.RemLoc(ss.where + where)
}
func (ss *subStore) ListLocs(where string) ([]string, error) {
	return ss.data.ListLocs(ss.where + where)
}
func (ss *subStore) ListData(where string) ([]string, error) {
	return ss.data.ListData(ss.where + where)
}
