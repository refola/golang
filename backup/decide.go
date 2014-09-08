// Do a bunch of algorithmic stuff to convert backup config data structures into backup commands and sort backups in the appropriate order.

package backup

import (
	"io/ioutil"
)

// Check if one of the strings (paths) contains the other.
func oneContainsOther(a, b string) bool {
	if len(a) > len(b) {
		return b == a[:len(b)]
	} else {
		return a == b[:len(a)]
	}
}

// Build the config(s) for a single origin-destination-backuptype combination.
func makeBakConfigs(cfg *Config, origin Origin, dest Destination, bakType BackupType) []BackupConfig {
	var bcfgs []BackupConfig
	for _, from := range origin.From {
		for _, to := range dest.To {
			var bcfg BackupConfig
			bcfg.Args = append(bcfg.Args, cfg.BackupOptions["main"]...)
			bcfg.Args = append(bcfg.Args, []string(cfg.BackupOptions[cfg.Filesystems[dest.Filesystem]])...)
			bcfg.Filters = append(bcfg.Filters, origin.Rules[bakType]...)
			bcfg.From = from
			bcfg.FromName = origin.Name
			bcfg.To = to
			bcfg.ToName = dest.Name
			bcfgs = append(bcfgs, bcfg)
		}
	}
	return bcfgs
}

// Filter backups by validity -- parent-child relations, non-existent folders, (etc?)
func validBackups(cfgs []BackupConfig) []BackupConfig {
	var valids []BackupConfig
	for _, cfg := range cfgs {
		// check for parent-child relations in backup to/from folders
		if oneContainsOther(cfg.From, cfg.To) {
			continue
		}

		// check for nonexistent/inaccessible folders
		if _, err := ioutil.ReadDir(cfg.From); err != nil {
			continue
		}
		if _, err := ioutil.ReadDir(cfg.To); err != nil {
			continue
		}

		// All tests passed! It's a valid BackupConfig.
		valids = append(valids, cfg)
	}
	return valids
}

// Build all valid backups.
func MakeBackups(cfg *Config) []BackupConfig {
	var cfgs []BackupConfig

	for _, dest := range cfg.Destinations {
		for _, dtype := range dest.Types {
			for _, or := range cfg.Origins {
				if _, okay := or.Rules[dtype]; okay {
					cfgs = append(cfgs, makeBakConfigs(cfg, or, dest, dtype)...)
				}
			}
		}
	}

	return validBackups(cfgs)
}
