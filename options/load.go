package options

import (
	"encoding/json"
	"log"
	"os"
)

func Load(name string) (opt Options, err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	dec := json.NewDecoder(f)
	err = dec.Decode(&opt)
	return
}

func Must(opt Options, err error) Options {
	if err != nil {
		log.Fatal(err)
	}
	return opt
}
