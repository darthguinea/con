package config

import (
	"encoding/json"
	"os"
)

// Load (path string) interface - load config file and return interface
func Load(path string) interface{} {
	var v interface{}
	fi, err := os.Open(path)
	defer fi.Close()

	if err != nil {
		return nil
	}
	jsonParser := json.NewDecoder(fi)
	jsonParser.Decode(&v)
	return v
}
