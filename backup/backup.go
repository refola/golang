// Do the actual backups.

package backup

import (
	"fmt"
	"os"
	"os/exec"
)

// Run a backup -- currently just calls rsync.
// TODO: Add support for other backup programs, like btrfs
func Backup(cfg BackupConfig) error {
	// Add terminal slashes and name-based subfolder for destination
	cfg.From = cfg.From + "/"
	cfg.To = cfg.To + "/" + cfg.FromName + "/"

	fmt.Printf("\n\n\nBacking up %s to %s.\n", cfg.FromName, cfg.ToName)
	return rsync(cfg)
}
