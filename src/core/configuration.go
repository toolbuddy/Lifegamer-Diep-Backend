package core

import (
	"encoding/json"
	"os"
)

// define the ServerConfiguration struct
type ServerConfiguration struct {
	Host string
	Port string
}
// define the Configuration struct
type Configuration struct {
	Server ServerConfiguration
}
// define the loadFromFile function in Configuration pointer
func (c *Configuration) loadFromFile() error {
	file, _ := os.Open("src/config/main.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&c)
	if err != nil {
		return err
	}
	return nil
}
