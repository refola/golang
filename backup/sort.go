// Stuff for sorting backups...

package backup

import (
	"fmt"
	"github.com/refola/golang/util"
)

// Type for slice of BackupConfigs for use with a (topological) sorting algorithm
type bCfgSlice []BackupConfig

func (b bCfgSlice) Len() int {
	return len(b)
}

// Is b[i] before, after, or indeterminately ordered with respect to b[j]?
func (b bCfgSlice) Compare(i, j int) util.Comparation {
	// b[i] should be backed up before b[j] if one of b[i].To and b[j].From is inside the other.
	if oneContainsOther(b[i].To, b[j].From) {
		return util.Less
	} else if oneContainsOther(b[i].From, b[j].To) {
		return util.Greater
	} else {
		return util.Other
	}
}

// Sort by origin and then destination names, for the sort package.
func (b bCfgSlice) Less(i, j int) bool {
	// Determine if a should be sorted before b.
	before := func(a, b string) bool {
		// Compare the bytes.
		for i := 0; i < len(a) && i < len(b); i++ {
			if a[i] != b[i] {
				return a[i] < b[i]
			}
		}
		// true iff a=b[:len(a)] and false iff b=a[:len(b)]
		return len(a) < len(b)
	}
	if b[i].FromName != b[j].FromName {
		return before(b[i].FromName, b[j].FromName)
	} else {
		return before(b[i].ToName, b[j].ToName)
	}
}
func (b bCfgSlice) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// Sort the backups so that each one that if, for backups X and Y, doing X affects the content of Y, then X is done before Y.
func SortBackups(baks []BackupConfig) {
	//util.Sort(bCfgSlice(baks))
	err := util.Topological(bCfgSlice(baks))
	if err != nil {
		fmt.Printf("Error in topological sort: %v\n Continuing anyway.", err)
	}
}
