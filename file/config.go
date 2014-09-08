// Get configuration stuff for other Refola things

package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

var ConfigBase string

func init() {
	u, err := user.Current()
	if err != nil {
		fmt.Println("Error getting user!\n" + err.Error())
		os.Exit(1)
	}
	ConfigBase = u.HomeDir + "/.config/refola"
}

// unmarshal ConfigBase/pkgName/config.json into the value pointed at by receiver
func FromJSON(pkgName string, receiver interface{}) error {
	jsonPath := ConfigBase + "/" + pkgName + "/config.json"
	b, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, receiver)
}
