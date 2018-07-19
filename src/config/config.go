package config

import (
	"encoding/json"
	"os"
)

// Load (path string) interface - load config file and return interface
func Load(c interface{}, path string) {
	fi, err := os.Open(path)
	defer fi.Close()

	if err != nil {
		return
	}

	jsonParser := json.NewDecoder(fi)
	jsonParser.Decode(&c)
}
