package main

import (
	"fmt"
	"os"

	"github.com/martindrlik/play/config"
)

func loadConfig(name string) (config.Config, error) {
	f, err := os.Open(name)
	if err != nil {
		return config.Config{}, fmt.Errorf("unable to load config %q: %w", name, err)
	}
	defer f.Close()
	return config.Load(f)
}
