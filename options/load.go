package options

import (
	"encoding/json"
	"os"
)

// Load loads options from file given by name.
func Load(name string) (opt Options, err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	dec := json.NewDecoder(f)
	err = dec.Decode(&opt)
	return
}
