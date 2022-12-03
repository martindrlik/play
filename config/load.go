package config

import (
	"encoding/json"
	"io"
)

var (
	DefaultRequestLimit uint = 250
)

// Load loads json configuration given by reader r.
func Load(r io.Reader) (config Config, err error) {
	dec := json.NewDecoder(r)
	err = dec.Decode(&config)
	if err == nil {
		if config.RequestLimit <= 0 {
			config.RequestLimit = DefaultRequestLimit
		}
	}
	return
}
