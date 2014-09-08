// Do the actual backups.

package backup

import (
	"fmt"
	"os"
	"os/exec"
)

// Run a backup -- currently just calls rsync.
// TODO: Priority:Last -- Add support for other backup programs?
func Backup(cfg BackupConfig) error {
	// Add terminal slashes and name-based subfolder for destination
	cfg.From = cfg.From + "/"
	cfg.To = cfg.To + "/" + cfg.FromName + "/"

	fmt.Printf("\n\n\nBacking up %s to %s.\n", cfg.FromName, cfg.ToName)
	return rsync(cfg)
}

// Call rsync to do the real work.
func rsync(cfg BackupConfig) error {
	filterPath := cfgDir() + "/rules/rsync"

	// rsync args are kinda fussy, but at least I got programatic multiple filters working, unlike in Bash.
	for _, v := range cfg.Filters {
		cfg.Args = append(cfg.Args, fmt.Sprintf("--filter=merge %s/%s", filterPath, v))
	}
	cfg.Args = append(cfg.Args, cfg.From, cfg.To)
	cmd := exec.Command("rsync", cfg.Args...)

	// show rsync's output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if DEBUG {
		fmt.Printf("Debug: Not actually doing backup. Here's what would've been used.\n%v\n", cmd.Args)
		return nil
	}

	fmt.Printf("Running %v.\n", cmd.Args)
	return cmd.Run()
}
