// Refola Backup runs something to perform every valid backup as configured in some file(s). Currently it runs rsync and is configured via JSON.

package main

import (
	"github.com/refola/golang/backup"
	"fmt"
	"os"
)

// Show error and exit program if err!=nil
func e(action string, err error) {
	if err != nil {
		fmt.Printf("Error during %s: %s\n", action, err.Error())
		os.Exit(1)
	}
}

func main() {
	// TODO for 1.0: document, make elegant, implement sorting, show status of each origin/destination location
	//backup.DEBUG = true
	fmt.Println("Running Refola Backup.")

	// Config
	cfg, err := backup.GetConfig()
	e("config getting", err)

	// Getting all valid backups
	baks := backup.MakeBackups(cfg)
	//backup.SortBackups(baks)

	// Doing backups
	for _, b := range baks {
		e("backing up", backup.Backup(b))
	}

	// Summarizing actions
	fmt.Println("\n\n\nBackups done:")
	for i, b := range baks {
		fmt.Printf("%d:\t%s to %s.\n", i, b.FromName, b.ToName)
	}
}
