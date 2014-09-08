// Copy the backup config stuff from files into data structures.

package backup

import (
	"code.google.com/p/refola/file"
	//"os/user" // TODO: Use this to replace {user} and {home} in future configs
)

// Substitutes different things for [] expressions, as follows:
//	[.]		configuration directory (where config.json is)
// TODO	[user]		username of running user
// TODO	[home]		home folder of running user
// TODO	[home=user]	home folder of named user
// TODO	[count=#]	appropriate numbered subfolder, tracked by a file....
// TODO	[.foldername]	the folder for "foldername" as defined ... somewhere
func bracketSub(s string) string {
	if len(s) < 3 {
		return s
	} else if s[:3] == "[.]" {
		return cfgDir() + s[3:]
	}
	return s
}

func cfgPkgName() string {
	name := "backup"
	if DEBUG {
		name += "/debug"
	}
	return name
}

func cfgDir() string {
	return file.ConfigBase + "/" + cfgPkgName()
}

//Get everything needed to figure out what commands are valid to run with this program.
func GetConfig() (*Config, error) {
	var cfg Config
	err := file.FromJSON(cfgPkgName(), &cfg)
	if err != nil {
		return nil, ebld("Could not retrieve config!", err)
	}

	for _, origin := range cfg.Origins {
		for i, v := range origin.From {
			origin.From[i] = bracketSub(v)
		}
	}
	for _, dest := range cfg.Destinations {
		for i, v := range dest.To {
			dest.To[i] = bracketSub(v)
		}
	}

	return &cfg, nil
}
