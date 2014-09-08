// General backup package stuff that hasn't (yet) been moved to a more descriptive file

package backup

import (
	"errors"
)

// Debug mode effects:
// 	GetConfig() tells refola/util.FromJSON to get config from "backup/debug"
// 	backup() just prints the command instead of backing up
var DEBUG = true

// A place to hold everything from config.json.
type Config struct {
	BackupOptions map[string]BackupOpts
	Filesystems   map[string]string // which file systems use which BackupOpts
	Origins       []Origin
	Destinations  []Destination
}

type BackupOpts []string // parameters to pass to backup program
// Everything we need to know about backing up from somewhere.
type Origin struct {
	Name  string                  // backup goes to "Destination.To/Origin.Name/"
	From  []string                // where to look for copying files from
	Rules map[BackupType][]string // BackupType-named lists of rules files
}

// Everything we need to know about backing up to somewhere.
type Destination struct {
	Name       string       // the name to display when backing up to Destination
	To         []string     // where to look for copying files to
	Types      []BackupType // rule type valid for this destination
	Filesystem string
}
type BackupType string // the types of backups to do -- sorta a tagging system for matching compatible origins and destinations and using the right rsync rules. separate type to avoid accidental mixing with other strings

// Info for doing one backup
type BackupConfig struct {
	Args     []string
	Filters  []string
	From     string
	FromName string
	To       string
	ToName   string
}

// Build a more descriptive error from a message and an existing error.
func ebld(msg string, err error) error {
	return errors.New(msg + " " + err.Error())
}
