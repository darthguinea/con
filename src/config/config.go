package config

import (
	"encoding/json"
	"io/ioutil"

	"../log"
)

// Load (path string) interface - load config file and return interface
func Load(c interface{}, path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("Error loading file [%v]", path)
	}

	json.Unmarshal(data, &c)
}
